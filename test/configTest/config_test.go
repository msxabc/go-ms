package configTest

import (
	"gt-go-ms/config"
	"io/ioutil"
	"testing"
	"os"
	"syscall"
	"time"
)

var goodData string = `{
				"log": 
						{
							"facility" : "kern", 
							"level" : "debug"
						}
			 }`


var badData string = `{"abc": "123"}`



func TestDefaultConfig(t * testing.T){
	data := badData
	err := ioutil.WriteFile("./testConfig.json", []byte(data), os.ModePerm)

	c, err := config.New("./testConfig.json")
	//give enough time to setup SIGHUP handling
	time.Sleep(2000 * time.Millisecond)

	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	if (c.Log == nil){
		t.Error("default failed, got nil config")
		t.FailNow()
	}

	if (c.Log.Facility != "local0"){
		t.Error("Expected facility local0, but got " + c.Log.Facility)
	}

	if (c.Log.Level != "info"){
		t.Error("Expected facility info, but got " + c.Log.Level)
	}

}

func TestGoodConfig(t *testing.T){
	data := goodData

	err := ioutil.WriteFile("./testConfig.json", []byte(data), os.ModePerm)
	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	c, err := config.New("./testConfig.json")


	if (c.Log.Facility != "local0"){
		t.Error("Expected facility local0, but got " + c.Log.Facility)
	}

	if (c.Log.Level != "info"){
		t.Error("Expected facility info, but got " + c.Log.Level)
	}

	t.Log("Signal to reload good configuration")
	sendSIGHUP()

	if (c.Log == nil){
		t.Error("Reload failed, got nil config")
		t.FailNow()
	}

	if (c.Log.Facility != "kern"){
		t.Error("Expected facility kern, but got " + c.Log.Facility)
	}

	if (c.Log.Level != "debug"){
		t.Error("Expected facility debug, but got " + c.Log.Level)
	}
}

func TestBadConfig(t *testing.T){
	data := badData

	err := ioutil.WriteFile("./testConfig.json", []byte(data), os.ModePerm)
	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	c, err := config.New("./testConfig.json")

	if (c.Log.Facility != "kern"){
		t.Error("Expected facility kern, but got " + c.Log.Facility)
	}

	if (c.Log.Level != "debug"){
		t.Error("Expected facility debug, but got " + c.Log.Level)
	}

	t.Log("Signal to reload bad configuration")
	sendSIGHUP()


	c, err = config.New("./testConfig.json")

	if (err != nil){
		t.Error(err)
	}

	//config should contain previous value
	if (c.Log.Facility != "kern"){
		t.Error("Expected facility kern, but got " + c.Log.Facility)
	}

	if (c.Log.Level != "debug"){
		t.Error("Expected facility debug, but got " + c.Log.Level)
	}

}

func TestRenew(t *testing.T){

	data := goodData

	err := ioutil.WriteFile("./testConfigABC.json", []byte(data), os.ModePerm)
	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	c, err := config.New("./testConfigABC.json")

	if (err == nil){
		t.Error("Expected error for calling New again")
	}

	//config should contain previous value
	if (c.Log.Facility != "kern"){
		t.Error("Expected facility kern remain unchanged, but got " + c.Log.Facility)
	}

	if (c.Log.Level != "debug"){
		t.Error("Expected facility remain unchanged, but got " + c.Log.Level)
	}


}

var appConfig string = `{"propertyFile": "./properties.json"}`
var appData string = `{"string": "123", "int": 1, "boolean": true}`

func TestAppConfig(t *testing.T){

	err := ioutil.WriteFile("./testConfig.json", []byte(appConfig), os.ModePerm)
	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	err = ioutil.WriteFile("./properties.json", []byte(appData), os.ModePerm)
	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	t.Log("Signal to reload app configuration")
	sendSIGHUP()

	v, _:= config.GetAppConfig("string")

	if s, ok := v.(string); !ok {
		t.Error("Failed string test, expected string type return")
		t.Log(v)
	}else{
		if s != "123" {
			t.Error("Failed string test, expected 123 got " + s)
		}
	}

	v, _ = config.GetAppConfig("int")

	//golang json marshals number to float64 by default
	if s, ok := v.(float64); !ok {
		t.Error("Failed int test, expected int type return")
	}else{
		if s != 1 {
			t.Error("Failed int test, expected 1")
			t.Log(s)
		}
	}

	v, _ = config.GetAppConfig("boolean")

	if s, ok := v.(bool); !ok {
		t.Error("Failed boolean test, expected boolean type return")
		t.Log(v)
	}else{
		if s != true {
			t.Error("Failed boolean test, expected true got false")
		}
	}

	v, err = config.GetAppConfig("undefined")

	if (err == nil){
		t.Error("Failed undefined test, expected error got data")
		t.Log(v)
	}

}


func TestReload(t *testing.T){
	data := badData

	err := ioutil.WriteFile("./testConfig.json", []byte(data), os.ModePerm)
	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	config.New("./testConfig.json")
	ch := make(chan bool)
	config.CallMeWhenReload(ch)
	sendSIGHUP()
	rl := <-ch
	if (true != rl) {
		t.Error("Expected reload signal true, but got false")
	}
}



func sendSIGHUP(){
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	time.Sleep(2000 * time.Millisecond)
}