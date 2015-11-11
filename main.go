package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/gogap/spirit"
)

var (
	app *cli.App
)

func main() {
	app = cli.NewApp()
	app.Name = "spirit-tool"
	app.Authors = []cli.Author{{"zeal", "xujinzheng@gmail.com"}}
	app.Usage = "help user easily to use spirit"

	app.Commands = []cli.Command{
		//commandUpdate(update),
		//commandFind(find),
		commandCreate(create),
	}

	app.Run(os.Args)
}

func update(context *cli.Context) {
	return
}

func find(context *cli.Context) {
	return
}

func create(context *cli.Context) {
	var tmpl *template.Template
	var err error

	defer func() {
		if err != nil {
			spirit.Logger().Error(err)
			os.Exit(128)
		}
	}()

	goPath := context.String("gopath")
	if goPath == "" {
		err = fmt.Errorf("could not get GOPATH")
		return
	}
	spirit.Logger().Infof("GOPATH: %s", goPath)

	goSrc := goPath + "/src"

	appDir := context.String("path")
	if appDir == "" {
		err = fmt.Errorf("please input your app path, like: github.com/your_orgs/project_name ")
		return
	}

	appPath := fmt.Sprintf("%s/%s", goSrc, appDir)
	if _, e := os.Stat(appPath); e != nil {
		if !os.IsNotExist(e) {
			err = e
			return
		}
	} else {
		err = fmt.Errorf("your project path %s already exist", appPath)
		return
	}

	configFile := context.String("config")

	if configFile == "" {
		err = fmt.Errorf("please input config file")
		return
	}

	configFileName := ""
	if fi, e := os.Stat(configFile); e != nil {
		err = e
		return
	} else {
		configFileName = fi.Name()
	}

	var confData []byte
	if confData, err = ioutil.ReadFile(configFile); err != nil {
		return
	}

	templateName := context.String("template")
	innerConfig := context.Bool("inner")

	tmplPathFmt := "github.com/gogap/spirit-tool/template/%s/main.go"

	tmplPath := path.Join(goSrc, fmt.Sprintf(tmplPathFmt, templateName))

	spirit.Logger().Infof("using template of %s: %s", templateName, tmplPath)

	if tmpl, err = template.New("main.go").Option("missingkey=error").Delims("//<-", "->//").ParseFiles(tmplPath); err != nil {
		return
	}

	spiritConf := spirit.SpiritConfig{}

	if err = json.Unmarshal(confData, &spiritConf); err != nil {
		return
	}

	var urns []string
	if urns, err = parseSpiritConfigUsingURN(spiritConf); err != nil {
		return
	}

	spirit.Logger().Infof("%d urns found", len(urns))

	extSources := context.StringSlice("source")

	sources := []string{
		path.Join(goSrc, "github.com/gogap/spirit-tool/source/offical.json"),
		path.Join(goSrc, "github.com/gogap/spirit-tool/source/third_party.json"),
	}

	sources = append(sources, extSources...)

	var pkgs []string
	if pkgs, err = urnsToPackages(urns, sources...); err != nil {
		return
	}

	getPkg := context.Bool("get")
	updatePkg := context.Bool("update")

	if getPkg {
		if updatePkg {
			spirit.Logger().Infof("downloading %d packages (go get -u) ...", len(pkgs))
		} else {
			spirit.Logger().Info("downloading %d packages (go get) ...", len(pkgs))
		}

		getPackages(updatePkg, pkgs...)
		spirit.Logger().Info("packages download finished.")
	}

	strArgs := context.StringSlice("args")
	tmplArgs := map[string]string{}

	for _, arg := range strArgs {
		arg = strings.TrimSpace(arg)
		if arg != "" {
			v := strings.Split(arg, "=")
			if len(v) != 2 {
				err = fmt.Errorf("the args to template format error, arg: %s", arg)
				return
			}
			tmplArgs[v[0]] = v[1]
		}
	}

	data := map[string]interface{}{
		"packages":    pkgs,
		"config":      string(confData),
		"innerConfig": innerConfig,
		"configFile":  configFileName,
		"args":        tmplArgs,
	}

	buffer := &bytes.Buffer{}
	if err = tmpl.Execute(buffer, data); err != nil {
		return
	}

	spirit.Logger().Infof("new app create at %s", appPath)

	if err = os.MkdirAll(appPath, os.FileMode(0755)); err != nil {
		return
	}

	srcPath := path.Join(appPath, "main.go")
	if err = ioutil.WriteFile(srcPath, buffer.Bytes(), os.FileMode(0644)); err != nil {
		return
	}

	spirit.Logger().Infof("create src: %s", srcPath)

	if !innerConfig {
		confPath := path.Join(appPath, configFileName)
		if err = ioutil.WriteFile(confPath, confData, os.FileMode(0644)); err != nil {
			return
		}
		spirit.Logger().Infof("create config: %s", confPath)
	}

	return
}

func parseSpiritConfigUsingURN(conf spirit.SpiritConfig) (urns []string, err error) {
	if urns = parseActorsUsingURN(
		conf.InputTranslators,
		conf.OutputTranslators,
		conf.Inboxes,
		conf.Outboxes,
		conf.Receivers,
		conf.Senders,
		conf.Routers,
		conf.Components,
		conf.LabelMatchers,
		conf.URNRewriters,
	); err != nil {
		return
	}

	for _, readerPool := range conf.ReaderPools {
		urns = append(urns, parseActorUsingURN(readerPool.ActorConfig, *readerPool.Reader)...)
	}

	for _, writerPool := range conf.WriterPools {
		urns = append(urns, parseActorUsingURN(writerPool.ActorConfig, *writerPool.Writer)...)
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

func urnsToPackages(urns []string, sourceFiles ...string) (packages []string, err error) {
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
			err = fmt.Errorf("urn of %s not exist", urn)
		} else {
			pkgs[pkg] = true
		}
	}

	for pkg, _ := range pkgs {
		packages = append(packages, pkg)
	}

	sort.Sort(sort.StringSlice(packages))

	return
}

func getPackages(update bool, pkgs ...string) {
	wg := sync.WaitGroup{}

	cmd := "go get -v "
	if update {
		cmd = "go get -u -v "
	}

	for _, pkg := range pkgs {
		getCmd := cmd + pkg
		wg.Add(1)
		go execCommand(getCmd, &wg)
	}
	wg.Wait()
}
