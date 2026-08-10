#!/usr/bin/env bash
#
# wasm.sh — build a Go/Ebitengine package for wasm.
#
# Outputs:
#   <outdir>/game.wasm
#   <outdir>/wasm_exec.js
#   <outdir>/main.html    (loader; the actual game frame)
#   <outdir>/index.html   (host page; iframes main.html)
#
set -euo pipefail

WIDTH=640
HEIGHT=480
TITLE=""
OUTDIR=""

usage() {
	cat <<'EOF'
Usage: build-wasm.sh [options] <package>

  <package>   Go package to build, e.g. ./example/snd

Options:
  -o DIR      output directory (default: bin/<package path>)
  -W N        iframe width in px (default: 640)
  -H N        iframe height in px (default: 480)
  -t TEXT     page title (default: package directory name)
  -h          show this help

Examples:
  ./tools/build/wasm.sh ./example/snd
  ./tools/build/wasm.sh -o dist -W 800 -H 600 -t "My Demo" ./cmd/game
EOF
}

while getopts ":o:W:H:t:ch" opt; do
	case "$opt" in
	o) OUTDIR="$OPTARG" ;;
	W) WIDTH="$OPTARG" ;;
	H) HEIGHT="$OPTARG" ;;
	t) TITLE="$OPTARG" ;;
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

if [ $# -ne 1 ]; then
	echo "error: expected exactly one package argument" >&2
	usage >&2
	exit 2
fi
PKG="$1"

command -v go >/dev/null 2>&1 || {
	echo "error: go not found in PATH" >&2
	exit 1
}

# Derive a sensible output directory and title from the package path.
REL="${PKG#./}"
REL="${REL%/}"
if [ -z "$REL" ] || [ "$REL" = "." ]; then
	REL="$(basename "$PWD")"
fi
[ -n "$OUTDIR" ] || OUTDIR="bin/$REL"
[ -n "$TITLE" ] || TITLE="$(basename "$REL")"

# Locate wasm_exec.js. Go moved it in 1.24: misc/wasm -> lib/wasm.
GOROOT="$(go env GOROOT)"
WASM_EXEC=""
for candidate in "$GOROOT/lib/wasm/wasm_exec.js" "$GOROOT/misc/wasm/wasm_exec.js"; do
	if [ -f "$candidate" ]; then
		WASM_EXEC="$candidate"
		break
	fi
done
if [ -z "$WASM_EXEC" ]; then
	echo "error: wasm_exec.js not found under $GOROOT (looked in lib/wasm and misc/wasm)." >&2
	echo "       Download the copy matching $(go env GOVERSION) from the golang/go repo." >&2
	exit 1
fi

if [ -d "$OUTDIR" ]; then
	echo ">> cleaning $OUTDIR"
	rm -rf -- "$OUTDIR"
fi
mkdir -p -- "$OUTDIR"

echo ">> building $PKG for js/wasm"
env GOOS=js GOARCH=wasm go build -o "$OUTDIR/game.wasm" "$PKG"

echo ">> copying $(basename "$WASM_EXEC") ($(go env GOVERSION))"
cp -- "$WASM_EXEC" "$OUTDIR/wasm_exec.js"

echo ">> writing main.html"
cat >"$OUTDIR/main.html" <<EOF
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
EOF

echo ">> writing index.html"
cat >"$OUTDIR/index.html" <<EOF
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>$TITLE</title>
</head>
<body>
<iframe src="main.html" allow="autoplay" width="$WIDTH" height="$HEIGHT" frameborder="0"></iframe>
</body>
</html>
EOF

SIZE="$(du -h "$OUTDIR/game.wasm" | cut -f1)"
echo
echo "done: $OUTDIR (game.wasm $SIZE)"
