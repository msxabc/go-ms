package commtest

import (
	"gt-go-ms/comm"
	"gt-go-ms/route"
	"gt-go-ms/config"
	"gt-go-ms/microservice"
	"testing"
	"time"
	"strings"
)

var addData = []byte(`{"numbers":[{"number":2}, {"number":6}]}`)

func TestSetup(t *testing.T){
	_, err := config.New("./testConfig.json")
	//give enough time to setup SIGHUP handling
	time.Sleep(2000 * time.Millisecond)

	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	ms, err := microservice.New()
	if (err != nil){
		t.Error(err)
		t.FailNow()
	}

	ms.AddEndpoint(MakeAddEndpoint(), "add", DecodeAddRequest(), EncodeAddResponse())
	go ms.Start(8080)

	//give enough time for the service to start
	time.Sleep(2000*time.Millisecond)
}

func TestSend(t *testing.T){
	t.Log("undefined route test")
	resp:= comm.Send("undefinedRoute", addData)
	if (resp.Err == nil) {
		t.Error("Expected error for bad route but got response")
	}


	t.Log("defined ms route test")
	resp = comm.Send("Add", addData)

	if (resp.Err != nil) {
		t.Error(resp.Err)
		t.Log(route.GetAll())
	}

	s := string(resp.Data[:])
	if (strings.Index(s, `"sum":8`) < 0){
		t.Error("Send message to add service failed with return data: " + s)
	}

	t.Log("defined route returns error http status test")
	resp = comm.Send("Httperror", addData)
	if (resp.Err == nil) {
		t.Error("Expected error for bad route but got response")
	}
}

func TestSendAsync(t *testing.T){
	

	undefinedResp := make(chan comm.Response)
	addResp := make (chan comm.Response)
	errorResp := make (chan comm.Response)

	comm.SendAsync("undefinedRoute", addData, undefinedResp)
	comm.SendAsync("Add", addData, addResp)
	comm.SendAsync("Httperror", addData, errorResp)

	t.Log("undefined route test")
	resp := <- undefinedResp
	if (resp.Err == nil) {
		t.Error("Expected error for bad route but got response")
	}


	t.Log("defined ms route test")
	resp = <- addResp
	if (resp.Err != nil) {
		t.Error(resp.Err)
		t.Log(route.GetAll())
	}

	s := string(resp.Data[:])
	if (strings.Index(s, `"sum":8`) < 0){
		t.Error("Send message to add service failed with return data: " + s)
	}

	t.Log("defined route returns error http status test")
	resp = <- errorResp
	if (resp.Err == nil) {
		t.Error("Expected error for bad route but got response")
	}
}



