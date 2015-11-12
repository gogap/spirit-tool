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
//<-range $_, $pkg := .packages->////<-printf "\t_ \"%s\"\n" $pkg.URI->////<-end->////<-printf ")"->//


const TemplateVersion = "0.0.1"

const CreateTime string //<-printf "= \"%s\"" .create_time->//

var configFile string //<-printf "= \"%s\"" .config_filename->//

var (
	innerConfig bool //<-printf "= %v" .args.inner_config->//
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
	config string //<-if .args.inner_config->////<-printf "= `%s`" .config->////<-end->//
)
