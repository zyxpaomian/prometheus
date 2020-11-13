#!/usr/bin/env bash
#
# Build React web UI.
# Run from repository root.
set -e
set -u

echo $0

if ! [[ "$0" =~ "D:\Coding\myprometheus\scripts\build_react_app.sh" ]]; then
	echo "must be run from repository root"
	exit 255
fi

cd web/ui/react-app

echo "building React app"
PUBLIC_URL=. yarn build
rm -rf ../static/react
mv build ../static/react
