// +build vartime

package group

import (
	"github.com/dedis/kyber/group/curve25519"
	"github.com/dedis/kyber/group/nist"
)

func init() {
	curve25519 := curve25519.NewAES128SHA256Ed25519(false)
	suites[curve25519.String()] = curve25519

	p256 := nist.NewAES128SHA256P256()
	suites[p256.String()] = p256

	qr512 := nist.NewAES128SHA256QR512()
	suites[qr512.String()] = qr512
}
