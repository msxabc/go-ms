package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const defaults string = `{
	"log": {
		"facility" : "local0",
		"level": "info"
	}
}`

type Config struct {
	Log *logConfig `json:"log"`
}

type logConfig struct {
	Facility string `json:"facility"` 
	Level string `json:"level"`
}

var c *Config
var once sync.Once
var configFile string = "/etc/gt-go-ms.json"


func New(filename string) (*Config, error){
	var err error = nil
	once.Do(func() {
		if (filename != ""){
			configFile = filename
		}
		c = &Config{}
		err = json.Unmarshal([]byte(defaults), c)
	
		if (err != nil){
			return 
		}

		e := loadConfig()
		if (e != nil) {
			err = errors.New("Unable to init config from: " + configFile + e.Error())
			return
		}

		//setup reload logic based on HUP signal
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP)
		

		go func() {
			for {
		        sig := <-sigs
		        if (sig == syscall.SIGHUP) {
		        	e := loadConfig()
			        if (e != nil){
			        	log.Println("Error: Unable to reload config parameters. Application will continue ot use current configuration parameters.  You can try to restart this application to apply new changes")
			        }	
			    }
		    }
		}()

    })

	return c, err
}

func Get() *Config {
	return c
}


func loadConfig() (error){
	data, err := ioutil.ReadFile(configFile) 
	if (err != nil){
		return err
	}
	err = json.Unmarshal(data, c)
	
	if (err != nil){
		return err
	}


	return nil
}