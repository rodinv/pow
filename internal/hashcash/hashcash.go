package hashcash

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/rodinv/errors"
)

const (
	randomBytes      = 8
	timeFormat       = "2006-01-02"
	hashcashV1Length = 8
)

// ProofOfWork base hashcash algorithm struct
type ProofOfWork struct {
	version int
	bits    int64
	target  *big.Int
	storage StorageProvider
}

// StorageProvider provides work with hash storage
type StorageProvider interface {
	IsSpent(key string) bool
}

func New(b int64, storage StorageProvider) (*ProofOfWork, error) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b))

	return &ProofOfWork{
		version: 1,
		bits:    b,
		target:  target,
		storage: storage,
	}, nil
}

// GetHeader gets the header by resource
func (p *ProofOfWork) GetHeader(resource string) (string, error) {
	token := make([]byte, randomBytes)
	_, err := rand.Read(token)
	if err != nil {
		return "", errors.Wrap(err, "getting random bytes")
	}

	return fmt.Sprintf("%d:%d:%s:%s::%s:",
		p.version,
		p.bits,
		time.Now().UTC().Format(timeFormat),
		resource,
		base64.StdEncoding.EncodeToString(token),
	), nil
}

// Compute computes hash
func (p *ProofOfWork) Compute(header string) string {
	for nonce := int64(0); nonce < math.MaxInt64; nonce++ {
		data := fmt.Sprintf("%s%d", header, nonce)

		hash := sha256.Sum256([]byte(data))
		hashInt := big.NewInt(0)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(p.target) == -1 {
			return data
		}
	}

	return ""
}

// Verify verifies incoming hash
func (p *ProofOfWork) Verify(data string) error {
	vals := strings.Split(data, ":")
	if len(vals) != hashcashV1Length {
		return errors.Errorf("wrong hashcash len %d, hash %s", len(vals), data)
	}

	// verify hash
	hash := sha256.Sum256([]byte(data))
	hashInt := big.NewInt(0)
	hashInt.SetBytes(hash[:])

	if hashInt.Cmp(p.target) != -1 {
		return errors.New("wrong hash")
	}

	// verify time
	date, err := time.Parse(timeFormat, vals[2])
	if err != nil {
		return errors.Errorf("wrong date %s", vals[2])
	}
	if date.After(time.Now().AddDate(0, 0, 2)) {
		return errors.New("date is too far into the future")
	}

	// verify resource
	resource := vals[3]
	if net.ParseIP(resource) == nil {
		return errors.New("wrong resource")
	}

	// check store
	if p.storage.IsSpent(data) {
		return errors.New("hash has already been used up")
	}

	return nil
}

// GetBitsFromHeader gets bits from incoming header
func GetBitsFromHeader(header string) (int64, error) {
	vals := strings.Split(header, ":")
	if len(vals) != hashcashV1Length {
		return -1, errors.Errorf("wrong hashcash len %d", len(vals))
	}

	bits, err := strconv.ParseInt(vals[1], 10, 0)
	if err != nil {
		return -1, errors.Wrapf(err, "can't parse %d to int", vals[1])
	}

	return bits, nil
}
