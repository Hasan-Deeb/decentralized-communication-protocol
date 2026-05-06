package peer

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// NodeIdentity represents the cryptographic identity of a P2P node
type NodeIdentity struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	PeerID     string
}

// GenerateIdentity creates a new Ed25519 key pair and derives the PeerID
func GenerateIdentity() (*NodeIdentity, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 keys: %w", err)
	}

	// PeerID is the hex string of the SHA-256 hash of the Public Key
	hash := sha256.Sum256(pub)
	peerID := hex.EncodeToString(hash[:])

	return &NodeIdentity{
		PublicKey:  pub,
		PrivateKey: priv,
		PeerID:     peerID,
	}, nil
}
