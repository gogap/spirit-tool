package main

import (
	"os"

	"github.com/codegangsta/cli"
)

type cliAction func(context *cli.Context)

func commandUpgrade(action cliAction) cli.Command {
	return cli.Command{
		Name:      "upgrade",
		ShortName: "",
		Usage:     "Upgrade spirit-tool and reinstall",
		Action:    action,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "verbosity, v",
				Usage: "How much troubleshooting info to print (1~5)",
			},
		},
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
			}, cli.BoolFlag{
				Name:  "force, f",
				Usage: "is your app is exist, it will overwrite it",
			}, cli.StringFlag{
				Name:  "rev, r",
				Usage: "packages revision config filepath, json format, e.g.: {\"github.com/gogap/spirit\":\"master\"}",
			}, cli.IntFlag{
				Name:  "verbosity, v",
				Usage: "How much troubleshooting info to print (1~5)",
			},
		},
	}
}

func commandRun(action cliAction) cli.Command {
	return cli.Command{
		Name:      "run",
		ShortName: "",
		Usage:     "run your spirit config directly",
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
			}, cli.StringFlag{
				Name:  "template,t",
				Value: "classic",
				Usage: "which template to use, default is classic",
			}, cli.BoolFlag{
				Name:  "update, u",
				Usage: "run `go get -u` before run",
			}, cli.StringSliceFlag{
				Name:  "source, s",
				Usage: "your own source file",
			}, cli.StringSliceFlag{
				Name:  "args, a",
				Usage: "the args will pass into template, format: -a key=val, you could use `args.key` to get value",
			}, cli.StringFlag{
				Name:  "rev, r",
				Usage: "packages revision config filepath, json format, e.g.: {\"github.com/gogap/spirit\":\"master\"}",
			}, cli.BoolFlag{
				Name:  "detach, d",
				Usage: "Run spirit in background and print PID",
			}, cli.StringSliceFlag{
				Name:  "env, e",
				Usage: "Set environment variables",
			}, cli.IntFlag{
				Name:  "verbosity, v",
				Usage: "How much troubleshooting info to print (1~5)",
			},
		},
	}
}

func commandBuild(action cliAction) cli.Command {
	return cli.Command{
		Name:      "build",
		ShortName: "",
		Usage:     "build your spirit with config",
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
			}, cli.StringFlag{
				Name:  "template,t",
				Value: "classic",
				Usage: "which template to use, default is classic",
			}, cli.BoolFlag{
				Name:  "update, u",
				Usage: "run `go get -u` before build",
			}, cli.StringSliceFlag{
				Name:  "source, s",
				Usage: "your own source file",
			}, cli.StringSliceFlag{
				Name:  "args, a",
				Usage: "the args will pass into template, format: -a key=val, you could use `args.key` to get value",
			}, cli.StringFlag{
				Name:  "rev, r",
				Usage: "packages revision config filepath, json format, e.g.: {\"github.com/gogap/spirit\":\"master\"}",
			}, cli.StringFlag{
				Name:  "output, o",
				Usage: "the binary output path",
			}, cli.IntFlag{
				Name:  "verbosity, v",
				Usage: "How much troubleshooting info to print (1~5)",
			},
		},
	}
}
