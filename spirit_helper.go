package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogap/spirit"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

var (
	ErrNoURNPackageSourceFound = errors.New("no urn packages source found")
	ErrConfigFileNameIsEmpty   = errors.New("config file name is empty")
)

type SpiritHelper struct {
	conf           spirit.SpiritConfig
	configFile     string
	configFileName string
	originalConfig []byte

	RefURNs     []string
	RefPackages []Package
}

func (p *SpiritHelper) LoadSpiritConfig(filename string) (err error) {

	if filename == "" {
		err = ErrConfigFileNameIsEmpty
		return
	}

	if fi, e := os.Stat(filename); e != nil {
		err = e
		return
	} else {
		p.configFile = filename
		p.configFileName = fi.Name()
	}

	if p.originalConfig, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(p.originalConfig, &p.conf); err != nil {
		return
	}

	return
}

func (p *SpiritHelper) CreateProject(createOpts CreateOptions, tmplArgs map[string]interface{}) (err error) {
	if err = createOpts.Validate(); err != nil {
		return
	}

	goSrc := path.Join(createOpts.GoPath, "src")

	if err = p.parse(goSrc, createOpts.Sources); err != nil {
		return
	}

	// download packages
	if createOpts.GetPackages {
		if err = p.GetPackages(createOpts.PackagesRevision, createOpts.UpdatePackages); err != nil {
			return
		}
	}

	// make project dir
	projectPath := path.Join(goSrc, createOpts.ProjectPath)
	if path.IsAbs(createOpts.ProjectPath) {
		projectPath = createOpts.ProjectPath
	}

	if fi, e := os.Stat(projectPath); e != nil {
		if !strings.Contains(e.Error(), "no such file or directory") &&
			!os.IsNotExist(e) {
			err = e
			return
		}
	} else if !fi.IsDir() {
		err = fmt.Errorf("your project path %s already exist, but it is not a directory", projectPath)
		return
	} else if createOpts.ForceWrite {
		spirit.Logger().Warnf("project path %s already exist, it will be overwrite", projectPath)
	} else {
		err = fmt.Errorf("your project path %s already exist", projectPath)
		return
	}

	if err = os.MkdirAll(projectPath, os.FileMode(0755)); err != nil {
		if !os.IsNotExist(err) {
			return
		} else if !createOpts.ForceWrite {
			return
		}
		err = nil
	}

	// render code template
	tmplPathFmt := "github.com/gogap/spirit-tool/template/%s/main.go"
	tmplArgsPathFmt := "github.com/gogap/spirit-tool/template/%s/args.json"

	tmplPath := path.Join(goSrc, fmt.Sprintf(tmplPathFmt, createOpts.TemplateName))
	spirit.Logger().Infof("using template of %s: %s", createOpts.TemplateName, tmplPath)

	tmplArgsPath := path.Join(goSrc, fmt.Sprintf(tmplArgsPathFmt, createOpts.TemplateName))
	spirit.Logger().Infof("using template args of %s: %s", createOpts.TemplateName, tmplArgsPath)

	var tmpl *template.Template
	if tmpl, err = template.New("main.go").Option("missingkey=error").Delims("//<-", "->//").ParseFiles(tmplPath); err != nil {
		return
	}

	internalArgs := map[string]interface{}{}
	if argData, e := ioutil.ReadFile(tmplArgsPath); e == nil {
		if err = json.Unmarshal(argData, &internalArgs); err != nil {
			return
		}
	}

	if tmplArgs != nil {
		for k, v := range tmplArgs {
			internalArgs[k] = v
		}
	}

	buffer := &bytes.Buffer{}
	if err = tmpl.Execute(buffer, map[string]interface{}{
		"create_options":  createOpts,
		"packages":        p.RefPackages,
		"config":          p.configFile,
		"config_filename": p.configFileName,
		"create_time":     time.Now(),
		"args":            internalArgs}); err != nil {
		return
	}

	srcPath := path.Join(projectPath, "main.go")
	if err = ioutil.WriteFile(srcPath, buffer.Bytes(), os.FileMode(0644)); err != nil {
		return
	}

	confPath := path.Join(projectPath, p.configFileName)
	if err = ioutil.WriteFile(confPath, p.originalConfig, os.FileMode(0644)); err != nil {
		return
	}

	// format code for sort import packages order
	if _, err = execCommand("go fmt " + srcPath); err != nil {
		return
	}

	spirit.Logger().Infof("project created at %s\n", projectPath)

	return
}

