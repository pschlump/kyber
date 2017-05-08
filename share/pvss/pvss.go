// Package pvss implements public verifiable secret sharing as introduced in
// "A Simple Publicly Verifiable Secret Sharing Scheme and its Application to
// Electronic Voting" by Berry Schoenmakers. In comparison to regular verifiable
// secret sharing schemes, PVSS enables any third party to verify shares
// distributed by a dealer using zero-knowledge proofs. PVSS runs in three steps:
//  1. The dealer creates a list of encrypted public verifiable shares using
//     EncShares() and distributes them to the trustees.
//  2. Upon the announcement that the secret should be released, each trustee
//     uses DecShare() to first verify and, if valid, decrypt his share.
//  3. Once a threshold of decrypted shares has been released, anyone can
//     verify them and, if enough shares are valid, recover the shared secret
//     using RecoverSecret().
// For concrete examples see pvss_test.go.
package pvss

import (
	"errors"
	"hash"

	"github.com/dedis/crypto"
	"github.com/dedis/crypto/proof"
	"github.com/dedis/crypto/share"
	"github.com/dedis/crypto/util/random"
)

// NOTE. Here, we need need the full suite definition for using the proof/
// package while we only use Hash.
// 1. the pkg has still access to lot more than what the package requires (so
// the opposite of what we wanted to achieve)
// 2. Again, full suite redeclaration
// 3. Possibly troublesome inconsistencies: why would vss use a different suite
// than share/poly.go, while vss USES share/poly.go as a building block.
type Suite interface {
	crypto.Group
	Hash() hash.Hash
	Cipher(key []byte, options ...interface{}) crypto.Cipher
	crypto.Encoding
}

// Some error definitions.
var errorTooFewShares = errors.New("not enough shares to recover secret")
var errorDifferentLengths = errors.New("inputs of different lengths")
var errorEncVerification = errors.New("verification of encrypted share failed")
var errorDecVerification = errors.New("verification of decrypted share failed")

// PubVerShare is a public verifiable share.
type PubVerShare struct {
	S share.PubShare  // Share
	P proof.DLEQProof // Proof
}

// EncShares creates a list of encrypted publicly verifiable PVSS shares for
// the given secret and the list of public keys X using the sharing threshold
// t and the base point H. The function returns the list of shares and the
// public commitment polynomial.
func EncShares(suite Suite, H crypto.Point, X []crypto.Point, secret crypto.Scalar, t int) ([]*PubVerShare, *share.PubPoly, error) {
	n := len(X)
	encShares := make([]*PubVerShare, n)

	// Create secret sharing polynomial
	priPoly := share.NewPriPoly(suite, t, secret, random.Stream)

	// Create secret set of shares
	priShares := priPoly.Shares(n)

	// Create public polynomial commitments with respect to basis H
	pubPoly := priPoly.Commit(H)

	// Prepare data for encryption consistency proofs ...
	indices := make([]int, n)
	values := make([]crypto.Scalar, n)
	HS := make([]crypto.Point, n)
	for i := 0; i < n; i++ {
		indices[i] = priShares[i].I
		values[i] = priShares[i].V
		HS[i] = H
	}

	// Create NIZK discrete-logarithm equality proofs
	proofs, _, sX, err := proof.NewDLEQProofBatch(suite, HS, X, values)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < n; i++ {
		ps := &share.PubShare{indices[i], sX[i]}
		encShares[i] = &PubVerShare{*ps, *proofs[i]}
	}

	return encShares, pubPoly, nil
}

// VerifyEncShare checks that the encrypted share sX satisfies
// log_{H}(sH) == log_{X}(sX) where sH is the public commitment computed by
// evaluating the public commitment polynomial at the encrypted share's index i.
func VerifyEncShare(suite Suite, H crypto.Point, X crypto.Point, sH crypto.Point, encShare *PubVerShare) error {
	if err := encShare.P.Verify(suite, H, X, sH, encShare.S.V); err != nil {
		return errorEncVerification
	}
	return nil
}

