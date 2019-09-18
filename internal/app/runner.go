package app

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/airenas/google-test-runner/internal/app/google"
	"github.com/pkg/errors"
)

func runGoogleTest(f string, wDir string) (*google.TestResult, string, error) {
	file, err := ioutil.TempFile("", "google_test_runner")
	if err != nil {
		return nil, "", errors.Wrap(err, "Can't create temp file ")
	}
	defer os.Remove(file.Name())

	cmd := f + " --gtest_output=json:" + file.Name()
	output, errCmd := runCommand(cmd, wDir)
	gr, err := readJSON(file.Name())
	if err != nil {
		return nil, output, errors.Wrap(errCmd, errors.Wrap(err, "Can't decode json.\n").Error())
	}
	return gr, output, nil
}

func runCommand(command string, wDir string) (string, error) {
	cmdArr := strings.Split(command, " ")
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	cmd.Dir = wDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		errR := errors.Wrap(err, string(output))
		return string(output), errR
	}
	return string(output), nil
}

func readJSON(fn string) (*google.TestResult, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)
	var g google.TestResult
	err = dec.Decode(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
