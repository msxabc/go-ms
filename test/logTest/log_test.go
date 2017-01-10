package logTest

import (
	"gt-go-ms/log"
	"testing"
)



func TestLog(t *testing.T){
	t.Log("Starting log test") 
	l := log.Get()
	if (l == nil){
		t.Error("Error happened, returned nil log")
		t.FailNow()
	}
	defer l.Close()

	
	l.Emerg("Emerg")
	l.Alert("Alert")
	l.Err("Err")
	l.Warn("Warn")
	l.Info("Info")
	l.Debug("Debug")
	l.Crit("Crit")
	
}