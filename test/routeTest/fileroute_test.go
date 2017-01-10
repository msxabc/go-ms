package routeTest

import (
	"gt-go-ms/route"
	"gt-go-ms/config"
	"testing"
	"time"
	"syscall"
	"os"
	"io/ioutil"
)


var goodRoute1Config string = `{
				"routeType": "file",
				"routeFile": "./route1.json"
			 }`

var goodRoute2Config string = `{
				"routeType": "file",
				"routeFile": "./route2.json"
			 }`

var badRouteConfig string = `{
				"routeType": "file",
				"routeFile": "./badroute.json"
			 }`



func TestGoodRoute(t *testing.T){
	setup(t, goodRoute1Config)

	if r, err := route.Get("abc"); r != "www.abc.com"{
		t.Error("Expected abc route, got " + r)
		if err != nil {
			t.Error (err)
		}
	}

	if r, err := route.Get("123"); err == nil {
		t.Error("Expected error for empty route, but got " + r)
	}

}

func TestRouteReload(t *testing.T){
	c := make (chan bool)
	config.CallMeWhenReload(c)

	if r, err := route.Get("abc"); r != "www.abc.com"{
		t.Error("Expected abc route, got " + r)
		if err != nil {
			t.Error (err)
		}
	}

	if r, err := route.Get("123"); err == nil {
		t.Error("Expected error for empty route, but got " + r)
	}

	setup(t, goodRoute2Config)
	sendSIGHUP()

	if r, err := route.Get("abc"); err == nil{
		t.Error("Expected nil abc route after reload, but got " + r)
	}

	if r, _ := route.Get("123"); r != "abc123" {
		t.Error("Expected good route for 123 after reload, but got " + r)
	}
	
}

func TestNonExistRoute(t *testing.T){
	setup(t, badRouteConfig)
	sendSIGHUP()

	//things should remain the same as previous test
	if r, err:= route.Get("abc"); err == nil{
		t.Error("Expected nil abc route after reload, but got " + r)
	}

	if r,_ := route.Get("123"); r != "abc123" {
		t.Error("Expected good route for 123 after reload, but got " + r)
	}
}




func setup(t *testing.T, configData string){
	err := ioutil.WriteFile("./testConfig.json", []byte(configData), os.ModePerm)
	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	_, err = config.New("./testConfig.json")
	//give enough time to setup SIGHUP handling
	time.Sleep(2000 * time.Millisecond)

	if (err != nil){
		t.Error(err)
		t.FailNow()
	}
}

func sendSIGHUP(){
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	time.Sleep(2000 * time.Millisecond)
}