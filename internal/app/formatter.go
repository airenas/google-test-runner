package app

import (
	"fmt"
	"io"
	"time"

	"github.com/airenas/google-test-runner/internal/app/google"
	"github.com/gookit/color"
)

var red = color.FgRed.Render
var green = color.FgGreen.Render

type formatterWriter struct {
	writer         io.Writer
	showOnlyFailed bool
}

type stats struct {
	allTest         int
	allTestDuration time.Duration
	allSuits        int
	failed          int
	succeded        int
}

func newFormatterWriter(w io.Writer, showOnlyFailed bool) *formatterWriter {
	return &formatterWriter{writer: w, showOnlyFailed: showOnlyFailed}
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
	fmt.Fprintf(f.writer, "%s Starting %s\n", green("[==========]"), file)
}

func (f *formatterWriter) ShowSuiteFailure(file string, err error) {
	fmt.Fprintf(f.writer, "%s File failed: %s\n%s\n", red("[  FAILED  ]"), file, err.Error())
}

// [----------] 2 tests from Accenter
// [ RUN      ] Accenter.init
// [       OK ] Accenter.init (244 ms)
// [ RUN      ] Accenter.fail_init
// [       OK ] Accenter.fail_init (0 ms)
// [----------] 2 tests from Accenter (244 ms total)
func (f *formatterWriter) ShowTests(data *google.TestResult) {
	for _, ts := range data.Testsuites {
		fmt.Fprintf(f.writer, "%s %d test%s from %s\n", green("[----------]"), len(ts.Testsuite), sOrEmpty(len(ts.Testsuite)), ts.Name)
		for _, t := range ts.Testsuite {
			if len(t.Failures) == 0 {
				if !f.showOnlyFailed {
					fmt.Fprintf(f.writer, "%s %s.%s\n", green("[ RUN      ]"), ts.Name, t.Name)
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
		fmt.Fprintf(f.writer, "%s %d test%s from %s (%s total)\n", green("[----------]"), len(ts.Testsuite), sOrEmpty(len(ts.Testsuite)), ts.Name,
			ts.Time)
	}
}