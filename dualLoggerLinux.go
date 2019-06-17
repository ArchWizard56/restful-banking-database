// +build linux

package main

import (
	//"fmt"
	"log"
    "log/syslog"
    "os"
    //"time"
)

var SystemSyslog *syslog.Writer

func SetupLogging () {
    var err error
    SystemSyslog, err = syslog.New(syslog.LOG_INFO,os.Args[0])
    if err != nil {
    log.Fatal(err)
    }
}

func DualDebug (m string) {
    SystemSyslog.Debug(m)
    log.Print(m)
}
func DualInfo (m string) {
    SystemSyslog.Write([]byte(m))
    log.Print(m)
}
func DualNotice (m string) {
    SystemSyslog.Notice(m)
    log.Print(m)
}
func DualWarning (m string) {
    SystemSyslog.Warning(m)
    log.Print(m)
}
func DualErr (err error) {
    SystemSyslog.Err(err.Error())
    log.Fatal(err)
}
