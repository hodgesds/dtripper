# Dtripper
golang HTTP over DNS

This project was inspired by a [fun idea](https://github.com/hodgesds/ztripper) I had awhile go and I wanted to expand upon it.

# Example
Running the DNS server:
```sh
   go run server/server.go -d
```

Running the websocket server:
```sh
   go run server/ws_server.go -d
```

Run the example http client:
```sh
   go run example/combined.go
```

Do a request over DNS:
```sh
   go run example/combined.go -url dns://foo.com
```

Do a request over websocket:
```sh
   go run example/combined.go -url ws://foo.com
```
