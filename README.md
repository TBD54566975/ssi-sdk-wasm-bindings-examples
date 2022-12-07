![wasm usage gif](./ssi-sdk-wasm-example.gif)

^ shows using go binding in the browser to generate a DID in JS land

# Building JS WASM bindings

```bash
GOOS=js GOARCH=wasm go build -tags jwx_es256k -o ./static/main.wasm sample-app #gosetup
```

# Running Web Server to test bindings 
```
go run webserver/main.go
```

