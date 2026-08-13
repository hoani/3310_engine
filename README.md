
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
./tools/build/build.sh wasm ./examples/sprite 
```

The tools/build/build.sh tool should be copied to your project.
