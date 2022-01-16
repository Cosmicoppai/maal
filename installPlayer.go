package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func installPlayer() error {
	scriptPath, err := filepath.Abs("scripts/updater.bat")
	if err != nil {
		log.Println(err)
		return err
	}
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "start", scriptPath)
	var _err bytes.Buffer
	cmd.Stderr = &_err
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if HandleInstallError(err, fmt.Sprintf("Error while installing the mpv Player: %s", _err.String())) {
		return err
	}

	processState, _ := cmd.Process.Wait()
	fmt.Println(processState.ExitCode(), processState.Exited())
	if processState.Exited() && (processState.ExitCode() == 0) {
		pathVar := os.Getenv("PATH")
		playerPATH, _ := os.Getwd()
		newPATH := fmt.Sprintf("%s;%s", pathVar, playerPATH)
		err = os.Setenv("PATH", newPATH)
		if HandleInstallError(err, "Error while adding the executable to path") {
			return err
		}

		return nil
	}
	return errors.New("installation process not completed")
}
