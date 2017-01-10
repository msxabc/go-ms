package log

import (
	"log/syslog"
	"gt-go-ms/config"
	"os"
	"sync"
	"strings"
)

type logger struct {
 logWriter *syslog.Writer
 mutex *sync.Mutex
}

var l *logger
var once sync.Once

type Logger interface{
	Crit(s string) error
	Emerg(s string)error
	Alert(s string) error
	Err(s string) error
	Warn(s string) error
	Info(s string) error
	Debug(s string) error
	Close()
}


func new() (Logger){
	
	once.Do(func() {
		w, e := syslog.New(getFacility() | getSeverity(), os.Args[0])
		m := &sync.Mutex{}
		if e == nil {
        	l = &logger{w, m}
        }else {
        	panic("Service logger creation failed:" + e.Error())
        }
    })
	return l
}

func Get() (Logger){
	if (l == nil) {new()}
	return l
}

func (l *logger) Crit(s string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if (getSeverity() >= syslog.LOG_CRIT) { return l.logWriter.Crit(s) }
	return nil 
}

func (l *logger) Emerg(s string) error{
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if (getSeverity() >= syslog.LOG_EMERG) {return l.logWriter.Emerg(s) }
	return nil 
}

func (l *logger) Alert(s string) error{
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if (getSeverity() >= syslog.LOG_ALERT) {return l.logWriter.Alert(s)}
	return nil 
}

func (l *logger) Err (s string) error{
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if (getSeverity() >= syslog.LOG_ERR) {return l.logWriter.Err(s)}
	return nil 
}

func (l *logger) Warn (s string) error{
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if (getSeverity() >= syslog.LOG_WARNING) {return l.logWriter.Warning(s)}
	return nil 
}

func (l *logger) Info (s string) error{
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if (getSeverity() >= syslog.LOG_INFO) {return l.logWriter.Info(s)}
	return nil 
}

func (l *logger) Debug (s string) error{
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if (getSeverity() >= syslog.LOG_DEBUG) {return l.logWriter.Debug(s)}
	return nil 
}

func (l *logger) Close(){
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logWriter.Close()
}


func getFacility() syslog.Priority{
	var f syslog.Priority
	c := config.Get()

	if c == nil { return syslog.LOG_KERN}

	switch strings.ToLower(c.Log.Facility) {
		case "kern": f = syslog.LOG_KERN
		case "user": f = syslog.LOG_USER
		case "daemon": f = syslog.LOG_DAEMON
		case "auth": f = syslog.LOG_AUTH
		case "syslog": f = syslog.LOG_SYSLOG
		case "lpr": f = syslog.LOG_LPR
		case "news": f = syslog.LOG_NEWS
		case "uucp": f = syslog.LOG_UUCP
		case "cron": f = syslog.LOG_CRON
		case "authpriv": f = syslog.LOG_AUTHPRIV
		case "ftp": f = syslog.LOG_FTP
		case "local0": f = syslog.LOG_LOCAL0
		case "local1": f = syslog.LOG_LOCAL1
		case "local2": f = syslog.LOG_LOCAL2
		case "local3": f = syslog.LOG_LOCAL3
		case "local4": f = syslog.LOG_LOCAL4
		case "local5": f = syslog.LOG_LOCAL5
		case "local6": f = syslog.LOG_LOCAL6
		case "local7": f = syslog.LOG_LOCAL7
		default: f = syslog.LOG_KERN
	}

	return f
}

func getSeverity () syslog.Priority{
	var s syslog.Priority
	c := config.Get()

	if c == nil { return syslog.LOG_INFO}

	switch strings.ToLower(c.Log.Level){
		case "emerg": s = syslog.LOG_EMERG
		case "alert": s = syslog.LOG_ALERT
		case "crit": s = syslog.LOG_CRIT
		case "err": s = syslog.LOG_ERR
		case "notice": s = syslog.LOG_NOTICE
		case "info": s = syslog.LOG_INFO
		case "debug": s = syslog.LOG_DEBUG
		default: s = syslog.LOG_INFO
	}

	return s
}


