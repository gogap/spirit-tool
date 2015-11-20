package main

import (
	"errors"
)

var (
	ErrGoPathIsEmpty     = errors.New("go path is empty")
	ErrProjectDirIsEmpty = errors.New("project dir is empty")
	ErrNoTemplateName    = errors.New("no template name")
)

type CreateOptions struct {
	TemplateName     string
	GoPath           string
	ProjectPath      string
	GetPackages      bool
	UpdatePackages   bool
	ForceWrite       bool
	Sources          []string
	PackagesRevision map[string]string
	IsTempPath       bool
}

func (p *CreateOptions) Validate() (err error) {
	if p.GoPath == "" {
		err = ErrGoPathIsEmpty
		return
	}

	if p.ProjectPath == "" {
		err = ErrProjectDirIsEmpty
		return
	}

	if p.TemplateName == "" {
		err = ErrNoTemplateName
		return
	}

	return
}

type ProjectOptions struct {
	IsInnerConfig     bool
	DefaultConfigName string
	Packages          []string
	Args              map[string]string
}
