package dec

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"

	kz "github.com/klauspost/compress/zstd"
	zc "github.com/takanoriyanagitani/go-zstd-cat"
	. "github.com/takanoriyanagitani/go-zstd-cat/util"
)

func CopyZstdReader(zr *kz.Decoder, wtr io.Writer) (int64, error) {
	return io.Copy(wtr, zr)
}

func CopyReaderZstd(
	rdr io.Reader, wtr io.Writer, opts ...kz.DOption,
) (int64, error) {
	zrdr, e := kz.NewReader(rdr, opts...)
	if nil != e {
		return 0, e
	}
	defer zrdr.Close()

	return CopyZstdReader(zrdr, wtr)
}

type Config zc.DecodeConfig

func (c Config) ToConcurrency() (i int, valid bool) {
	switch c.Concurrency.Valid {
	case true:
		return c.Concurrency.V, 0 <= c.Concurrency.V
	default:
		return 0, false
	}
}

func (c Config) ToOpts() []kz.DOption {
	var ret []kz.DOption

	concurrency, valid := c.ToConcurrency()
	if valid {
		ret = append(ret, kz.WithDecoderConcurrency(concurrency))
	}
	return ret
}

func (c Config) Copy(rdr io.Reader, wtr io.Writer) (int64, error) {
	return CopyReaderZstd(rdr, wtr, c.ToOpts()...)
}

func (c Config) ToStdinToStdout() IO[int64] {
	return func(_ context.Context) (int64, error) {
		var br io.Reader = bufio.NewReader(os.Stdin)

		var bw *bufio.Writer = bufio.NewWriter(os.Stdout)
		i, e := c.Copy(br, bw)
		return i, errors.Join(e, bw.Flush())
	}
}
