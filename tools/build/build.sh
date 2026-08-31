#!/usr/bin/env bash
#
# build.sh — build a Go/Ebitengine package for wasm and/or desktop targets.
#
# Usage: build.sh <target> <package>
#   target: wasm | darwin | linux | windows  (one per invocation)
#
# Output layout (bin/<pkg>/ by default):
#   wasm/                     game.wasm, wasm_exec.js, main.html, index.html
#   <goos>/<goarch>/<name>    the binary (<name>.exe on windows)
#   <name>-wasm.zip           archives sit at the root, one per built target
#   <name>-<goos>-<goarch>.zip
#
set -euo pipefail

TITLE=""
BASEDIR=""
ARCHES=""
LDFLAGS="-s -w"
ALLOW_CROSS=0
DO_ZIP=1

usage() {
	cat <<'USAGEEOF'
Usage: build.sh [options] <target> <package>

  <target>    wasm | darwin | linux | windows
              One target per invocation. Note: darwin/linux require cgo.
  <package>   Go package to build, e.g. ./example/snd

Options:
  -o DIR      base output directory (default: bin/<package path>)
  -a LIST     comma-separated GOARCH list for desktop targets
              (default: amd64,arm64)
  -t TEXT     page title, wasm only (default: package directory name)
  -n          skip building the .zip archives
  -x          attempt cgo cross-compiles that normally get skipped
  -g          keep debug symbols (omit the default -ldflags "-s -w")
  -h          show this help

Examples:
  ./tools/build/build.sh wasm ./example/snd
  ./tools/build/build.sh -a arm64 darwin ./cmd/game
  ./tools/build/build.sh linux ./cmd/game          # on a Linux box
  ./tools/build/build.sh -o dist -W 800 -H 600 -t "My Demo" wasm ./cmd/game
USAGEEOF
}

while getopts ":o:a:W:H:t:nxgh" opt; do
	case "$opt" in
	o) BASEDIR="$OPTARG" ;;
	a) ARCHES="$OPTARG" ;;
	t) TITLE="$OPTARG" ;;
	n) DO_ZIP=0 ;;
	x) ALLOW_CROSS=1 ;;
	g) LDFLAGS="" ;;
	h)
		usage
		exit 0
		;;
	:)
		echo "error: -$OPTARG requires an argument" >&2
		exit 2
		;;
	\?)
		echo "error: unknown option -$OPTARG" >&2
		usage >&2
		exit 2
		;;
	esac
done
shift $((OPTIND - 1))

if [ $# -ne 2 ]; then
	echo "error: expected a target and a package" >&2
	usage >&2
	exit 2
fi
TARGET="$1"
PKG="$2"

case "$TARGET" in
wasm | darwin | linux | windows) ;;
*)
	echo "error: unknown target '$TARGET'" >&2
	usage >&2
	exit 2
	;;
esac

[ -n "$ARCHES" ] || ARCHES="amd64,arm64"
IFS=',' read -r -a ARCH_LIST <<<"$ARCHES"

command -v go >/dev/null 2>&1 || {
	echo "error: go not found in PATH" >&2
	exit 1
}

# An exported GOROOT overrides the toolchain's own root. TinyGo builds a merged
# GOROOT of symlinks under its cache dir, and if that path is still exported,
# plain `go build` tries to compile TinyGo's stdlib fork with the gc compiler
# and dies on TinyGo-only packages (runtime/interrupt, internal/futex) and
# TinyGo-only pragmas (//go:extern). Drop it; the go binary in PATH knows its
# own root.
if [ -n "${GOROOT:-}" ]; then
	echo ">> ignoring exported GOROOT=$GOROOT (using the toolchain's own root)" >&2
	unset GOROOT
fi

