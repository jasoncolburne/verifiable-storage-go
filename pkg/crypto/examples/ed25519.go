package examples

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"sync"

	interfaces "github.com/jasoncolburne/verifiable-storage-go/pkg/crypto"
)

type Ed25519VerificationKey struct {
	publicKey string
}

func NewEd25519VerificationKey(publicKey string) *Ed25519VerificationKey {
	return &Ed25519VerificationKey{
		publicKey: publicKey,
	}
}

func (e Ed25519VerificationKey) Public() (string, error) {
	return e.publicKey, nil
}

func (e Ed25519VerificationKey) Verifier() interfaces.Verifier {
	return NewEd25519Verifier()
}

type Ed25519Verifier struct{}

func NewEd25519Verifier() *Ed25519Verifier {
	return &Ed25519Verifier{}
}

func (e Ed25519Verifier) Verify(signature, publicKey string, message []byte) error {
	keyBytes, err := base64.URLEncoding.DecodeString(publicKey)
	if err != nil {
		return err
	}

	edKey := ed25519.PublicKey(keyBytes[1:])

	sigBytes, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	if !ed25519.Verify(edKey, message, sigBytes[2:]) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

type Ed25519 struct {
	signingKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewEd25519(seed []byte) (*Ed25519, error) {
	var publicKey ed25519.PublicKey
	var signingKey ed25519.PrivateKey
	var err error

	if seed == nil {
		publicKey, signingKey, err = ed25519.GenerateKey(nil)
		if err != nil {
			return nil, err
		}
	} else {
		var ok bool

		signingKey = ed25519.NewKeyFromSeed(seed)
		publicKey, ok = signingKey.Public().(ed25519.PublicKey)
		if !ok {
			return nil, fmt.Errorf("incorrect key type")
		}
	}

	return &Ed25519{
		signingKey: signingKey,
		publicKey:  publicKey,
	}, nil
}

func (e Ed25519) Identity() (string, error) {
	return e.Public()
}

func (e Ed25519) Sign(message []byte) (string, error) {
	sig := ed25519.Sign(e.signingKey, message)
	buffer := [66]byte{}

	copy(buffer[2:], sig)

	b64 := base64.URLEncoding.EncodeToString(buffer[:])
	qb64 := []rune(b64)
	qb64[0] = '0'
	qb64[1] = 'B'

	return string(qb64), nil
}

func (e Ed25519) Verifier() interfaces.Verifier {
	return NewEd25519Verifier()
}

func (e Ed25519) Public() (string, error) {
	buffer := [33]byte{}

	copy(buffer[1:], e.publicKey)
	b64 := base64.URLEncoding.EncodeToString(buffer[:])
	qb64 := []rune(b64)
	qb64[0] = 'B'

	return string(qb64), nil
}

type VerificationKeyStore struct {
	mu   sync.RWMutex
	keys map[string]interfaces.VerificationKey
}

func NewVerificationKeyStore() *VerificationKeyStore {
	return &VerificationKeyStore{
		keys: make(map[string]interfaces.VerificationKey),
	}
}

func (s *VerificationKeyStore) Add(identity string, key interfaces.VerificationKey) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.keys[identity] = key
}

func (s *VerificationKeyStore) Get(identity string) (interfaces.VerificationKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key, exists := s.keys[identity]
	if !exists {
		return nil, fmt.Errorf("key not found for identity: %s", identity)
	}
	return key, nil
}
