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
	ShowSuiteFailure(string, string, error)
	ShowTests(*google.TestResult, string)
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
	workerQueueLimit := make(chan bool, config.workersCount) 

	for _, line := range fileNames {
		workerQueueLimit <- true // try get access to work
		go func(file string) {
			defer tasksWg.Done()
			defer func() { <-workerQueueLimit }() // decrease working queue

			config.formatter.ShowSuiteStart(file)

			gr, output, err := runGoogleTest(file, config.workingDir)
			if err != nil {
				config.formatter.ShowSuiteFailure(file, output, err)
			} else {
				resMutex.Lock()
				defer resMutex.Unlock()
				res = append(res, gr)
				config.formatter.ShowTests(gr, output)
			}
		}(line)
	}

	tasksWg.Wait()

	config.formatter.ShowStatistics(res)
}
