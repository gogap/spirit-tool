package main

import (
	"path"

	"github.com/gogap/spirit"
)

type Package struct {
	gosrc    string
	URI      string
	Revision string
}

func (p *Package) Get(update bool) (err error) {

	baseCMD := "go get "
	if verbosity > 0 {
		baseCMD = "go get -v "
	}

	cmd := baseCMD + p.URI
	if update {
		cmd = baseCMD + "-u " + p.URI
	}

	var out []byte

	if out, err = execCommand(cmd); err != nil {
		spirit.Logger().Errorln(string(out))
		return
	}

	if p.Revision == "" {
		return
	}

	checkoutCMD := "git -C " + path.Join(p.gosrc, p.URI) + " checkout " + p.Revision

	if out, err = execCommand(checkoutCMD); err != nil {
		spirit.Logger().Errorln(string(out))
		return
	}

	return
}
