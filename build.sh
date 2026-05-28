#!/bin/bash
set -euo pipefail

APP_NAME="StickerDownloadBot"
BUILD_TIME=$(date "+%Y-%m-%d %H:%M:%S")

echo "Building $APP_NAME ..."
echo "buildTime: $BUILD_TIME"

go build -trimpath -ldflags="-s -w -X 'main.buildTime=$BUILD_TIME'" -o "$APP_NAME" .

echo "Done: $APP_NAME"
