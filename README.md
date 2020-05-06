# gologger

## Install 
`go get github.com/kaizer666/gologger`

## Usage

```go
package main

import GoLogger `github.com/kaizer666/gologger`

var logger GoLogger.Logger

func main() {
    logger = GoLogger.Logger{}
    logger.SetLogLevel(0)
    logger.SetLogFileName("main.log")
    err := logger.Init()
    panic(err)
}
```
