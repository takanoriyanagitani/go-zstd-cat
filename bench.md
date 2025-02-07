# Simple Benchmark

## Encoder

| tool   | level   | rate     | size | ratio(rate) | ratio(size) |
|:------:|:-------:|:--------:|:----:|:-----------:|:-----------:|
| gozstd | Fast    | 675 MB/s | 900M | 173%        | 0.91x       |
| gozstd | Default | 644 MB/s | 895M | 165%        | 0.91x       |
| gozstd | Better  | 384 MB/s | 834M |  98%        | 0.85x       |
| gozstd | Best    |  93 MB/s | 718M |  24%        | 0.72x       |
| gzip   | Best    |  29 MB/s | 766M |   7%        | 0.78x       |
| gzip   | Fast    | 391 MB/s | 984M | (100%)      | (1x)        |

## Decoder

| tool    | level   | rate     | size | ratio(rate) |
|:-------:|:-------:|:--------:|:----:|:-----------:|
| (cat)   | -       | 4.0 GB/s | 6.1G | (200%)      |
| gozscat | Fast    | 1.3 GB/s | 900M |  65%        |
| gozscat | Default | 1.2 GB/s | 895M |  60%        |
| gozscat | Better  | 1.4 GB/s | 834M |  70%        |
| gozscat | Best    | 1.5 GB/s | 718M |  75%        |
| gzip    | Best    | 1.9 GB/s | 766M |  95%        |
| gzip    | Fast    | 2.0 GB/s | 984M | (100%)      |

