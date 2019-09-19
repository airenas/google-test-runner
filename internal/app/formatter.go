package app

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/airenas/google-test-runner/internal/app/google"
	"github.com/gookit/color"
)

var red = color.FgRed.Render
var green = color.FgGreen.Render

type formatterWriter struct {
	writer            io.Writer
	showOnlyFailed    bool
	showGTestOutput   bool
	skipTestSuites    bool
	skipTestCaseStart bool
	skipTestCases     bool
	skipTestStart     bool
}

type stats struct {
	allTest         int
	allTestDuration time.Duration
	allSuits        int
	failed          int
	succeded        int
}

func newFormatterWriter(w io.Writer, formatInfo string) *formatterWriter {
	result := &formatterWriter{writer: w}
	result.showOnlyFailed = strings.Contains(formatInfo, "f")
	result.showGTestOutput = strings.Contains(formatInfo, "o")
	result.skipTestSuites = strings.Contains(formatInfo, "s")
	result.skipTestCaseStart = strings.Contains(formatInfo, "a")
	result.skipTestCases = strings.Contains(formatInfo, "c")
	result.skipTestStart = strings.Contains(formatInfo, "t")
	return result
}

//
// [==========] 97 tests from 2 test cases ran. (1757 ms total)
// [  PASSED  ] 93 tests.
// [  FAILED  ] 4 tests, listed below:
// [  FAILED  ] AccenterTest.accent_gydomasis
// [  FAILED  ] AccenterTest.accent_fill_from_verb
// [  FAILED  ] AccenterTest.accent_fill_from_verb2
// [  FAILED  ] AccenterTest.accent_auklejamasis
func (f *formatterWriter) ShowStatistics(data []*google.TestResult) {
	st := stats{}
	failed := make([]string, 0)
	for _, gt := range data {
		for _, ts := range gt.Testsuites {
			st.allSuits++
			d, _ := time.ParseDuration(ts.Time)
			st.allTestDuration = st.allTestDuration + d
			for _, t := range ts.Testsuite {
				st.allTest++
				if len(t.Failures) == 0 {
					st.succeded++
				} else {
					st.failed++
					failed = append(failed, ts.Name+"."+t.Name)
				}

			}
		}
	}
	fmt.Fprintf(f.writer, "\n%s %d test%s from %d test case%s ran. (%s total)\n", green("[==========]"),
		st.allTest, sOrEmpty(st.allTest), st.allSuits, sOrEmpty(st.allSuits), durationAsStr(st.allTestDuration))
	if st.succeded > 0 {
		fmt.Fprintf(f.writer, "%s %d test%s\n", green("[  PASSED  ]"), st.succeded, sOrEmpty(st.succeded))
	}
	if st.failed > 0 {
		fmt.Fprintf(f.writer, "%s %d test%s, listed below:\n", red("[  FAILED  ]"), st.failed, sOrEmpty(st.failed))
	}
	for _, ft := range failed {
		fmt.Fprintf(f.writer, "%s %s\n", red("[  FAILED  ]"), ft)
	}
}

func sOrEmpty(i int) string {
	if i > 1 {
		return "s"
	}
	return ""
}

func durationAsStr(i time.Duration) string {
	return i.String()
}

func (f *formatterWriter) ShowSuiteStart(file string) {
	if !f.skipTestSuites {
		fmt.Fprintf(f.writer, "%s Starting %s\n", green("[==========]"), file)
	}
}

func (f *formatterWriter) ShowSuiteFailure(file string, output string, err error) {
	f.printGTestOutput(output)
	fmt.Fprintf(f.writer, "%s File failed: %s\n%s\n", red("[  FAILED  ]"), file, err.Error())
}

func (f *formatterWriter) printGTestOutput(output string) {
	if f.showGTestOutput {
		fmt.Fprintf(f.writer, "\n\n\n<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>\n")
		fmt.Fprintf(f.writer, "<<<<<<<<< START GTEST OUTPUT >>>>>>>>>\n")
		fmt.Fprintf(f.writer, "<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>\n\n")
		fmt.Fprintf(f.writer, "%s\n", output)
		fmt.Fprintf(f.writer, "<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>\n")
		fmt.Fprintf(f.writer, "<<<<<<<<<< END GTEST OUTPUT >>>>>>>>>>\n")
		fmt.Fprintf(f.writer, "<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>\n\n\n")
	}
}

// [----------] 2 tests from Accenter
// [ RUN      ] Accenter.init
// [       OK ] Accenter.init (244 ms)
// [ RUN      ] Accenter.fail_init
// [       OK ] Accenter.fail_init (0 ms)
// [----------] 2 tests from Accenter (244 ms total)
func (f *formatterWriter) ShowTests(data *google.TestResult, output string) {
	f.printGTestOutput(output)
	for _, ts := range data.Testsuites {
		if !(f.skipTestCases || f.skipTestCaseStart) {
			fmt.Fprintf(f.writer, "%s %d test%s from %s\n", green("[----------]"), len(ts.Testsuite), sOrEmpty(len(ts.Testsuite)), ts.Name)
		}
		for _, t := range ts.Testsuite {
			if len(t.Failures) == 0 {
				if !f.showOnlyFailed {
					if !f.skipTestStart {
						fmt.Fprintf(f.writer, "%s %s.%s\n", green("[ RUN      ]"), ts.Name, t.Name)
					}
					fmt.Fprintf(f.writer, "%s %s.%s (%s)\n", green("[       OK ]"), ts.Name, t.Name, t.Time)
				}
			} else {
				fmt.Fprintf(f.writer, "%s %s.%s\n", green("[ RUN      ]"), ts.Name, t.Name)
				for _, fl := range t.Failures {
					fmt.Fprintf(f.writer, "%s %s\n", fl.Failure, fl.Type)
				}
				fmt.Fprintf(f.writer, "%s %s.%s (%s)\n", red("[  FAILED  ]"), ts.Name, t.Name, t.Time)
			}

		}
		if !f.skipTestCases {
			fmt.Fprintf(f.writer, "%s %d test%s from %s (%s total)\n", green("[----------]"), len(ts.Testsuite), sOrEmpty(len(ts.Testsuite)), ts.Name,
				ts.Time)
		}
	}
}
