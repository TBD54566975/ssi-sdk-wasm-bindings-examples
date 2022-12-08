# Minimal demo of WASM bindings of ssi-sdk into javascript

https://github.com/TBD54566975/ssi-sdk is the home of self soverign stuff at TBD, implented in golang. We wan to use this from the web as well, this minimal demo shows how.

# Building JS WASM bindings

```bash
GOOS=js GOARCH=wasm go build -tags jwx_es256k -o ./static/main.wasm sample-app #gosetup
```

# Running Web Server to test bindings 
```
go run webserver/main.go
```
`webserver/index.html` shows how to load the `main.wasm` binary and access it as a javascript function.


`main.go` has the glue code that adds golang functions to javascript and shows the pattern to follow.