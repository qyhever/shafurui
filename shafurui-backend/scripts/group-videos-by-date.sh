#!/usr/bin/env bash

set -euo pipefail

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <source-dir> <output-dir>"
  exit 1
fi

SOURCE_DIR="$1"
OUTPUT_DIR="$2"

if [ ! -d "$SOURCE_DIR" ]; then
  echo "source dir not found: $SOURCE_DIR"
  exit 1
fi

if ! command -v exiftool >/dev/null 2>&1; then
  echo "exiftool is required. Install it first."
  exit 1
fi

mkdir -p "$OUTPUT_DIR"

mtime_date() {
  local file="$1"

  case "$(uname -s 2>/dev/null || echo unknown)" in
    Darwin|FreeBSD)
      stat -f "%Sm" -t "%Y-%m-%d" "$file"
      ;;
    *)
      date -r "$file" "+%Y-%m-%d"
      ;;
  esac
}

find "$SOURCE_DIR" -type f \( \
  -iname "*.mp4" -o \
  -iname "*.mov" -o \
  -iname "*.m4v" -o \
  -iname "*.webm" \
\) -print0 |
while IFS= read -r -d '' video; do
  shot_date="$(
    exiftool -s3 -d "%Y-%m-%d" \
      -CreateDate \
      -MediaCreateDate \
      -TrackCreateDate \
      "$video" 2>/dev/null |
      sed -n '1p'
  )"

  if [ -z "$shot_date" ]; then
    shot_date="$(mtime_date "$video")"
  fi

  target_dir="$OUTPUT_DIR/$shot_date"
  target_file="$target_dir/$(basename "$video")"

  mkdir -p "$target_dir"

  if [ -e "$target_file" ]; then
    echo "skip exists: $target_file"
    continue
  fi

  cp -p "$video" "$target_file"
  echo "copy: $video -> $target_file"
done
