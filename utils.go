package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func runCommand(cmd string) (output string, err error) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	return string(out), err
}

func replaceData(path string, originalString string, newString string) (err error) {
	if fileData, err := ioutil.ReadFile(path); err != nil {
		log.Fatal(err)
	} else {
		newContent := []byte(strings.Replace(string(fileData), originalString, newString, -1))
		err = ioutil.WriteFile(path, newContent, 0)
	}
	return err
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
