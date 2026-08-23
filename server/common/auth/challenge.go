package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Challenge struct {
	Domain    string
	URI       string
	ChainID   int64
	Address   common.Address
	Nonce     string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Message   string
}

func NewChallenge(
	domain string,
	uri string,
	chainID int64,
	address common.Address,
	now time.Time,
	ttl time.Duration,
) (Challenge, error) {
	parsed, err := url.Parse(uri)
	if err != nil || !validChallengeScheme(parsed) || !strings.EqualFold(parsed.Hostname(), domain) {
		return Challenge{}, errors.New("domain and secure URI do not match")
	}
	if chainID != 84532 {
		return Challenge{}, errors.New("chain must be Base Sepolia")
	}
	if address == (common.Address{}) || ttl <= 0 {
		return Challenge{}, errors.New("address and TTL are required")
	}
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return Challenge{}, err
	}
	challenge := Challenge{
		Domain: domain, URI: uri, ChainID: chainID, Address: address,
		Nonce: hex.EncodeToString(nonceBytes), IssuedAt: now.UTC(), ExpiresAt: now.Add(ttl).UTC(),
	}
	challenge.Message = challenge.render()
	return challenge, nil
}

func validChallengeScheme(parsed *url.URL) bool {
	if parsed.Scheme == "https" {
		return true
	}
	if parsed.Scheme != "http" {
		return false
	}
	host := parsed.Hostname()
	return strings.EqualFold(host, "localhost") || net.ParseIP(host).IsLoopback()
}

func PersonalSignHash(message string) []byte {
	return accounts.TextHash([]byte(message))
}

func Verify(challenge Challenge, signatureHex string, address common.Address, now time.Time) error {
	if challenge.ChainID != 84532 || challenge.Message != challenge.render() {
		return errors.New("challenge fields do not match message")
	}
	if address != challenge.Address {
		return errors.New("challenge address mismatch")
	}
	if !now.Before(challenge.ExpiresAt) {
		return errors.New("challenge expired")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(signatureHex, "0x"))
	if err != nil || len(signature) != crypto.SignatureLength {
		return errors.New("invalid signature encoding")
	}
	if signature[64] >= 27 {
		signature[64] -= 27
	}
	publicKey, err := crypto.SigToPub(PersonalSignHash(challenge.Message), signature)
	if err != nil || crypto.PubkeyToAddress(*publicKey) != address {
		return errors.New("signature address mismatch")
	}
	return nil
}

func (c Challenge) render() string {
	return fmt.Sprintf(
		"%s wants you to sign in to o1 Launchpad.\n\nURI: %s\nChain ID: %d\nAddress: %s\nNonce: %s\nIssued At: %s\nExpiration Time: %s",
		c.Domain, c.URI, c.ChainID, c.Address.Hex(), c.Nonce,
		c.IssuedAt.UTC().Format(time.RFC3339), c.ExpiresAt.UTC().Format(time.RFC3339),
	)
}
