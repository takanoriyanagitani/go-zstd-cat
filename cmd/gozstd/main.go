package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	zc "github.com/takanoriyanagitani/go-zstd-cat"
	. "github.com/takanoriyanagitani/go-zstd-cat/util"
	ek "github.com/takanoriyanagitani/go-zstd-cat/zstd/enc/klauspost"
)

func envValByKey(key string) IO[string] {
	return func(_ context.Context) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("env var %s missing", key)
		}
	}
}

var concurrency IO[int] = Bind(
	envValByKey("ENV_ZSTD_ENC_CONCURRENCY"),
	Lift(strconv.Atoi),
)

var optcon IO[sql.Null[int]] = Ok(concurrency)

var encodeLevel IO[zc.EncodeLevel] = Bind(
	envValByKey("ENV_ENCODE_LEVEL"),
	Lift(func(level string) (zc.EncodeLevel, error) {
		return zc.EncodeLevelFromStr(level), nil
	}),
).Or(Of(zc.EncodeLevelDefault))

var encodeCfg IO[zc.EncodeConfig] = Bind(
	optcon,
	func(oc sql.Null[int]) IO[zc.EncodeConfig] {
		return Bind(
			encodeLevel,
			Lift(func(lv zc.EncodeLevel) (zc.EncodeConfig, error) {
				return zc.EncodeConfig{
					Concurrency: oc,
					EncodeLevel: lv,
				}, nil
			}),
		)
	},
)

var cfg IO[ek.Config] = Bind(
	encodeCfg,
	Lift(func(dc zc.EncodeConfig) (ek.Config, error) {
		return ek.Config(dc), nil
	}),
)

var stdin2stdout IO[Void] = Bind(
	cfg,
	func(dc ek.Config) IO[Void] {
		return Bind(
			dc.ToStdinToStdout(),
			Lift(func(_ int64) (Void, error) { return Empty, nil }),
		)
	},
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return stdin2stdout(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
