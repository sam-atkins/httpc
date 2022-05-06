[![CI](https://github.com/sam-atkins/httpc/actions/workflows/ci.yml/badge.svg)](https://github.com/sam-atkins/httpc/actions/workflows/ci.yml)

# httpc

A Go HTTP Client.

## Usage

Import the package.

```go
import "github.com/sam-atkins/httpc"
```

### Example GET requests

```go
headers := map[string]string{"X-Auth-Token": "topSecretToken"}
res, err := httpc.Get("https://api.com/api/v1/example/").AddHeaders(headers).Do()
```

```go
type simpleJSON struct {
    Data []struct {
        ExampleKey string `json:"exampleKey"`
    } `json:"data"`
    Status string `json:"status"`
}
var sj simpleJSON

err := Get("https://api.com/api/v1/example/").Load(&sj)
```

### Example POST requests

```go
url := "https://api.com/api/v1/example/"
type requestBody struct {
    Text  string
    Token string
}
body := &requestBody{
    Text:  "this is some text",
    Token: "mySecretToken",
}
res, err := Post(url, body).Do()
```

```go
url := "https://api.com/api/v1/example/"
type requestBody struct {
    Text  string
}
body := &requestBody{
    Text:  "this is some text",
}

type simpleJSON struct {
    Data []struct {
        ExampleKey string `json:"exampleKey"`
    } `json:"data"`
    Status string `json:"status"`
}
var sj simpleJSON

headers := map[string]string{"X-Auth-Token": "topSecretToken"}

err := Post(url, body).AddHeaders(headers).Load(&sj)
```
