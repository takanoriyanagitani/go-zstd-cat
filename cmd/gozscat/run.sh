#!/bin/sh

sample=/usr/share/dict/words
input=sample.d/sample.zst

geninput(){
	echo generating input...
	cat "${sample}" | ENV_ENCODE_LEVEL=Best ../gozstd/gozstd > "${input}"
}

test -f "${input}" || geninput

cat "${input}" | ./gozscat | wc -c
