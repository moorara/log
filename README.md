[![Go Doc][godoc-image]][godoc-url]
[![Build Status][workflow-image]][workflow-url]
[![Go Report Card][goreport-image]][goreport-url]
[![Test Coverage][coverage-image]][coverage-url]
[![Maintainability][maintainability-image]][maintainability-url]

# log

This package can be used for production-ready logging in Go applications.
It hides the complexity of configuring and using the state-of-the-arts loggers
by providing a **single interface** that is _easy-to-use_ and _hard-to-misuse_!

## Quick Start

You can either log using an _instance_ logger or the _singleton_ logger.
After creating an instance logger, you need to set the singleton logger using the `SetSingleton` method once.
The instance logger can be further used to create more contextualized loggers as the children of the root logger.

### [zap](https://github.com/uber-go/zap)

```go
package main

import "github.com/moorara/log"

func main() {
  // Creating a zap logger
  logger := log.NewZap(log.Options{
    Name:        "my-service",
    Version:     "0.1.0",
    Environment: "production",
    Region:      "us-east-1",
    Tags: map[string]string{
      "domain": "auth",
    },
  })

  // Initializing the singleton logger
  log.SetSingleton(logger)

  // Logging using the singleton logger
  log.Infof("starting server on port %d ...", 8080)

  // Logging using the instance logger
  logger.Info("request received.",
    "tenantId", "aaaaaaaa",
    "requestId", "bbbbbbbb",
  )
}
```

Output logs from stdout:

```json
{"level":"info","timestamp":"2020-04-24T12:39:04.506116-04:00","caller":"example/main.go:21","message":"starting server on port 8080 ...","domain":"auth","environment":"production","logger":"my-service","region":"us-east-1","version":"0.1.0"}
{"level":"info","timestamp":"2020-04-24T12:39:04.506268-04:00","caller":"example/main.go:24","message":"request received.","domain":"auth","environment":"production","logger":"my-service","region":"us-east-1","version":"0.1.0","tenantId":"aaaaaaaa","requestId":"bbbbbbbb"}
```

### [go-kit](https://github.com/go-kit/kit/tree/master/log)

```go
package main

import "github.com/moorara/log"

func main() {
  // Creating a kit logger
  logger := log.NewKit(log.Options{
    Name:        "my-service",
    Version:     "0.1.0",
    Environment: "production",
    Region:      "us-east-1",
    Tags: map[string]string{
      "domain": "auth",
    },
  })

  // Initializing the singleton logger
  log.SetSingleton(logger)

  // Logging using the singleton logger
  log.Infof("starting server on port %d ...", 8080)

  // Logging using the instance logger
  logger.Info("request received.",
    "tenantId", "aaaaaaaa",
    "requestId", "bbbbbbbb",
  )
}
```

Output logs from stdout:

```json
{"caller":"main.go:21","domain":"auth","environment":"production","level":"info","logger":"my-service","message":"starting server on port 8080 ...","region":"us-east-1","timestamp":"2020-04-24T12:39:53.05221-04:00","version":"0.1.0"}
{"caller":"main.go:24","domain":"auth","environment":"production","level":"info","logger":"my-service","message":"request received.","region":"us-east-1","requestId":"bbbbbbbb","tenantId":"aaaaaaaa","timestamp":"2020-04-24T12:39:53.052529-04:00","version":"0.1.0"}
```


[godoc-url]: https://pkg.go.dev/github.com/moorara/log
[godoc-image]: https://godoc.org/github.com/moorara/log?status.svg
[workflow-url]: https://github.com/moorara/log/actions
[workflow-image]: https://github.com/moorara/log/workflows/Main/badge.svg
[goreport-url]: https://goreportcard.com/report/github.com/moorara/log
[goreport-image]: https://goreportcard.com/badge/github.com/moorara/log
[coverage-url]: https://codeclimate.com/github/moorara/log/test_coverage
[coverage-image]: https://api.codeclimate.com/v1/badges/5401f0f63ecbb401202f/test_coverage
[maintainability-url]: https://codeclimate.com/github/moorara/log/maintainability
[maintainability-image]: https://api.codeclimate.com/v1/badges/5401f0f63ecbb401202f/maintainability
