package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gogap/spirit"
)

var (
	app *cli.App
)

var (
	verbosity = 0
)

func main() {
	app = cli.NewApp()
	app.Name = "spirit-tool"
	app.Authors = []cli.Author{{"zeal", "xujinzheng@gmail.com"}}
	app.Usage = "help user easily to use spirit"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		commandUpgrade(upgrade),
		//commandFind(find),
		commandRun(run),
		commandCreate(create),
		commandBuild(build),
	}

	app.Run(os.Args)
}

func upgrade(context *cli.Context) {
	var err error

	defer func() {
		if err != nil {
			spirit.Logger().Error(err)
			os.Exit(128)
		}
	}()

	var out []byte
	cmd := "go get -u github.com/gogap/spirit-tool"
	if out, err = execCommand(cmd); err != nil {
		spirit.Logger().Errorln(err)
		return
	}
	spirit.Logger().Infoln(out)

	cmd = "go install github.com/gogap/spirit-tool"
	if out, err = execCommand(cmd); err != nil {
		spirit.Logger().Errorln(err)
		return
	}
	spirit.Logger().Infoln(out)

	return
}

func find(context *cli.Context) {
	return
}

func create(context *cli.Context) {
	verbosity = context.Int("verbosity")
	spirit.Logger().Level = logrus.Level(verbosity)

	var err error

	defer func() {
		if err != nil {
			spirit.Logger().Error(err)
			os.Exit(128)
		}
	}()

	goPath := context.String("gopath")
	projectPath := context.String("path")
	configFile := context.String("config")
	extSources := context.StringSlice("source")
	getPkg := context.Bool("get")
	updatePkg := context.Bool("update")
	strArgs := context.StringSlice("args")
	forceWrite := context.Bool("force")
	templateName := context.String("template")
	revConfig := context.String("rev")

	if goPath == "" {
		err = fmt.Errorf("could not get GOPATH")
		return
	}

	spirit.Logger().Infof("GOPATH: %s", goPath)

	if projectPath == "" {
		err = fmt.Errorf("please input your project path, like: github.com/your_orgs/project_name ")
		return
	}

	if configFile == "" {
		err = fmt.Errorf("please input config file")
		return
	}

	sources := []string{
		path.Join(goPath, "src", "github.com/gogap/spirit-tool/source/offical.json"),
		path.Join(goPath, "src", "github.com/gogap/spirit-tool/source/third_party.json"),
	}

	sources = append(sources, extSources...)

	tmplArgs := map[string]interface{}{}

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

	helper := SpiritHelper{}

	if err = helper.LoadSpiritConfig(configFile); err != nil {
		return
	}

	var rev map[string]string
	if revConfig != "" {
		loadKeyValueJSON(revConfig, &rev)
	}

	createOpts := CreateOptions{
		TemplateName:     templateName,
		GoPath:           goPath,
		ProjectPath:      projectPath,
		GetPackages:      getPkg,
		UpdatePackages:   updatePkg,
		ForceWrite:       forceWrite,
		Sources:          sources,
		PackagesRevision: nil,
	}

	if err = helper.CreateProject(createOpts, tmplArgs); err != nil {
		return
	}

	return
}

func run(context *cli.Context) {
	verbosity = context.Int("verbosity")
	spirit.Logger().Level = logrus.Level(verbosity)

	var err error

	defer func() {
		if err != nil {
			spirit.Logger().Error(err)
			os.Exit(128)
		}
	}()

	goPath := context.String("gopath")
	configFile := context.String("config")
	extSources := context.StringSlice("source")

	updatePkg := context.Bool("update")
	strArgs := context.StringSlice("args")
	templateName := context.String("template")
	revConfig := context.String("rev")
	detach := context.Bool("detach")
	envs := context.StringSlice("env")

	if goPath == "" {
		err = fmt.Errorf("could not get GOPATH")
		return
	}

	spirit.Logger().Infof("GOPATH: %s", goPath)

	if configFile == "" {
		err = fmt.Errorf("please input config file")
		return
	}

	sources := []string{
		path.Join(goPath, "src", "github.com/gogap/spirit-tool/source/offical.json"),
		path.Join(goPath, "src", "github.com/gogap/spirit-tool/source/third_party.json"),
	}

	sources = append(sources, extSources...)

	tmplArgs := map[string]interface{}{}

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

	helper := SpiritHelper{}

	if err = helper.LoadSpiritConfig(configFile); err != nil {
		return
	}

	tmpDir := ""
	if tmpDir, err = ioutil.TempDir("", "spirit-tool."); err != nil {
		return
	}

	var rev map[string]string
	if revConfig != "" {
		loadKeyValueJSON(revConfig, &rev)
	}

	createOpts := CreateOptions{
		TemplateName:     templateName,
		GoPath:           goPath,
		ProjectPath:      tmpDir,
		GetPackages:      true,
		UpdatePackages:   updatePkg,
		ForceWrite:       true,
		Sources:          sources,
		PackagesRevision: rev,
	}

	if err = helper.RunProject(createOpts, detach, envs, tmplArgs); err != nil {
		return
	}

	return
}

func build(context *cli.Context) {
	verbosity = context.Int("verbosity")
	spirit.Logger().Level = logrus.Level(verbosity)

	var err error

	defer func() {
		if err != nil {
			spirit.Logger().Error(err)
			os.Exit(128)
		}
	}()

	goPath := context.String("gopath")
	configFile := context.String("config")
	extSources := context.StringSlice("source")

	updatePkg := context.Bool("update")
	strArgs := context.StringSlice("args")
	templateName := context.String("template")
	revConfig := context.String("rev")
	output := context.String("output")

	if goPath == "" {
		err = fmt.Errorf("could not get GOPATH")
		return
	}

	spirit.Logger().Infof("GOPATH: %s", goPath)

	if configFile == "" {
		err = fmt.Errorf("please input config file")
		return
	}

	sources := []string{
		path.Join(goPath, "src", "github.com/gogap/spirit-tool/source/offical.json"),
		path.Join(goPath, "src", "github.com/gogap/spirit-tool/source/third_party.json"),
	}

	sources = append(sources, extSources...)

	tmplArgs := map[string]interface{}{}

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

	helper := SpiritHelper{}

	if err = helper.LoadSpiritConfig(configFile); err != nil {
		return
	}

	tmpDir := ""
	if tmpDir, err = ioutil.TempDir("", "spirit-tool."); err != nil {
		return
	}

	var rev map[string]string
	if revConfig != "" {
		loadKeyValueJSON(revConfig, &rev)
	}

	createOpts := CreateOptions{
		TemplateName:     templateName,
		GoPath:           goPath,
		ProjectPath:      tmpDir,
		GetPackages:      true,
		UpdatePackages:   updatePkg,
		ForceWrite:       true,
		Sources:          sources,
		PackagesRevision: rev,
	}

	if !path.IsAbs(output) {
		fp, _ := filepath.Abs(os.Args[0])
		output = path.Join(path.Dir(fp), output)
	}

	if err = helper.BuildProject(createOpts, output, tmplArgs); err != nil {
		return
	}

	return
}

func loadKeyValueJSON(filename string, v *map[string]string) (err error) {
	var revData []byte
	if revData, err = ioutil.ReadFile(filename); err != nil {
		return
	} else if err = json.Unmarshal(revData, v); err != nil {
		return
	}
	return
}
