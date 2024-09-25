package hashing

import (
	"errors"

	"github.com/mrrizkin/omniscan/system/config"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type Hashing interface {
	GenerateHash(str string) (string, error)
	CompareHash(str, hash string) (bool, error)
}

func Argon2(config config.Config) Hashing {
	return newArgon(
		uint32(config.HASH_MEMORY),
		uint32(config.HASH_ITERATIONS),
		uint32(config.HASH_KEY_LEN),
		uint32(config.HASH_SALT_LEN),
		uint8(config.HASH_PARALLELISM),
	)
}
