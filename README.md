[![CI](https://github.com/sam-atkins/httpc/actions/workflows/ci.yml/badge.svg)](https://github.com/sam-atkins/httpc/actions/workflows/ci.yml)

# httpc

A Go HTTP Client.

## Usage

Example GET request:

```go
headers := map[string]string{"X-Auth-Token": "topSecretToken"}
res, err := httpc.Get("https://api.com/api/v1/example/").AddHeaders(headers).Do()
```
