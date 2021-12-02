package secp

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/icodeface/chain-kit/filecoin/sigs"
	"github.com/minio/blake2b-simd"
)

// PrivateKeyBytes is the size of a serialized private key.
const PrivateKeyBytes = 32

type secpSigner struct{}

func (secpSigner) GenPrivate() ([]byte, error) {
	key, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	privkey := make([]byte, PrivateKeyBytes)
	blob := key.D.Bytes()

	// the length is guaranteed to be fixed, given the serialization rules for secp2561k curve points.
	copy(privkey[PrivateKeyBytes-len(blob):], blob)

	return privkey, nil
}

func (secpSigner) ToPublic(pk []byte) ([]byte, error) {
	x, y := btcec.S256().ScalarBaseMult(pk)
	return elliptic.Marshal(btcec.S256(), x, y), nil
}

func (secpSigner) Sign(pk []byte, msg []byte) ([]byte, error) {
	b2sum := blake2b.Sum256(msg)
	p, _ := btcec.PrivKeyFromBytes(btcec.S256(), pk)
	sig, err := btcec.SignCompact(btcec.S256(), p, b2sum[:], false)
	if err != nil {
		return nil, err
	}

	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v

	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (secpSigner) Verify(signature []byte, a address.Address, msg []byte) error {
	b2sum := blake2b.Sum256(msg)

	var sig = make([]byte, 65)
	copy(sig, signature)

	v := sig[64] + 27
	copy(sig[1:], sig[:64])
	sig[0] = v

	pk, _, err := btcec.RecoverCompact(btcec.S256(), sig, b2sum[:])
	if err != nil {
		return err
	}
	pubk := pk.SerializeUncompressed()

	if err != nil {
		return err
	}

	maybeaddr, err := address.NewSecp256k1Address(pubk)
	if err != nil {
		return err
	}

	if a != maybeaddr {
		return fmt.Errorf("signature did not match")
	}

	return nil
}

func init() {
	sigs.RegisterSignature(crypto.SigTypeSecp256k1, secpSigner{})
}
