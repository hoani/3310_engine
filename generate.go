package main

//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname EffortsPro -package somepx -output example/assets/fonts/somepx/efforts_pro.go example/assets/fonts/somepx/ttf/EffortsPro.ttf  -yadvance 8 -dpi 96
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname Tiny -package mwelch -output example/assets/fonts/mwelch/tiny.go example/assets/fonts/mwelch/ttf/tiny.ttf -size 6 -yadvance 5 -dpi=72
