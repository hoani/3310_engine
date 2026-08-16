package main

//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname EffortsPro -package font -output font/efforts_pro.go font/ttf/EffortsPro.ttf -yadvance 8 -dpi 96
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname Tiny -package font -output font/tiny.go font/ttf/tiny.ttf -size 6 -yadvance 5 -dpi=72
