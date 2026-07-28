package main

//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname EffortsPro -package font -output font/efforts_pro.go font/ttf/EffortsPro.ttf -dpi 96
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname Tiny -package font -output font/tiny.go font/ttf/tiny.ttf -size 6 -dpi=72

//go:generate go run ./tools/img2p5 -in sprites/mono -out sprites/pgm
//go:generate go run ./tools/img2p5 -gray -in sprites/gray -out sprites/pgm
