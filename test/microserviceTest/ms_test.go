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



func TestMSSetup(t * testing.T){
	ms, err := microservice.New()
	if (err != nil){
		t.Error(err)
		t.FailNow()
	}
	
	ms.AddEndpoint(MakeEchoEndpoint(), "echo", DecodeEchoRequest(), EncodeEchoResponse())
	ms.AddEndpoint(MakeAddEndpoint(), "add", DecodeAddRequest(), EncodeAddResponse())
	go ms.Start(8080)

	//give enough time for the service to start
	time.Sleep(2000*time.Millisecond)
}

func TestEchoEndpoint(t * testing.T){
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

func TestDoubleEchoEndpoint(t * testing.T){
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

func TestBadEchoEndpoint(t * testing.T){
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


var addData = []byte(`{"numbers":[{"number":2}, {"number":6}]}`)
var addBadData = []byte(`{"numbers":[{"number":"abc123"}, {"number":6}]}`)

func TestAddEndpoint(t *testing.T){
	t.Log("Test: add endpoint")
	req, err := http.NewRequest("POST", "http://localhost:8080/add", bytes.NewBuffer(addData))
	client := &http.Client{}
    resp, err := client.Do(req)

	if (err != nil){
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	s := string(body[:])
	if (strings.Index(s, `"sum":8`) < 0){
		t.Error("add failed with return data: " + s)
	}
}

func TestBadAddEndpoint(t *testing.T){
	t.Log("Test: add endpoint bad data")
	req, err := http.NewRequest("POST", "http://localhost:8080/add", bytes.NewBuffer(addBadData))
	client := &http.Client{}
    resp, err := client.Do(req)

	if (err != nil){
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	s := string(body[:])
	if (strings.Index(s, `"sum":`) >= 0){
		t.Error("add bad data failed with return data: " + s)
	}

}

