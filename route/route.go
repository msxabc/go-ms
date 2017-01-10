package route

import (
	"gt-go-ms/config"
	"gt-go-ms/log"
	"errors"
	"encoding/json"
	"sync"
	"io/ioutil"
	"strings"
	//"fmt"
)


var services map[string]string
var once sync.Once //each ms only have one copy of routing rules
var reloadSig chan bool


func Get (serviceName string) (string, error) {
	if services == nil {
		err := new()
		if (err != nil) {return "", err}
	}

	v, ok := services[serviceName]
	if (!ok) {return "", errors.New("Unable to find route for " + serviceName)}

	return v, nil
}

func new() error{
	var err error
	once.Do(func() {
		reloadSig = make(chan bool)

		err = load()
		if (err == nil){
			
			config.CallMeWhenReload(reloadSig)

			go func() {
				for {
			        r := <-reloadSig
			        if (r == true) {
			        	e := load()
				        if (e != nil){
				        	log.Get().Err("Unable to reload routing rules.  Routing rules unchanged")
				        }	
				    }
			    }
			}()
		}
	})

	return err
}

func load() error{
	c := config.Get()
	switch(strings.ToLower(c.RouteType)){
	case "file":
		return loadFromFile(c.RouteFile)
	default:
		return errors.New("No route rules defined")
	}

	return nil
}


func loadFromFile(fileName string) error{
	if (fileName == "") {return errors.New("No route file defined")}
	data, err := ioutil.ReadFile(fileName) 
	if (err != nil){
		return err
	}
	routes := make(map[string]string)
	err = json.Unmarshal(data, &routes)
	
	if (err != nil){
		return err
	}
	//reset service route map, if there is a previous one it will be garbage collected
	services = routes

	return nil
}

func GetAll() string {
	var s string = ""
	for k, v := range services{
		s = s + "\n" + k + ":" + v + ","
	}

	return s
}