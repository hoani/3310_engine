
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
./tools/build/wasm.sh ./examples/sprite 
```

The tools/build/wasm.sh tool should be copied to your project.
