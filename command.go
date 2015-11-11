package main

import (
	"os"

	"github.com/codegangsta/cli"
)

type cliAction func(context *cli.Context)

func commandUpdate(action cliAction) cli.Command {
	return cli.Command{
		Name:      "update",
		ShortName: "",
		Usage:     "Update spirit-tool and reinstall",
		Action:    action,
	}
}

func commandFind(action cliAction) cli.Command {
	return cli.Command{
		Name:      "find",
		ShortName: "",
		Usage:     "Find package by urn",
		Action:    action,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "urn, u",
				Value: new(cli.StringSlice),
				Usage: "urn",
			},
		},
	}
}

func commandCreate(action cliAction) cli.Command {
	return cli.Command{
		Name:      "create",
		ShortName: "",
		Usage:     "Create src from template, it will create into $GOPATH/{PATH}",
		Action:    action,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "gopath",
				Value: os.Getenv("GOPATH"),
				Usage: "default gopath is get from $GOPATH",
			}, cli.StringFlag{
				Name:  "config, c",
				Value: "",
				Usage: "config file",
			}, cli.BoolFlag{
				Name:  "inner, i",
				Usage: "embed config content",
			}, cli.StringFlag{
				Name:  "path, p",
				Value: "",
				Usage: "",
			}, cli.StringFlag{
				Name:  "template,t",
				Value: "classic",
				Usage: "which template to use, default is classic",
			}, cli.BoolFlag{
				Name:  "get, g",
				Usage: "automatic get packages by `go get` command",
			}, cli.BoolFlag{
				Name:  "update, u",
				Usage: "if get flag is ture, it will use `go get -u`",
			}, cli.StringSliceFlag{
				Name:  "source, s",
				Usage: "your own source file",
			}, cli.StringSliceFlag{
				Name:  "args, a",
				Usage: "the args will pass into template, format: -a key=val, you could use `args.key` to get value",
			},
		},
	}
}
