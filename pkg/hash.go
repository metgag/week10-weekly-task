package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashParams struct {
	Memory  uint32
	Time    uint32 // iterations
	Thread  uint8  // parallelism
	SaltLen uint32
	KeyLen  uint32
}

func NewHashParams() *HashParams {
	return &HashParams{}
}

func (h *HashParams) UseRecommended() {
	h.Memory = 64 * 1024
	h.Time = 3
	h.Thread = 2
	h.SaltLen = 16
	h.KeyLen = 32
}

// genHash
func (h *HashParams) GenerateFromPassword(password string) (string, error) {
	salt, err := h.generateRandomBytes()
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, h.Time, h.Memory, h.Thread, h.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.Memory, h.Time, h.Thread, b64Salt, b64Hash)

	return encodedHash, nil
}

// genSalt
func (h *HashParams) generateRandomBytes() ([]byte, error) {
	// salt
	b := make([]byte, h.SaltLen)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (h *HashParams) ComparePasswordAndHash(password, encodedHash string) (isMatch bool, err error) {
	salt, hash, err := h.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, h.Time, h.Memory, h.Thread, h.KeyLen)

	if subtle.ConstantTimeCompare(hash, otherHash) == 0 {
		return false, nil
	}

	return true, nil
}

func (h *HashParams) decodeHash(encodedHash string) (salt []byte, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, err
	}
	if vals[1] != "argon2id" {
		return nil, nil, err
	}

	var version int
	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, err
	}

	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &h.Memory, &h.Time, &h.Thread); err != nil {
		return nil, nil, err
	}

	getSalt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}
	h.SaltLen = uint32(len(getSalt))

	getHash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, err
	}
	h.KeyLen = uint32(len(getHash))

	return getSalt, getHash, nil
}
