package zscat

import (
	"database/sql"
)

type DecodeConfig struct {
	Concurrency sql.Null[int]
}

var DecodeConfigDefault DecodeConfig = DecodeConfig{
	Concurrency: sql.Null[int]{V: 0, Valid: false},
}

type EncodeLevel int

const (
	EncodeLevelUnspecified EncodeLevel = iota
	EncodeLevelFast        EncodeLevel = iota
	EncodeLevelDefault     EncodeLevel = iota
	EncodeLevelBetter      EncodeLevel = iota
	EncodeLevelBest        EncodeLevel = iota
)

var EncodeLevelMap map[string]EncodeLevel = map[string]EncodeLevel{
	"Fast":    EncodeLevelFast,
	"Default": EncodeLevelDefault,
	"Better":  EncodeLevelBetter,
	"Best":    EncodeLevelBest,
}

func EncodeLevelFromStr(s string) EncodeLevel {
	val, found := EncodeLevelMap[s]
	switch found {
	case true:
		return val
	default:
		return EncodeLevelUnspecified
	}
}

type EncodeConfig struct {
	Concurrency sql.Null[int]
	EncodeLevel
}

var EncodeConfigDefault EncodeConfig = EncodeConfig{
	Concurrency: sql.Null[int]{V: 0, Valid: false},
	EncodeLevel: EncodeLevelDefault,
}
