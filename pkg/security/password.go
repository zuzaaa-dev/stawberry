package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type param struct {
	memory   uint32
	time     uint32
	threads  uint8
	saltSize uint32
	keyLen   uint32
}

const (
	saltSize  = 16
	argonTime = 1
	threads   = 4
	keyLen    = 32
	memory    = 64 * 1024
)

func HashArgon2id(password string) (string, error) {
	salt, err := generateRandomBytes(saltSize)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, memory, threads, keyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, argonTime, threads, b64Salt, b64Hash,
	)

	return encodedHash, nil
}

func ComparePasswordAndArgon2id(password, encodedHash string) (bool, error) {
	p, salt, hash, err := decodeArgon2id(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threads, p.keyLen)

	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return []byte(""), err
	}
	return b, nil
}

func decodeArgon2id(encodedHash string) (param, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return param{}, nil, nil, errors.New("invalid encoded hash")
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return param{}, nil, nil, err
	}

	if version != argon2.Version {
		return param{}, nil, nil, errors.New("incompatible hash version")
	}

	var p param
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.time, &p.threads)
	if err != nil {
		return param{}, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return param{}, nil, nil, err
	}
	p.saltSize = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return param{}, nil, nil, err
	}
	p.keyLen = uint32(len(hash))

	return p, salt, hash, nil
}
