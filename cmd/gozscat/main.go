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
	dk "github.com/takanoriyanagitani/go-zstd-cat/zstd/dec/klauspost"
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
	envValByKey("ENV_ZSTD_DEC_CONCURRENCY"),
	Lift(strconv.Atoi),
).Or(Of(0))

var decodeCfg IO[zc.DecodeConfig] = Bind(
	concurrency,
	Lift(func(i int) (zc.DecodeConfig, error) {
		return zc.DecodeConfig{
			Concurrency: sql.Null[int]{V: i, Valid: true},
		}, nil
	}),
)

var cfg IO[dk.Config] = Bind(
	decodeCfg,
	Lift(func(dc zc.DecodeConfig) (dk.Config, error) {
		return dk.Config(dc), nil
	}),
)

var stdin2stdout IO[Void] = Bind(
	cfg,
	func(dc dk.Config) IO[Void] {
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
