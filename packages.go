package main

import (
	"path"
)

type Package struct {
	gosrc    string
	URI      string
	Revision string
}

func (p *Package) Get(update bool) (err error) {
	cmd := "go get -v " + p.URI
	if update {
		cmd = "go get -u -v " + p.URI
	}

	if _, err = execCommand(cmd); err != nil {
		return
	}

	if p.Revision == "" {
		return
	}

	checkoutCMD := "git -C " + path.Join(p.gosrc, p.URI) + " checkout " + p.Revision

	if _, err = execCommand(checkoutCMD); err != nil {
		return
	}

	return
}
