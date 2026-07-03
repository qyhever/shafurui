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
  -iname "*.mp4" -o \
  -iname "*.mov" -o \
  -iname "*.m4v" -o \
  -iname "*.webm" \
\) -print0 |
while IFS= read -r -d '' video; do
  dir="$(dirname "$video")"
  filename="$(basename "$video")"
  base="${filename%.*}"
  cover="$dir/$base.jpg"

  if [ -f "$cover" ]; then
    echo "skip exists: $cover"
    continue
  fi

  echo "generate: $video -> $cover"

  if ! ffmpeg -hide_banner -loglevel error -y \
    -ss 00:00:01 \
    -i "$video" \
    -frames:v 1 \
    -q:v 3 \
    "$cover" </dev/null; then
    rm -f "$cover"
    echo "failed: $video"
  fi
done
