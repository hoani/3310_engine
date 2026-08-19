// Command img2pbm converts PNG sprites into PBM (P4) or PGM (P5) files.
//
// Two outputs are available for PBM:
// * Monochrome - outputs a single black and while file
// * Masked Monochrome - outputs two files, one for color, one as a transparency mask
//
// For PGM (P5)
// Three states are encoded in the gray value so a single
// file stays fully previewable in tools which support p5.
//
//	  0 = ink        (opaque + dark  -> black on the display)
//	255 = paper      (opaque + light -> white on the display)
//	128 = transparent(alpha 0        -> mask off; skip when blitting)
//
// Because the display treats 1 as black and PGM treats 0 as black, ink maps
// to 0 here and previews match the screen with no inversion.
//
//	//go:generate go run ./tools/png2pbm -p4 -in assets/png -out assets/pbm
//	//go:generate go run ./tools/png2pbm -p4x -in assets/png -out assets/pbm_duo
//	//go:generate go run ./tools/png2pbm -p5 -in assets/png -out assets/pgm
//
// then run `go generate ./...`.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const (
	valInk         = 0    // opaque + dark
	valPaper       = 0xf9 // opaque + light
	valTransparent = 0xff // alpha below the cut
)

type format struct {
	p4Ink     bool
	p4Alpha   bool
	p5        bool
	threshold int
	alphaCut  int
}

func main() {
	in := flag.String("in", ".", "directory of source PNG files")
	out := flag.String("out", ".", "directory to write generated files")
	threshold := flag.Int("threshold", 128, "luminance 0-255 below which an opaque pixel is ink")
	alphaCut := flag.Int("alpha", 128, "alpha 0-255 at or above which a pixel is opaque")
	p4 := flag.Bool("p4", false, "generate p4 pbm")
	p4x := flag.Bool("p4x", false, "generate extended p4 pbm with ink and transparency")
	p5 := flag.Bool("p5", false, "generate grayscale pgm")
	flag.Parse()

	f := &format{
		threshold: *threshold,
		alphaCut:  *alphaCut,
	}
	if *p4 || *p4x {
		f.p4Ink = true
	}
	if *p4x {
		f.p4Alpha = true
	}
	if *p5 {
		f.p5 = true
	}

	if err := run(*in, *out, f); err != nil {
		fmt.Fprintf(os.Stderr, "img2p5: %e", err)
		os.Exit(1)
	}
}

func run(inDir, outDir string, f *format) error {
	entries, err := os.ReadDir(inDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	count := 0
	imgs := []string{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.EqualFold(filepath.Ext(e.Name()), ".png") {
			fmt.Fprintf(os.Stdout, "skipping: %s", e.Name())
			continue
		}

		generated, err := f.generate(inDir, outDir, e)
		if err != nil {
			return err
		}
		imgs = append(imgs, generated...)

		count++
	}

	embName := filepath.Base(outDir) + "_" + filepath.Base(inDir) + "_" + f.Name()
	embName += ".go"
	writeEmbeds(outDir, embName, imgs)

	fmt.Printf("png2pbm: converted %d file(s)\n", count)
	return nil
}

func (f *format) Name() string {
	name := ""
	if f.p4Ink {
		name += "p4"
	}
	if f.p4Alpha {
		name += "a"
	}
	if f.p5 {
		name += "p5"
	}
	return name
}

func (f *format) generate(inDir, outDir string, e os.DirEntry) ([]string, error) {
	imgs := []string{}
	if f.p4Ink {
		img, err := f.generateP4(inDir, outDir, e)
		if err != nil {
			return []string{}, err
		}
		imgs = append(imgs, img)
	}
	if f.p4Alpha {
		img, err := f.generateP4Alpha(inDir, outDir, e)
		if err != nil {
			return []string{}, err
		}
		imgs = append(imgs, img)
	}
	if f.p5 {
		img, err := f.generateP5(inDir, outDir, e)
		if err != nil {
			return []string{}, err
		}
		imgs = append(imgs, img)
	}

	return imgs, nil
}

func (f *format) generateP4(inDir, outDir string, e os.DirEntry) (string, error) {
	src := filepath.Join(inDir, e.Name())
	imgName := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name())) + ".pbm"
	dst := filepath.Join(outDir, imgName)
	if err := f.convertP4(src, dst, false); err != nil {
		return "", fmt.Errorf("%s: %w", e.Name(), err)
	}
	fmt.Printf("png2pbm: %s -> %s\n", src, dst)
	return imgName, nil
}

