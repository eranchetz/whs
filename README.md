# whs
Wait HTTP Server | An HTTP Server that holds the request hostage for n seconds

#### defaults
* wait = 0s
* ignore-sigterm = true

## Run 

```go
docker run -ti --rm  -p 8080:8080 eranchetz/whs -wait 2 -ignore-sigterm false
```

# Use 

```bash
curl localhost:8080/?wait=750ms
```

wait is parsed by these [rules](https://golang.org/pkg/time/#ParseDuration)