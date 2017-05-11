package key

import (
	"crypto/cipher"
	"encoding/base64"

	"github.com/dedis/kyber"
	"github.com/dedis/kyber/util/random"
)

// KeyPair represents a public/private keypair
// together with the ciphersuite the key was generated from.
type KeyPair struct {
	Suite  kyber.Suite  // Ciphersuite this keypair is for
	Public kyber.Point  // Public key
	Secret kyber.Scalar // Secret key
}

// NewKeyPair directly creates a secret/public key pair
func NewKeyPair(suite kyber.Suite) *KeyPair {
	kp := new(KeyPair)
	kp.Gen(suite, random.Stream)
	return kp
}

// Generate a fresh public/private keypair with the given ciphersuite,
// using a given source of cryptographic randomness.
func (p *KeyPair) Gen(suite kyber.Suite, random cipher.Stream) {
	p.Suite = suite
	p.Secret = suite.NewKey(random)
	p.Public = suite.Point().Mul(nil, p.Secret)
}

// PubId returns the base64-encoded HashId for this KeyPair's public key.
func (p *KeyPair) PubId() string {
	buf, _ := p.Public.MarshalBinary()
	hash := kyber.Sum(p.Suite, buf)
	return base64.RawURLEncoding.EncodeToString(hash)
}
