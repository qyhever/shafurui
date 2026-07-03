#!/usr/bin/env bash

set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <video-root-dir>"
  exit 1
fi

VIDEO_ROOT="$1"

if [ ! -d "$VIDEO_ROOT" ]; then
  echo "video root dir not found: $VIDEO_ROOT"
  exit 1
fi

if ! command -v ffmpeg >/dev/null 2>&1; then
  echo "ffmpeg is required. Install it first."
  exit 1
fi

find "$VIDEO_ROOT" -type f \( \
  -iname "*.mov" -o \
  -iname "*.m4v" \
\) -print0 |
while IFS= read -r -d '' video; do
  dir="$(dirname "$video")"
  filename="$(basename "$video")"
  base="${filename%.*}"
  output="$dir/$base.mp4"
  tmp_output="$output.tmp.mp4"

  if [ -f "$output" ]; then
    echo "skip exists: $output"
    continue
  fi

  echo "remux: $video -> $output"

  rm -f "$tmp_output"

  if ffmpeg -hide_banner -loglevel error -y \
    -i "$video" \
    -map 0:v:0 \
    -map 0:a:0? \
    -c copy \
    -movflags +faststart \
    -map_metadata 0 \
    "$tmp_output" </dev/null; then
    mv "$tmp_output" "$output"
    rm -f "$video"
  else
    rm -f "$tmp_output"
    echo "failed: $video"
  fi
done