func (p *SpiritHelper) GetPackages(pkgRevision map[string]string, update bool) (err error) {
	for _, pkg := range p.RefPackages {
		if pkgRevision != nil {
			if revision, exist := pkgRevision[pkg.URI]; exist {
				pkg.Revision = revision
			}
		}
		if err = pkg.Get(update); err != nil {
			return
		}
	}
	return
}

func (p *SpiritHelper) BuildProject(createOpts CreateOptions, name string, tmplArgs map[string]interface{}) (err error) {

	if err = p.CreateProject(createOpts, tmplArgs); err != nil {
		return
	}

	cmd := "go build -o "
	if verbosity > 0 {
		cmd = "go build -v -o "
	}

	if _, err = execCommandWithDir(cmd+name+" "+path.Join(createOpts.ProjectPath, "main.go"), createOpts.ProjectPath); err != nil {
		return
	}

	return
}

func (p *SpiritHelper) RunProject(createOpts CreateOptions, detach bool, envs []string, tmplArgs map[string]interface{}) (err error) {
	if err = p.BuildProject(createOpts, "main", tmplArgs); err != nil {
		return
	}

	if cmder, e := execute(path.Join(createOpts.ProjectPath, "main"), createOpts.ProjectPath, !detach, envs); e != nil {
		err = e
		return
	} else if !detach {
		cmder.Wait()
	} else {
		spirit.Logger().Infof("PID: %d\n", cmder.Process.Pid)
	}

	return
}

func (p *SpiritHelper) parse(gosrc string, sources []string) (err error) {
	if sources == nil || len(sources) == 0 {
		err = ErrNoURNPackageSourceFound
		return
	}

	var urns []string

	if urns = parseActorsUsingURN(
		p.conf.InputTranslators,
		p.conf.OutputTranslators,
		p.conf.Inboxes,
		p.conf.Outboxes,
		p.conf.Receivers,
		p.conf.Senders,
		p.conf.Routers,
		p.conf.Components,
		p.conf.LabelMatchers,
		p.conf.URNRewriters,
	); err != nil {
		return
	}

	for _, readerPool := range p.conf.ReaderPools {
		urns = append(urns, parseActorUsingURN(readerPool.ActorConfig)...)
		if readerPool.Reader != nil {
			urns = append(urns, parseActorUsingURN(*readerPool.Reader)...)
		}
	}

	for _, writerPool := range p.conf.WriterPools {
		urns = append(urns, parseActorUsingURN(writerPool.ActorConfig)...)
		if writerPool.Writer != nil {
			urns = append(urns, parseActorUsingURN(*writerPool.Writer)...)
		}
	}

	p.RefURNs = urns

	if p.RefPackages, err = urnsToPackages(gosrc, urns, sources...); err != nil {
		return
	}

	return
}

func parseActorsUsingURN(confs ...[]spirit.ActorConfig) (urns []string) {
	for _, conf := range confs {
		for _, c := range conf {
			urns = append(urns, c.URN)
		}
	}
	return
}

func parseActorUsingURN(actorConfs ...spirit.ActorConfig) (urns []string) {
	for _, conf := range actorConfs {
		urns = append(urns, conf.URN)
	}
	return
}

func urnsToPackages(gosrc string, urns []string, sourceFiles ...string) (packages []Package, err error) {
	urnPkgMap := map[string]string{}

	for _, sourceFile := range sourceFiles {
		var data []byte

		if data, err = ioutil.ReadFile(sourceFile); err != nil {
			return
		}

		sourceConf := SourceConfig{}
		if err = json.Unmarshal(data, &sourceConf); err != nil {
			return
		}

		for _, urnPkg := range sourceConf.Packages {
			if oldVal, exist := urnPkgMap[urnPkg.URN]; exist {
				if oldVal != urnPkg.Pkg {
					err = fmt.Errorf("source have duplicate urn pkg, urn:%s, pkg1:%s, pkg2: %s, file: %s", urnPkg.URN, oldVal, urnPkg.Pkg, sourceFile)
					return
				}
			}
			urnPkgMap[urnPkg.URN] = urnPkg.Pkg
		}
	}

	pkgs := map[string]bool{}

	for _, urn := range urns {
		if pkg, exist := urnPkgMap[urn]; !exist {
			err = fmt.Errorf("no package from any source of urn: %s", urn)
			return
		} else {
			pkgs[pkg] = true
		}
	}

	for pkg, _ := range pkgs {
		packages = append(packages, Package{gosrc: gosrc, URI: pkg, Revision: ""})
	}

	return
}
