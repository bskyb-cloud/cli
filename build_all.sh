#!/usr/bin/env bash

# run bin/replace-sha and add nimbus to Version in app_constants.go: Version = "6.22.2-nimbus-e480663"
# change the version

set -e -x

export VERSION="v6.22.2"
export OUTDIR=$PWD/out

GOARCH=amd64 GOOS=windows ./bin/build && cp $OUTDIR/cf $OUTDIR/cf-windows-amd64.exe
GOARCH=386 GOOS=windows ./bin/build && cp $OUTDIR/cf $OUTDIR/cf-windows-386.exe
GOARCH=amd64 GOOS=linux ./bin/build  && cp $OUTDIR/cf $OUTDIR/cf-linux-amd64
GOARCH=386 GOOS=linux ./bin/build  && cp $OUTDIR/cf $OUTDIR/cf-linux-386
GOARCH=amd64 GOOS=darwin ./bin/build  && cp $OUTDIR/cf $OUTDIR/cf-darwin-amd64

cd $OUTDIR

tar czf "cf-linux-amd64-$VERSION.tar.gz" cf-linux-amd64
tar czf "cf-linux-386-$VERSION.tar.gz" cf-linux-386
tar czf "cf-darwin-amd64-$VERSION.tar.gz" cf-darwin-amd64
zip "cf-windows-386-$VERSION.zip" cf-windows-386.exe
zip "cf-windows-amd64-$VERSION.zip" cf-windows-amd64.exe