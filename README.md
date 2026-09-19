# 3310 Engine

This engine was originally built for the [Nokia 3310 Game Jam 8](https://itch.io/jam/nokiajam8).

I used the engine to build [Trick Shot II](https://itch.io/jam/nokiajam8/rate/4922590) which won first place in the Jam.

The intention was to build an engine which runs both on desktop/browser and on a microcontroller so games can be played on a real 3310.


## Fonts

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
