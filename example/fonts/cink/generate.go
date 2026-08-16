package cink

//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname BittypixMonospace -package cink -output bittypix_monospace.go ttf/BittypixMonospace.otf -yadvance 8 -dpi 48

//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname CodersCrux -package cink -output coders_crux.go ttf/CodersCrux.otf -yadvance 8 -dpi 96
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname DiaryOfAn8bitMage -package cink -output diary_of_an_8bit_mage.go ttf/DiaryOfAn8bitMage.otf -yadvance 9 -dpi 54
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname Frogotype -package cink -output frogotype.go ttf/Frogotype.ttf -yadvance 9 -dpi 96
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname MesseDuesseldorf -package cink -output messe_duesseldorf.go ttf/MesseDuesseldorf.ttf -yadvance 10 -dpi 70
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname NineteenEightySeven -package cink -output nineteen_eighty_seven.go ttf/NineteenEightySeven.otf -yadvance 8 -dpi 42
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname NineteenNinetySeven -package cink -output nineteen_ninety_seven.go ttf/NineteenNinetySeven.otf -yadvance 9 -dpi 54
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname NineteenNinetySix -package cink -output nineteen_ninety_six.go ttf/NineteenNinetySix.otf -yadvance 8 -dpi 54
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname NineteenNinetyThree -package cink -output nineteen_ninety_three.go ttf/NineteenNinetyThree.otf -yadvance 8 -dpi 54
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname Notalot18 -package cink -output notalot18.go ttf/Notalot18.ttf -yadvance 7 -dpi 42
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname Notalot25 -package cink -output notalot25.go ttf/Notalot25.ttf -yadvance 7 -dpi 48
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname PixelOrGTFO -package cink -output pixel_or_gtfo.go ttf/Pixel-Or-GTFO.ttf -yadvance 7 -dpi 48
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname PrincessSavesYou -package cink -output princess_saves_you.go ttf/princess-saves-you.otf -yadvance 9 -dpi 54
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname Rygarde -package cink -output rygarde.go ttf/Rygarde.ttf -yadvance 8 -dpi 48
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname SuperLegendBoy -package cink -output super_legend_boy.go ttf/super-legend-boy.otf -yadvance 8 -dpi 54
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname TeenyTinyPixls -package cink -output teeny_tiny_pixls.go ttf/TeenyTinyPixls.otf -yadvance 6 -dpi 30
//go:generate go run tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@v0.7.0 -fontname TinyAndChunky -package cink -output tiny_and_chunky.go ttf/tiny-and-chunky.ttf -yadvance 6 -dpi 30
