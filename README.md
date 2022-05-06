[![CI](https://github.com/sam-atkins/httpc/actions/workflows/ci.yml/badge.svg)](https://github.com/sam-atkins/httpc/actions/workflows/ci.yml)

# httpc

A Go HTTP Client.

## Usage

Example GET requests:

```go
import "github.com/sam-atkins/httpc"

headers := map[string]string{"X-Auth-Token": "topSecretToken"}
res, err := httpc.Get("https://api.com/api/v1/example/").AddHeaders(headers).Do()
```

```go
import "github.com/sam-atkins/httpc"

type simpleJSON struct {
    Data []struct {
        ExampleKey string `json:"exampleKey"`
    } `json:"data"`
    Status string `json:"status"`
}
var sj simpleJSON

err := Get("https://api.com/api/v1/example/").Load(&sj)
```