// VerifyEncShareBatch provides the same functionality as VerifyEncShare but for
// slices of encrypted shares. The function returns the valid encrypted shares
// together with the corresponding public keys.
func VerifyEncShareBatch(suite Suite, H crypto.Point, X []crypto.Point, sH []crypto.Point, encShares []*PubVerShare) ([]crypto.Point, []*PubVerShare, error) {
	if len(X) != len(sH) || len(sH) != len(encShares) {
		return nil, nil, errorDifferentLengths
	}
	var K []crypto.Point // good public keys
	var E []*PubVerShare // good encrypted shares
	for i := 0; i < len(X); i++ {
		if err := VerifyEncShare(suite, H, X[i], sH[i], encShares[i]); err == nil {
			K = append(K, X[i])
			E = append(E, encShares[i])
		}
	}
	return K, E, nil
}

// DecShare first verifies the encrypted share against the encryption
// consistency proof and, if valid, decrypts it and creates a decryption
// consistency proof.
func DecShare(suite Suite, H crypto.Point, X crypto.Point, sH crypto.Point, x crypto.Scalar, encShare *PubVerShare) (*PubVerShare, error) {
	if err := VerifyEncShare(suite, H, X, sH, encShare); err != nil {
		return nil, err
	}
	G := suite.Point().Base()
	V := suite.Point().Mul(encShare.S.V, suite.Scalar().Inv(x)) // decryption: x^{-1} * (xS)
	ps := &share.PubShare{encShare.S.I, V}
	P, _, _, err := proof.NewDLEQProof(suite, G, V, x)
	if err != nil {
		return nil, err
	}
	return &PubVerShare{*ps, *P}, nil
}

// DecShareBatch provides the same functionality as DecShare but for slices of
// encrypted shares. The function returns the valid encrypted and decrypted
// shares as well as the corresponding public keys.
func DecShareBatch(suite Suite, H crypto.Point, X []crypto.Point, sH []crypto.Point, x crypto.Scalar, encShares []*PubVerShare) ([]crypto.Point, []*PubVerShare, []*PubVerShare, error) {
	if len(X) != len(sH) || len(sH) != len(encShares) {
		return nil, nil, nil, errorDifferentLengths
	}
	var K []crypto.Point // good public keys
	var E []*PubVerShare // good encrypted shares
	var D []*PubVerShare // good decrypted shares
	for i := 0; i < len(encShares); i++ {
		if ds, err := DecShare(suite, H, X[i], sH[i], x, encShares[i]); err == nil {
			K = append(K, X[i])
			E = append(E, encShares[i])
			D = append(D, ds)
		}
	}
	return K, E, D, nil
}

// VerifyDecShare checks that the decrypted share sG satisfies
// log_{G}(X) == log_{sG}(sX). Note that X = xG and sX = s(xG) = x(sG).
func VerifyDecShare(suite Suite, G crypto.Point, X crypto.Point, encShare *PubVerShare, decShare *PubVerShare) error {
	if err := decShare.P.Verify(suite, G, decShare.S.V, X, encShare.S.V); err != nil {
		return errorDecVerification
	}
	return nil
}

// VerifyDecShareBatch provides the same functionality as VerifyDecShare but for
// slices of decrypted shares. The function returns the the valid decrypted shares.
func VerifyDecShareBatch(suite Suite, G crypto.Point, X []crypto.Point, encShares []*PubVerShare, decShares []*PubVerShare) ([]*PubVerShare, error) {
	if len(X) != len(encShares) || len(encShares) != len(decShares) {
		return nil, errorDifferentLengths
	}
	var D []*PubVerShare // good decrypted shares
	for i := 0; i < len(X); i++ {
		if err := VerifyDecShare(suite, G, X[i], encShares[i], decShares[i]); err == nil {
			D = append(D, decShares[i])
		}
	}
	return D, nil
}

// RecoverSecret first verifies the given decrypted shares against their
// decryption consistency proofs and then tries to recover the shared secret.
func RecoverSecret(suite Suite, G crypto.Point, X []crypto.Point, encShares []*PubVerShare, decShares []*PubVerShare, t int, n int) (crypto.Point, error) {
	D, err := VerifyDecShareBatch(suite, G, X, encShares, decShares)
	if err != nil {
		return nil, err
	}
	if len(D) < t {
		return nil, errorTooFewShares
	}
	var shares []*share.PubShare
	for _, s := range D {
		shares = append(shares, &s.S)
	}
	return share.RecoverCommit(suite, shares, t, n)
}