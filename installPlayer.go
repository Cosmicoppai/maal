package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func installPlayer() (string, error) {
	var path string
	path1, err1 := exec.LookPath("mpv")
	path2, err2 := exec.LookPath("scripts/mpv")
	if (err1 != nil) && (err2 != nil) {

		scriptPath, err := filepath.Abs("scripts/updater.bat")
		if err != nil {
			log.Println(err)
			return "", err
		}
		cmd := exec.Command("powershell", "-NoProfile", "-Command", "start", scriptPath)
		var _err bytes.Buffer
		cmd.Stderr = &_err
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if HandleInstallError(err, fmt.Sprintf("Error while installing the mpv Player: %s", _err.String())) {
			return "", err
		}
		return "scripts/mpv", nil
	}
	if path1 == "" {
		path = path2
	} else {
		path = path1
	}
	return path, nil
}
