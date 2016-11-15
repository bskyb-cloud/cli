#!/usr/bin/env bash

# run bin/replace-sha and add nimbus to Version in app_constants.go: Version = "6.22.2-nimbus-e480663"

set -e -x

export OUTDIR=$PWD/out

GOARCH=amd64 GOOS=windows ./bin/build && cp $OUTDIR/cf $OUTDIR/cf-windows-amd64.exe
GOARCH=386 GOOS=windows ./bin/build && cp $OUTDIR/cf $OUTDIR/cf-windows-386.exe
GOARCH=amd64 GOOS=linux ./bin/build  && cp $OUTDIR/cf $OUTDIR/cf-linux-amd64
GOARCH=386 GOOS=linux ./bin/build  && cp $OUTDIR/cf $OUTDIR/cf-linux-386
GOARCH=amd64 GOOS=darwin ./bin/build  && cp $OUTDIR/cf $OUTDIR/cf-darwin-amd64
