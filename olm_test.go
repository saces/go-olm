package olm

import "testing"

import (
	"encoding/json"
)

func TestOlm(t *testing.T) {
	major, minor, patch := Version()
	t.Log("Version():", major, minor, patch)
}

func TestAccount(t *testing.T) {
	// Create a new account
	a1 := NewAccount()
	pickled1 := a1.Pickle([]byte("HELLO"))
	t.Log("Pickle():", pickled1)

	// Create a second account
	a11 := NewAccount()
	pickled11 := a11.Pickle([]byte("HELLO"))
	if pickled1 == pickled11 {
		t.Error("Two new accounts pickle to the same string")
	}

	// Load account from pickled account
	a2, err := AccountFromPickled(pickled1, []byte("HELLO"))
	if err != nil {
		t.Error(err)
	}
	// Pickle again to verify that we get the same as before
	pickled2 := a2.Pickle([]byte("HELLO"))
	t.Log("Pickle():", pickled2)
	a3, err := AccountFromPickled(pickled2, []byte("HELLO"))
	if err != nil {
		t.Error(err)
	}
	if pickled1 != pickled2 {
		t.Error("pickle(unpickle(pickle)) != pickle")
	}

	identityKeys := a1.IdentityKeys()
	t.Log("IdentityKeys():", identityKeys)

	// Sign a message
	signature := a1.Sign("HELLO WORLD")
	t.Log("a1.Sign():", signature)
	maxNumberOfOneTimeKeys := a1.MaxNumberOfOneTimeKeys()
	a1.GenOneTimeKeys(maxNumberOfOneTimeKeys)
	oneTimeKeys := a1.OneTimeKeys()
	t.Log("a1.OneTimeKeys():", oneTimeKeys)

	// Sign the same message again (signature should be different)
	message := "HELLO WORLD"
	s := a1.Sign(message)
	t.Logf("Sign(\"%s\"): %s", message, s)

	t.Log("a1.Clear():", a1.Clear())
	t.Log("a2.Clear():", a2.Clear())
	t.Log("a3.Clear():", a3.Clear())
}

type OneTimeKeys struct {
	Curve25519 map[string]string `json:"curve25519"`
}

type IdentityKeys struct {
	Curve25519 string `json:"curve25519"`
	Ed25519    string `json:"ed25519"`
}

func TestSession(t *testing.T) {
	// Generate new accounts
	a1 := NewAccount()
	a2 := NewAccount()

	// Generate one time keys for account 2, and get one of them
	a2.GenOneTimeKeys(a2.MaxNumberOfOneTimeKeys())
	a2OneTimeKeysJSON := a2.OneTimeKeys()
	var a2OneTimeKeys OneTimeKeys
	json.Unmarshal([]byte(a2OneTimeKeysJSON), &a2OneTimeKeys)
	//t.Log("Marshaled:", a2OneTimeKeysJSON)
	//t.Logf("Unmarshaled: %+v", a2OneTimeKeys)
	// Pick one One Time Key
	var a2OneTimeKey string
	for _, v := range a2OneTimeKeys.Curve25519 {
		a2OneTimeKey = v
		break
	}

	// Get identity key of account 2
	a2IdentityKeysJSON := a2.IdentityKeys()
	var a2IdentityKeys IdentityKeys
	json.Unmarshal([]byte(a2IdentityKeysJSON), &a2IdentityKeys)
	t.Log("a2IdentityKeys:", a2IdentityKeys)
	t.Log("a2OneTimeKey:", a2OneTimeKey)

	// From account 1, generate an outbout session towards account 2
	s1, err := a1.NewOutboundSession(a2IdentityKeys.Curve25519, a2OneTimeKey)
	if err != nil {
		t.Error(err)
	}
	pickled1 := s1.Pickle([]byte("HELLO"))
	t.Log("Pickle():", pickled1)
}

func TestUtility(t *testing.T) {
	u := NewUtility()

	// Hash a string
	h := u.Sha256("HELLO")
	t.Log("Sha256():", h)
	if h != "NzPNl3/46xi5hzV+Is7Zn0YJfzHssjnoeK5jdg6D5NU" {
		t.Error("Sha256 doesn't match")
	}

	// Verify a signed message
	message := "HELLO WORLD"
	key := "TdbnI8JjtbJW1h9dISHcZ7LTpMYIjKFiEBfKp8hxCeI"
	signature := "h6SV3IO8S0sOMyvUvgbQcLaPkP0utyXDFHMrAVoLZl87JG3z8thYo9L1jHusXtP+fXM9NB7E2p06udpmtIPHAQ"
	ok, err := u.Ed25519Verify(message, key, signature)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Log("Signature verification shouldn't have failed")
	}

	// Verify an incorrect signed message
	message = "GOOD BYE"
	ok, err = u.Ed25519Verify(message, key, signature)
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Log("Signature verification should have failed")
	}
}
