# whs
Wait HTTP Server | An HTTP Server that holds the request hostage for n seconds


## Run 

```go
docker run -ti --rm  -p 8080:8080 eranchetz/whs -wait 2
```