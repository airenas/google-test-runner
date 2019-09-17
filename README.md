# Google Test Runner

## Intro

It is a simple tools to run many googletest suits (separat executables) and collect overall statistics.

## Why

[C++ googletest](https://github.com/google/googletest) is great lib to test C++ code. But I run into a problem how to run many googletest executables locally and collect an overall statistics. Lets say tests are fast, but initialization for a suite takes some time. For exammple I have 20 test case suits. Initialization (*::testing::Test::SetUpTestCase*) for each suit takes about 1 sec. And I have about 500 test.

[`Ctest -j4`](https://cmake.org/) does the thing. But it runs every test separatelly, invokes *::testing::Test::SetUpTestCase* for each test. It ends up to 500 sec/4 (number of processes) = 125 sec.

I can run all tests using `find . -type f -name '*_test' -exec bash -c {} ';'`. But it does not show failures at the end.

## Installation

Before: You must have [go](https://golang.org/) >= 1.12 installed.

1) Clone the repo `git clone https://github.com/airenas/google-test-runner`

2) Build `cd google-test-runner/cmd/google-test-runner && go build`

3) Copy somewhere on your path `sudo cp google-test-runner /usr/local/bin`

## Usage

Run `google-test-runner -h` for a help.

1) Go to a dir containing googletest executables inside

2) Lets say `find . -type f -name '*_test'` - returns googletests with relative paths

3) Invoke the tests:

 * Standard: `find . -type f -name '*_test' | google-test-runner`
 * Filter as parameter: `google-test-runner -f ./**/**/*_test`
 * Show only failing tests: `google-test-runner -s -f ./**/**/*_test`
 * One worker: `google-test-runner -j 1 -s -f ./**/**/*_test`

---

## Author

**Airenas Vaičiūnas**

* [bitbucket.org/airenas](https://bitbucket.org/airenas)
* [github.com/airenas](https://github.com/airenas)
* [linkedin.com/in/airenas](https://www.linkedin.com/in/airenas/)

---

## License

Copyright © 2019, [Airenas Vaičiūnas](https://bitbucket.org/airenas).

---