func (f *format) generateP4Alpha(inDir, outDir string, e os.DirEntry) (string, error) {
	src := filepath.Join(inDir, e.Name())
	imgName := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name())) + "_mask.pbm"
	dst := filepath.Join(outDir, imgName)
	if err := f.convertP4(src, dst, true); err != nil {
		return "", fmt.Errorf("%s: %w", e.Name(), err)
	}
	fmt.Printf("png2pbm: %s -> %s\n", src, dst)
	return imgName, nil
}

func (f *format) generateP5(inDir, outDir string, e os.DirEntry) (string, error) {
	src := filepath.Join(inDir, e.Name())
	imgName := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name())) + ".pgm"
	dst := filepath.Join(outDir, imgName)
	if err := f.convertP5(src, dst); err != nil {
		return "", fmt.Errorf("%s: %w", e.Name(), err)
	}
	fmt.Printf("png2pbm: %s -> %s\n", src, dst)
	return imgName, nil
}

func encodeShade(gray int) byte {
	level := (gray*32 + 127) / 255 // round(gray/255 * 32) -> 0..32
	if level >= 32 {
		return valPaper
	}
	return byte(level << 3)
}

func (f *format) convertP4(src, dst string, alphaOnly bool) error {
	fi, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fi.Close()

	img, err := png.Decode(fi)
	if err != nil {
		return err
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	pix := NewBitWriter(w, h)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, a := img.At(x, y).RGBA()

			if int(a>>8) < f.alphaCut {
				pix.Put(false)
				continue
			}

			if alphaOnly {
				pix.Put(true)
				continue
			}

			// Perceived brighness
			lum := (299*int(r>>8) + 587*int(g>>8) + 114*int(bl>>8)) / 1000

			if lum < f.threshold {
				pix.Put(true)
			} else {
				pix.Put(false)
			}
		}
		pix.EndRow()
	}

	return writeP4(dst, w, h, pix.pix)
}

func (f *format) convertP5(src, dst string) error {
	fi, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fi.Close()

	img, err := png.Decode(fi)
	if err != nil {
		return err
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	pix := make([]byte, 0, w*h)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, a := img.At(x, y).RGBA()

			if int(a>>8) < f.alphaCut {
				pix = append(pix, valTransparent)
				continue
			}
			// Perceived brighness
			lum := (299*int(r>>8) + 587*int(g>>8) + 114*int(bl>>8)) / 1000
			pix = append(pix, encodeShade(lum))
		}
	}

	return writeP5(dst, w, h, pix)
}

func writeP5(path string, w, h int, pix []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	// P5 header: magic \n width space height \n maxval \n, then raw raster,
	// one byte per pixel, top-to-bottom, left-to-right, no row padding.
	fmt.Fprintf(bw, "P5\n%d %d\n255\n", w, h)
	bw.Write(pix)
	return bw.Flush()
}

func writeP4(path string, w, h int, pix []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	// P4 header: magic \n width space height \n, then raw raster,
	// one byte per 8 pixels, top-to-bottom, left-to-right, no row padding.
	fmt.Fprintf(bw, "P4\n%d %d\n", w, h)
	bw.Write(pix)
	return bw.Flush()
}

func writeEmbeds(path, name string, images []string) error {
	type Entry struct {
		VarName  string
		FileName string
	}

	type templateData struct {
		Package string
		Images  []Entry
	}

	data := templateData{
		Package: filepath.Base(path),
		Images:  make([]Entry, len(images)),
	}

	for i, img := range images {
		data.Images[i].FileName = img
		data.Images[i].VarName = toCamel(img)
	}

	var tmpl = template.Must(template.New("imgs").Parse(`// Generated by img2pbm;
package {{.Package}}

import (
	_ "embed"
)
{{range .Images}}
//go:embed "{{.FileName}}"
var {{.VarName}} string
{{end}}
`))

	tmplPath := filepath.Join(path, name)
	f, err := os.Create(tmplPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		_ = f.Close()
		return fmt.Errorf("executing template for %s: %w", name, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("closing file for %s: %w", name, err)
	}

	return nil
}

func toCamel(filename string) string {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	fields := strings.FieldsFunc(base, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	var b strings.Builder
	for _, f := range fields {
		r := []rune(f)
		b.WriteRune(unicode.ToUpper(r[0]))
		b.WriteString(string(r[1:]))
	}
	name := b.String()
	if name == "" {
		return "Img"
	}
	if unicode.IsDigit([]rune(name)[0]) { // identifiers can't start with a digit
		name = "Img" + name
	}
	return name
}
