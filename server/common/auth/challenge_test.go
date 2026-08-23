package auth

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestChallengeRejectsWrongDomainChainAddressAndExpiry(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	now := time.Unix(1_700_000_000, 0).UTC()
	challenge, err := NewChallenge("demo.example", "https://demo.example", 84532, address, now, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := crypto.Sign(PersonalSignHash(challenge.Message), key)
	if err != nil {
		t.Fatal(err)
	}
	signature := "0x" + hex.EncodeToString(sig)

	if err := Verify(challenge, signature, address, now.Add(time.Minute)); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	wrong := challenge
	wrong.Domain = "evil.example"
	if err := Verify(wrong, signature, address, now.Add(time.Minute)); err == nil {
		t.Fatal("expected domain mismatch")
	}
	wrong = challenge
	wrong.ChainID = 1
	if err := Verify(wrong, signature, address, now.Add(time.Minute)); err == nil {
		t.Fatal("expected chain mismatch")
	}
	if err := Verify(challenge, signature, crypto.CreateAddress(address, 1), now.Add(time.Minute)); err == nil {
		t.Fatal("expected address mismatch")
	}
	if err := Verify(challenge, signature, address, challenge.ExpiresAt); err == nil {
		t.Fatal("expected expiration error")
	}
}

func TestChallengeMessageBindsRequiredFields(t *testing.T) {
	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey)
	challenge, err := NewChallenge(
		"demo.example", "https://demo.example", 84532, address, time.Unix(1_700_000_000, 0).UTC(), 5*time.Minute,
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{"demo.example", "https://demo.example", "84532", address.Hex(), challenge.Nonce} {
		if !strings.Contains(challenge.Message, value) {
			t.Fatalf("message does not contain %q", value)
		}
	}
}

func TestChallengeAllowsOnlyLoopbackHTTP(t *testing.T) {
	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey)
	now := time.Unix(1_700_000_000, 0).UTC()
	if _, err := NewChallenge("localhost", "http://localhost:3000", 84532, address, now, time.Minute); err != nil {
		t.Fatalf("loopback HTTP should be allowed in local development: %v", err)
	}
	if _, err := NewChallenge("demo.example", "http://demo.example", 84532, address, now, time.Minute); err == nil {
		t.Fatal("non-loopback HTTP must be rejected")
	}
}
