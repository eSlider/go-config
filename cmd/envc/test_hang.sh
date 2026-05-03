#!/usr/bin/env bash
# Reproduce and explain envc "hangs": default input is stdin (-); io.ReadAll
# blocks until EOF. A TTY waits for Ctrl-D; /dev/zero never EOF (infinite read).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

ENVC="${TMPDIR:-/tmp}/envc-hang-test.$$"
cleanup() { rm -f "$ENVC"; }
trap cleanup EXIT

echo "== build envc -> $ENVC"
go build -o "$ENVC" ./cmd/envc

echo
echo "== 1) expected to finish quickly (YAML on stdin, explicit -)"
printf 'k: 1\n' | "$ENVC" convert --from yaml --to json --input - --output - | head -c 200
echo
echo "   (ok)"

echo
echo "== 2) empty stdin should finish quickly (EOF immediately)"
if ! out="$( "$ENVC" convert --from yaml --to json --input - --output - </dev/null 2>&1 )"; then
  echo "   exit non-zero (often empty YAML parse); stderr/out:"
  echo "$out" | sed 's/^/   | /'
  echo "   (ok — no hang, process terminated)"
else
  echo "$out" | sed 's/^/   /'
  echo "   (ok — no hang)"
fi

echo
echo "== 3) infinite stdin: should hit wall-clock limit (GNU timeout -> 124)"
if ! command -v timeout >/dev/null; then
  echo "   skip: no 'timeout' in PATH"
else
  set +e
  timeout -v 2s "$ENVC" convert --from yaml --to json --input - --output - </dev/zero >/dev/null 2>&1
  ec=$?
  set -e
  if [[ "$ec" -eq 124 ]]; then
    echo "   timeout fired (exit 124) — binary was blocked in read (no EOF on stdin)."
  elif [[ "$ec" -eq 137 ]] || [[ "$ec" -eq 143 ]]; then
    echo "   killed by signal (exit $ec) — still indicates blocked read until SIGKILL/SIGTERM."
  else
    echo "   unexpected exit $ec (expected 124 from GNU timeout on hang)"
    exit 1
  fi
fi

echo
echo "== why it happens"
cat <<'EOF'
   convert/get default to reading the input source once with io.ReadAll.
   For --input - (convert) or positional "-" (get), that is os.Stdin until EOF.

   - Interactive shell, no pipe: stdin is a TTY — blocked until you press Ctrl-D
     or pipe data in.
   - /dev/zero (or any endless reader): never EOF — blocked until killed or OOM.

   merge uses the same pattern for non-"-" inputs via bytesutil; URL fetches use
   a 30s HTTP client timeout (see internal/source.URL).

   Fix for users: pass a file or URL (--input path), or pipe finite data into stdin.
EOF

echo
echo "== done"
