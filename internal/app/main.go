package app

import (
	"sync"

	"github.com/airenas/google-test-runner/internal/app/google"
)

type reader interface {
	ReadExecutables() ([]string, error)
}

type formatter interface {
	ShowSuiteStart(string)
	ShowSuiteFailure(string, error)
	ShowTests(*google.TestResult)
	ShowStatistics([]*google.TestResult)
}

//Run is main entry point for te app
func Run(config *Settings) {
	fileNames, err := config.reader.ReadExecutables()
	if err != nil {
		panic(err)
	}

	var resMutex = &sync.Mutex{}
	res := make([]*google.TestResult, 0)

	var tasksWg sync.WaitGroup
	tasksWg.Add(len(fileNames))
	wLimit := make(chan bool, config.workersCount)

	for _, line := range fileNames {
		wLimit <- true
		go func(file string) {
			config.formatter.ShowSuiteStart(file)
			defer tasksWg.Done()
			defer func() { <-wLimit }()

			gr, err := runGoogleTest(file, config.workingDir)
			if err != nil {
				config.formatter.ShowSuiteFailure(file, err)
			} else {
				resMutex.Lock()
				defer resMutex.Unlock()
				res = append(res, gr)
				config.formatter.ShowTests(gr)
			}
		}(line)
	}

	tasksWg.Wait()
	config.formatter.ShowStatistics(res)
}
