#!/bin/bash

indir="ops-updater-https-source"

if [ $# -ne 1 ]
then
    echo "Usage: ./pack.sh v1.0.3"
else
    outfile="ops-updater-https-$1"
    cp -r ./$indir ./$outfile
    rm -rf $outfile/.git
    rm -rf $outfile/.gitignore
    rm -rf $outfile/.DS_STORE
    tar -zcvf $outfile.tar.gz $outfile
    rm -rf $outfile
fi
