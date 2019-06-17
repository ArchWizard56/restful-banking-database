// +build linux

package main

import (
	"fmt"
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
    m = fmt.Sprintf("DEBUG: %s", m)
    SystemSyslog.Debug(m)
    log.Print(m)
}
func DualInfo (m string) {
    m = fmt.Sprintf("INFO: %s", m)
    SystemSyslog.Write([]byte(m))
    log.Print(m)
}
func DualNotice (m string) {
    m = fmt.Sprintf("NOTICE: %s", m)
    SystemSyslog.Notice(m)
    log.Print(m)
}
func DualWarning (m string) {
    m = fmt.Sprintf("WARNING: %s", m)
    SystemSyslog.Warning(m)
    log.Print(m)
}
func DualErr (err error) {
    m := fmt.Sprintf("FATAL: %s", err.Error())
    SystemSyslog.Err(m)
    log.Fatal(m)
}
