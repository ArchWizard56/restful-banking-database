// +build !linux

package main

import (
	"log"
)

func SetupLogging () {
    log.Print("Note: Non-Linux systems don't support this implimentation of syslog; writing messages to Standard Output instead")
}

func DualDebug (m string) {
    log.Print(m)
}
func DualInfo (m string) {
    log.Print(m)
}
func DualNotice (m string) {
    log.Print(m)
}
func DualWarning (m string) {
    log.Print(m)
}
func DualErr (err error) {
    log.Fatal(err)
}