# Sanity-check the root the toolchain resolves to. TinyGo builds a merged GOROOT
# of symlinks under its cache dir, overlaying its own stdlib fork. If `go`
# resolves there, the gc compiler tries to build that fork and dies on
# TinyGo-only packages and pragmas. Catch it here rather than 60 lines later.
RESOLVED_GOROOT="$(go env GOROOT)"
TINYGO_ROOT=0
case "$RESOLVED_GOROOT" in *tinygo*) TINYGO_ROOT=1 ;; esac
[ -d "$RESOLVED_GOROOT/src/runtime/interrupt" ] && TINYGO_ROOT=1
if [ "$TINYGO_ROOT" -eq 1 ]; then
	echo "error: go resolves to what looks like a TinyGo goroot:" >&2
	echo "         $RESOLVED_GOROOT" >&2
	echo "       The gc toolchain cannot compile TinyGo's stdlib fork." >&2
	echo "       go binary:  $(command -v go)" >&2
	echo "       go version: $(go version)" >&2
	echo "       go env file: $(go env GOENV)" >&2
	echo "       Check that PATH points at a real Go install, not TinyGo's cache." >&2
	exit 1
fi

# Derive output directory, binary name and page title from the package path.
REL="${PKG#./}"
REL="${REL%/}"
if [ -z "$REL" ] || [ "$REL" = "." ]; then
	REL="$(basename "$PWD")"
fi
NAME="$(basename "$REL")"
[ -n "$BASEDIR" ] || BASEDIR="bin/$REL"
[ -n "$TITLE" ] || TITLE="$NAME"

HOST_OS="$(go env GOHOSTOS)"
ZIP_BIN="$(command -v zip || true)"
FAILED=""
SKIPPED=""

