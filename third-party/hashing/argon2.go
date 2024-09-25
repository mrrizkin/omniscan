package hashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLen      uint32
	saltLen     uint32
}

func newArgon(memory, iterations, keyLength, saltLength uint32, parallelism uint8) Hashing {
	return &argon{
		memory:      memory,
		iterations:  iterations,
		parallelism: parallelism,
		saltLen:     saltLength,
		keyLen:      keyLength,
	}
}

func (a *argon) GenerateHash(str string) (encodedHash string, err error) {
	salt, err := generateRandomBytes(a.saltLen)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(str),
		salt,
		a.iterations,
		a.memory,
		a.parallelism,
		a.keyLen,
	)

	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		a.memory,
		a.iterations,
		a.parallelism,
		base64Salt,
		base64Hash,
	)

	return encodedHash, nil
}

func (a *argon) CompareHash(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLen,
	)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *argon, salt, hash []byte, err error) {
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

	p = &argon{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLen = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLen = uint32(len(hash))

	return p, salt, hash, nil
}

func generateRandomBytes(length uint32) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
