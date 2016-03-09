#!/bin/bash

rm -f *.sha256
for f in $(ls); do shasum -a 256 $f > $f.sha256; done
