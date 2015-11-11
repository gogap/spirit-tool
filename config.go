package main

type URNPackage struct {
	URN string `json:"urn"`
	Pkg string `json:"pkg"`
}

type SourceConfig struct {
	UpdateTime string       `json:"update_time"`
	Packages   []URNPackage `json:"packages"`
}
