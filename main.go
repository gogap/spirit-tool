package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

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
		commandRun(run),
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

	createOpts := CreateOptions{
		TemplateName:     templateName,
		GoPath:           goPath,
		ProjectPath:      tmpDir,
		GetPackages:      true,
		UpdatePackages:   updatePkg,
		ForceWrite:       true,
		Sources:          sources,
		PackagesRevision: nil,
	}

	if err = helper.RunProject(createOpts, tmplArgs); err != nil {
		return
	}

	return
}
