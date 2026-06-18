#!/usr/bin/env bash
set -euo pipefail

URL="https://download.unity3d.com/download_unity/c04dd374db98/MacEditorTargetInstaller/UnitySetup-Mac-Support-for-Editor-5.3.6p8.pkg"
PKG="UnitySetup-Mac-Support-for-Editor-5.3.6p8.pkg"
WORKDIR="./Unity-5.3.6p8"
APP="./AmazingWorld.app"
CONFIG="./ServerConfig.xml"

COMPONENT="$WORKDIR/TargetSupport.pkg.tmp"
PAYLOAD="$COMPONENT/Payload"

if [ ! -f "$PKG" ]; then
    echo "Downloading $PKG..."
    curl -L --fail -o "$PKG" "$URL"
fi

rm -rf "$WORKDIR"

echo "Expanding package..."
pkgutil --expand "$PKG" "$WORKDIR"

if [ ! -f "$PAYLOAD" ]; then
    echo "Payload not found: $PAYLOAD" >&2
    exit 1
fi

echo "Extracting payload..."

(
    cd "$COMPONENT"

    rm -rf Payload.unpacked
    mkdir Payload.unpacked
    cd Payload.unpacked

    if ! ditto -x ../Payload . 2>/dev/null; then
        gunzip -dc ../Payload | cpio -idm
    fi
)

UNITY_PLAYER="$COMPONENT/Payload.unpacked/Variations/macosx64_nondevelopment_mono/UnityPlayer.app"

SRC_MONO="$UNITY_PLAYER/Contents/Frameworks/MonoEmbedRuntime/osx"
DST_MONO="$APP/Contents/Frameworks/MonoEmbedRuntime/osx"

SRC_PLAYER="$UNITY_PLAYER/Contents/MacOS/UnityPlayer"
DST_PLAYER="$APP/Contents/MacOS/AmazingWorld"

echo "Copying Mono runtime libraries..."
cp -f "$SRC_MONO"/*.dylib "$DST_MONO/"

echo "Replacing executable..."
cp -f "$SRC_PLAYER" "$DST_PLAYER"
chmod +x "$DST_PLAYER"

if [ ! -f "$CONFIG" ]; then
    echo "ServerConfig.xml not found: $CONFIG" >&2
    exit 1
fi

echo "Patching ServerConfig.xml..."
sed -i '' 's/user\.amazingworld\.com/springbay.amazingcore.org/g' "$CONFIG"

echo "Done."

