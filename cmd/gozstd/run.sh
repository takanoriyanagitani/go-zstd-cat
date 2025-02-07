#!/bin/sh

input=/usr/share/dict/words

cat "${input}" | cat                               | wc -c
cat "${input}" | ENV_ENCODE_LEVEL=Fast    ./gozstd | wc -c
cat "${input}" | ENV_ENCODE_LEVEL=Default ./gozstd | wc -c
cat "${input}" | ENV_ENCODE_LEVEL=        ./gozstd | wc -c
cat "${input}" | ENV_ENCODE_LEVEL=Better  ./gozstd | wc -c
cat "${input}" | ENV_ENCODE_LEVEL=Best    ./gozstd | wc -c
