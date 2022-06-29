# gracy

A small tool that helps you with gracefull shutdown. 

## Getting started
To download the package, run:
```bash
go get github.com/hanagantig/gracy
```

Import it in your program as:
```go
import "github.com/hanagantig/gracy"
```

A simple usage:
```go
myServer := http.NewServer()
gracy.AddCallback(func() error {
  return myServer.Stop()
})

err := gracy.Wait()
if err != nil {
	logger.Error("failed to gracefully shutdown server")
}
```