# make_zip <zipfile> <dir> — zip the flat contents of dir (no directory entry).
make_zip() {
	local zipfile="$1"
	local dir="$2"

	if [ "$DO_ZIP" -eq 0 ]; then
		return
	fi
	if [ -z "$ZIP_BIN" ]; then
		echo "   (no zip binary found; skipping $(basename "$zipfile"))" >&2
		return
	fi

	rm -f -- "$zipfile"
	# -j junks paths so the archive holds bare files, not the dir structure.
	# Safe here because every packaged dir is flat.
	if "$ZIP_BIN" -qj "$zipfile" "$dir"/*; then
		echo "   $zipfile ($(du -h "$zipfile" | cut -f1))"
	else
		echo "   FAILED to zip $zipfile" >&2
		FAILED="$FAILED zip:$(basename "$zipfile")"
	fi
}

build_desktop() {
	local goos="$1"
	local ext=""
	[ "$goos" = "windows" ] && ext=".exe"

	# Ebitengine's cgo needs differ per platform. Windows is pure Go (the glfw
	# and DirectX bindings use syscall/purego), so it cross-compiles from
	# anywhere with cgo off. darwin and linux define glfw.Window et al. in
	# cgo-gated files (internal/glfw/window_unix.go), so cgo must be ON or every
	# reference goes undefined — and that in turn needs a C toolchain plus the
	# platform's headers, which only the matching host has.
	local cgo=1
	[ "$goos" = "windows" ] && cgo=0

	if [ "$cgo" -eq 1 ] && [ "$goos" != "$HOST_OS" ] && [ "$ALLOW_CROSS" -eq 0 ]; then
		echo ">> skipping $goos: needs cgo, and this host is $HOST_OS"
		echo "   build it on a $goos host or in CI (-x to attempt anyway)"
		SKIPPED="$SKIPPED $goos"
		return
	fi

	local goarch archdir out
	for goarch in "${ARCH_LIST[@]}"; do
		archdir="$BASEDIR/$goos/$goarch"
		out="$archdir/$NAME$ext"

		# Clear only this arch, so building one arch doesn't wipe a sibling
		# built earlier on another machine.
		rm -rf -- "$archdir"
		mkdir -p -- "$archdir"

		echo ">> building $PKG for $goos/$goarch (CGO_ENABLED=$cgo)"
		if CGO_ENABLED="$cgo" GOOS="$goos" GOARCH="$goarch" \
			go build -trimpath ${LDFLAGS:+-ldflags "$LDFLAGS"} -o "$out" "$PKG"; then
			echo "   $out ($(du -h "$out" | cut -f1))"
			make_zip "$BASEDIR/$NAME-$goos-$goarch.zip" "$archdir"
		else
			echo "   FAILED: $goos/$goarch" >&2
			FAILED="$FAILED $goos/$goarch"
			rm -f -- "$BASEDIR/$NAME-$goos-$goarch.zip"
		fi
	done
}

build_wasm() {
	local outdir="$BASEDIR/wasm"

	# Locate wasm_exec.js. Go moved it in 1.24: misc/wasm -> lib/wasm.
	local goroot wasm_exec candidate
	goroot="$(go env GOROOT)"
	wasm_exec=""
	for candidate in "$goroot/lib/wasm/wasm_exec.js" "$goroot/misc/wasm/wasm_exec.js"; do
		if [ -f "$candidate" ]; then
			wasm_exec="$candidate"
			break
		fi
	done
	if [ -z "$wasm_exec" ]; then
		echo "error: wasm_exec.js not found under $goroot (looked in lib/wasm and misc/wasm)." >&2
		echo "       Download the copy matching $(go env GOVERSION) from the golang/go repo." >&2
		FAILED="$FAILED js/wasm"
		return
	fi

	rm -rf -- "$outdir"
	mkdir -p -- "$outdir"

	echo ">> building $PKG for js/wasm"
	if ! GOOS=js GOARCH=wasm \
		go build -trimpath ${LDFLAGS:+-ldflags "$LDFLAGS"} -o "$outdir/game.wasm" "$PKG"; then
		echo "   FAILED: js/wasm" >&2
		FAILED="$FAILED js/wasm"
		return
	fi

	echo ">> copying $(basename "$wasm_exec") ($(go env GOVERSION))"
	cp -- "$wasm_exec" "$outdir/wasm_exec.js"

	echo ">> writing main.html"
	cat >"$outdir/main.html" <<HTMLEOF
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>$TITLE</title>
<style>html,body{margin:0;padding:0;overflow:hidden}</style>
</head>
<body>
<script src="wasm_exec.js"></script>
<script>
const go = new Go();
WebAssembly.instantiateStreaming(fetch("./game.wasm"), go.importObject).then(result => {
	go.run(result.instance);
});
</script>
</body>
</html>
HTMLEOF

	echo ">> writing index.html"
	cat >"$outdir/index.html" <<HTMLEOF
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>$TITLE</title>
<style>
    html, body {
        margin: 0;
        padding: 0;
        overflow: hidden;
        width: 100%;
        height: 100%;
        background-color: #000; /* Keeps the background dark during transitions */
    }
    iframe {
        width: 100%;
        height: 100%;
        border: none;
        display: block;
    }

</style>
</head>
<body>
    <iframe src="main.html" allow="autoplay; fullscreen"></iframe>
</body>
</html>
HTMLEOF

	echo "   $outdir/game.wasm ($(du -h "$outdir/game.wasm" | cut -f1))"
	make_zip "$BASEDIR/$NAME-wasm.zip" "$outdir"
}

if [ "$TARGET" = "wasm" ]; then
	build_wasm
else
	build_desktop "$TARGET"
fi

echo
if [ -n "$SKIPPED" ]; then
	echo "skipped (need a matching host for cgo):$SKIPPED"
fi
if [ -n "$FAILED" ]; then
	echo "done with failures:$FAILED"
	echo "output: $BASEDIR"
	exit 1
fi
echo "done: $BASEDIR"
if [ "$TARGET" = "wasm" ]; then
	echo "serve $BASEDIR/wasm over HTTP — .wasm must be sent as application/wasm."
fi