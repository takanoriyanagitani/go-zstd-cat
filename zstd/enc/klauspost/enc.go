package dec

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"io"
	"os"

	kz "github.com/klauspost/compress/zstd"
	zc "github.com/takanoriyanagitani/go-zstd-cat"
	. "github.com/takanoriyanagitani/go-zstd-cat/util"
)

type Config zc.EncodeConfig

type Level zc.EncodeLevel

type LevelConversionMap map[zc.EncodeLevel]kz.EncoderLevel

var LevelConvMap LevelConversionMap = map[zc.EncodeLevel]kz.EncoderLevel{
	zc.EncodeLevelFast:    kz.SpeedFastest,
	zc.EncodeLevelDefault: kz.SpeedDefault,
	zc.EncodeLevelBetter:  kz.SpeedBetterCompression,
	zc.EncodeLevelBest:    kz.SpeedBestCompression,
}

func (l Level) Convert() sql.Null[kz.EncoderLevel] {
	var ret sql.Null[kz.EncoderLevel]

	mapd, found := LevelConvMap[zc.EncodeLevel(l)]
	if !found {
		return ret
	}

	ret.V = mapd
	ret.Valid = true
	return ret
}

func (c Config) ToConcurrency() (i int, valid bool) {
	switch c.Concurrency.Valid {
	case true:
		return c.Concurrency.V, 1 <= c.Concurrency.V
	default:
		return 0, false
	}
}

func (c Config) ToLevel() sql.Null[kz.EncoderLevel] {
	return Level(c.EncodeLevel).Convert()
}

func (c Config) ToOpts() []kz.EOption {
	var ret []kz.EOption

	var lvl sql.Null[kz.EncoderLevel] = c.ToLevel()
	if lvl.Valid {
		ret = append(ret, kz.WithEncoderLevel(lvl.V))
	}

	con, valid := c.ToConcurrency()
	if valid {
		ret = append(ret, kz.WithEncoderConcurrency(con))
	}

	return ret
}

func CopyToZstdWriter(rdr io.Reader, zw *kz.Encoder) (int64, error) {
	return io.Copy(zw, rdr)
}

func CopyToWriterZstd(
	rdr io.Reader,
	wtr io.Writer,
	opts ...kz.EOption,
) (int64, error) {
	enc, e := kz.NewWriter(wtr, opts...)
	if nil != e {
		return 0, e
	}
	i, e := CopyToZstdWriter(rdr, enc)
	return i, errors.Join(e, enc.Close())
}

func (c Config) Copy(rdr io.Reader, wtr io.Writer) (int64, error) {
	return CopyToWriterZstd(
		rdr,
		wtr,
		c.ToOpts()...,
	)
}

func (c Config) ToStdinToStdout() IO[int64] {
	return func(_ context.Context) (int64, error) {
		var br io.Reader = bufio.NewReader(os.Stdin)

		var bw *bufio.Writer = bufio.NewWriter(os.Stdout)
		i, e := c.Copy(br, bw)
		return i, errors.Join(e, bw.Flush())
	}
}
