package mstest

import (
	"gt-go-ms/microservice"
	"testing"
	"net/http"
	"io/ioutil"
	"strings"
	"bytes"
	"time"
)

var echoData = []byte(`{"data":"abc123", "command":"echo"}`)
var doubleData = []byte(`{"data":"abc123", "command":"double"}`)
var badData = []byte(`{"data":"abc123", "command":"google"}`)

func TestEchoMS(t * testing.T){

	go microservice.New(MakeEchoEndpoint(),
		"echo",
		DecodeRequest(),
		EncodeResponse(),
		8080)
	//give enough time for the server to run
	time.Sleep(2000 * time.Millisecond)


	t.Log("Test: echo command")
	req, err := http.NewRequest("POST", "http://localhost:8080/echo", bytes.NewBuffer(echoData))
	client := &http.Client{}
    resp, err := client.Do(req)

	if (err != nil){
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	s := string(body[:])
	if (strings.Index(s, `"v":"abc123"`) < 0){
		t.Error("echo command failed: " + s)
	}
}

func TestDoubleEchoMS(t * testing.T){


	t.Log("Test: double echo command")
	req, err := http.NewRequest("POST", "http://localhost:8080/echo", bytes.NewBuffer(doubleData))
	client := &http.Client{}
    resp, err := client.Do(req)

	if (err != nil){
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	s := string(body[:])
	if (strings.Index(s, `"v":"abc123abc123"`) < 0){
		t.Error("double command failed: " + s)
	}
}

func TestBadEchoMS(t * testing.T){

	t.Log("Test: double echo command")
	req, err := http.NewRequest("POST", "http://localhost:8080/echo", bytes.NewBuffer(badData))
	client := &http.Client{}
    resp, err := client.Do(req)

	if (err != nil){
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	s := string(body[:])
	if (strings.Index(s, `"err":"Unknown command"`) < 0){
		t.Error("bad command failed: " + s)
	}
}


