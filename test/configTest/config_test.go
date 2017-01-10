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