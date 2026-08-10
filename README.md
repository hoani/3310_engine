
For font generation, install tinyfont:

```
go install tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0
```


## Exporting to WASM

Testing:
```
wasmserve ./example/snd
```

Compiling

```
env GOOS=js GOARCH=wasm go build -o bin/example/snd/game.wasm ./example/snd
cp $(go env GOROOT)/lib/wasm/wasm_exec.js ./bin/example/snd/wasm_exec.js    
```

Add html files:

```html
<!-- main.html -->
<!DOCTYPE html>
<script src="wasm_exec.js"></script>
<script  allow="autoplay">
const go = new Go();
WebAssembly.instantiateStreaming(fetch("./game.wasm"), go.importObject).then(result => {
    go.run(result.instance);
});
</script>
```

```html
<!-- index.html -->
<!DOCTYPE html>
<iframe src="main.html" allow="autoplay" width="640" height="480"></iframe>
```
