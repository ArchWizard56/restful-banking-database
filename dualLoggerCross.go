// +build !linux

package main

import (
	"log"
    "fmt"
)

func SetupLogging () {
    DualNotice("Non-Linux systems don't support this implimentation of syslog; writing messages to Standard Output instead")
}

func DualDebug (m string) {
    m = fmt.Sprintf("DEBUG: %s", m)
    log.Print(m)
}
func DualInfo (m string) {
    m = fmt.Sprintf("INFO: %s", m)
    log.Print(m)
}
func DualNotice (m string) {
    m = fmt.Sprintf("NOTICE: %s", m)
    log.Print(m)
}
func DualWarning (m string) {
    m = fmt.Sprintf("WARNING: %s", m)
    log.Print(m)
}
func DualErr (err error) {
    m := fmt.Sprintf("FATAL: %s", err.Error())
    log.Fatal(m)
}
