#!/usr/bin/env bash

set -e -x

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