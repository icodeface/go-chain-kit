package ethereum

import "testing"

func TestWallet(t *testing.T) {
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, err := NewWallet(mnemonic)
	if err != nil {
		t.Error(err)
	}

	account, err := wallet.DeriveAccount(DerivePath(0, 0))
	if err != nil {
		t.Error(err)
	}
	if account.Address.Hex() != "0xC49926C4124cEe1cbA0Ea94Ea31a6c12318df947" {
		t.Error("wrong address")
	}

	account, err = wallet.DeriveAccount(DerivePath(0, 19))
	if err != nil {
		t.Error(err)
	}
	if account.Address.Hex() != "0x816A5f3ED3FB0DCb5C19A32C80cc9643fDB078EB" {
		t.Error("wrong address")
	}

}
