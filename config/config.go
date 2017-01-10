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
	RouteType string `json:"routeType"`
	RouteFile string  `json:"routeFile"`
	PropertyFile string `json:"propertyFile"`
}

type logConfig struct {
	Facility string `json:"facility"` 
	Level string `json:"level"`
}

var appProperties map[string]interface{}


var c *Config
var once sync.Once //each MS should have only one copy of config
var configFile string = "/etc/gt-agent.json"
var reloadChannels map[chan bool] interface{}


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
		
		reloadChannels = make(map[chan bool]interface{})

		go func() {
			for {
		        sig := <-sigs
		        if (sig == syscall.SIGHUP) {
		        	e := loadConfig()
			        if (e != nil){
			        	log.Println("Error: Unable to reload config parameters. Application will continue ot use current configuration parameters.  You can try to restart this application to apply new changes")
			        	log.Println(e)
			        }	
			        for c, _ := range reloadChannels {
			        	c<-true
			        }
			    }
		    }
		}()

    })

    if filename != configFile && err == nil {
    	//possibile reusing of New function with a different file
    	err = errors.New("Unable to create new configuraiton from " + filename + ", using current configuration created from : " + configFile)
    }

	return c, err
}

func Get() *Config{
	return c
}

func GetAppConfig(name string) (interface{}, error) {
	if v, ok := appProperties[name]; !ok{
		return "", errors.New("No config key defined for " + name)
	}else{
		return v, nil
	}

}

func CallMeWhenReload(c chan bool){
	reloadChannels[c] = nil
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

	if c.PropertyFile != ""{
		appProperties = make(map[string]interface{})
		data, err = ioutil.ReadFile(c.PropertyFile) 
		if (err != nil){
			return err
		}

		err = json.Unmarshal(data, &appProperties)
	
		if (err != nil){
			return err
		}
	}


	return nil
}