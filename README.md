# httpc

A Go HTTP Client.

## Usage

Example GET request:

```go
headers := map[string]string{"X-Auth-Token": c.token}
res, err := httpc.Get("https://api.com/api/v1/example/").AddHeaders(headers).Do()
```
