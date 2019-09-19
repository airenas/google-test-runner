package app

import (
	"bytes"
	"testing"

	"github.com/airenas/google-test-runner/internal/app/google"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestShowSuiteStart(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "")

	f.ShowSuiteStart("file")

	assert.Equal(t, "\x1b[32m[==========]\x1b[0m Starting file\n", sw.String())
}

func TestShowSuiteStart_Skip(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "s")

	f.ShowSuiteStart("file")

	assert.NotContains(t, sw.String(), "Starting file")
}

func TestShowSuite_Failure(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "")

	f.ShowSuiteFailure("file", "", errors.New("error"))

	assert.Equal(t, "\x1b[31m[  FAILED  ]\x1b[0m File failed: file\nerror\n", sw.String())
}

func TestShowSuite_FailureOutput(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "o")

	f.ShowSuiteFailure("file", "Output", errors.New("error"))

	assert.Contains(t, sw.String(), "Output")
}

func TestShowTests_AllPassed(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "")

	data := makeData()

	f.ShowTests(data, "")

	assert.Contains(t, sw.String(), "1 test from olia")
	assert.Contains(t, sw.String(), "RUN")
	assert.Contains(t, sw.String(), "OK")
}

func TestShowTests_SkipCasesStart(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "a")

	data := makeData()

	f.ShowTests(data, "")

	assert.NotContains(t, sw.String(), "from olia\n")
}

func TestShowTests_SkipCases(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "c")

	data := makeData()

	f.ShowTests(data, "")

	assert.NotContains(t, sw.String(), "from olia")
}

func TestShowTests_Output(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "o")

	data := makeData()

	f.ShowTests(data, "Output")

	assert.Contains(t, sw.String(), "Output")
}

func TestShowTests_AllPassed_Skip(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "f")

	data := makeData()

	f.ShowTests(data, "")

	assert.Contains(t, sw.String(), "1 test from olia")
	assert.NotContains(t, sw.String(), "RUN")
	assert.NotContains(t, sw.String(), "OK")
}

func TestShowTests_SkipStart(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "t")
	data := makeData()
	f.ShowTests(data, "")
	assert.NotContains(t, sw.String(), "RUN")
}

func TestShowTests_Failure(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "f")

	data := makeData()
	data.Testsuites[0].Testsuite[0].Failures = make([]google.Failure, 1)

	f.ShowTests(data, "")

	assert.Contains(t, sw.String(), "1 test from olia")
	assert.Contains(t, sw.String(), "RUN")
	assert.Contains(t, sw.String(), "FAILED")
}

func TestShowStatistics_AllPass(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "f")

	data := []*google.TestResult{makeData(), makeData()}

	f.ShowStatistics(data)

	assert.Contains(t, sw.String(), "2 tests from 2 test cases ran")
	assert.Contains(t, sw.String(), "PASSED")
	assert.NotContains(t, sw.String(), "FAILED")
}

func TestShowStatistics_Failed(t *testing.T) {
	sw := bytes.NewBufferString("")
	f := newFormatterWriter(sw, "f")

	data := []*google.TestResult{makeData(), makeData()}
	data[0].Testsuites[0].Testsuite[0].Failures = make([]google.Failure, 1)
	data[1].Testsuites[0].Testsuite[0].Failures = make([]google.Failure, 1)

	f.ShowStatistics(data)

	assert.Contains(t, sw.String(), "2 tests from 2 test cases ran")
	assert.NotContains(t, sw.String(), "PASSED")
	assert.Contains(t, sw.String(), "FAILED")
}

func makeData() *google.TestResult {
	data := google.TestResult{}
	data.Testsuites = make([]google.TestResult, 1)
	data.Testsuites[0].Testsuite = make([]google.Testsuite, 1)
	data.Testsuites[0].Name = "olia"
	data.Testsuites[0].Testsuite[0].Name = "testOlia"
	return &data
}
