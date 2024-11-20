package algorithm

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type ArgonHashing struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLen      uint32
	saltLen     uint32
}

func Argon2(memory, iterations, keyLen, saltLen uint32, parallelism uint8) *ArgonHashing {
	return &ArgonHashing{
		memory:      memory,
		iterations:  iterations,
		parallelism: parallelism,
		keyLen:      keyLen,
		saltLen:     saltLen,
	}
}

func (p *ArgonHashing) Generate(str string) (encodedHash string, err error) {
	salt, err := p.generateRandomBytes(p.saltLen)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(str),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLen,
	)

	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.memory,
		p.iterations,
		p.parallelism,
		base64Salt,
		base64Hash,
	)

	return encodedHash, nil
}

func (p *ArgonHashing) Compare(password, encodedHash string) (match bool, err error) {
	hashing, salt, hash, err := p.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		hashing.iterations,
		hashing.memory,
		hashing.parallelism,
		hashing.keyLen,
	)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (p *ArgonHashing) decodeHash(encodedHash string) (hashing *ArgonHashing, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	hashing = &ArgonHashing{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &hashing.memory, &hashing.iterations, &hashing.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	hashing.saltLen = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	hashing.keyLen = uint32(len(hash))

	return hashing, salt, hash, nil
}

func (h *ArgonHashing) generateRandomBytes(length uint32) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
