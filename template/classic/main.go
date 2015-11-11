package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/gogap/spirit"
)

//<-printf "import ("->//
//<-range $_, $pkg := .packages->////<-printf "\t_ \"%s\"\n" $pkg->////<-end->////<-printf ")"->//

var configFile string //<-printf "= \"%s\"" .configFile->//

var (
	innerConfig bool //<-printf "= %v" .innerConfig->//
)

func main() {
	var err error
	defer func() {
		if err != nil {
			spirit.Logger().Error(err)
			os.Exit(128)
		}
	}()

	if configFile != "" {
		var fileData []byte
		if fileData, err = ioutil.ReadFile(configFile); err != nil {
			return
		}

		config = string(fileData)
	}

	if config == "" {
		err = fmt.Errorf("config is empty")
		return
	}

	spiritConf := spirit.SpiritConfig{}

	if err = json.Unmarshal([]byte(config), &spiritConf); err != nil {
		return
	}

	if err = spiritConf.Validate(); err != nil {
		err = fmt.Errorf("spirit config validate failed, %s", err)
		return
	}

	var sp spirit.Spirit
	if sp, err = spirit.NewClassicSpirit(); err != nil {
		err = fmt.Errorf("create new classic spirit error, %s", err)
		return
	}

	if err = sp.Build(spiritConf); err != nil {
		err = fmt.Errorf("build classic spirit error, %s", err)
		return
	}

	var wg *sync.WaitGroup
	if wg, err = sp.Run(); err != nil {
		return
	}

	wg.Wait()
}

var (
	config string //<-if .innerConfig->////<-printf "= `%s`" .config->////<-end->//
)
