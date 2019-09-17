package main

import "github.com/airenas/google-test-runner/internal/app"

func main() {
	config, err := app.InitSettings()
	if err != nil {
		panic(err)
	}
	app.Run(config)
}
