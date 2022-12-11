package keys

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// Key represents an SSH public key in the authorized_keys format.
type Key struct {
	raw    string
	name   string
	pubkey ssh.PublicKey
}

// NewKey returns a new key based on the raw key string and the key name
// (comment). If the key can not be parsed as a valid authorized key, an error
// is returned.
func NewKey(raw, name string) (*Key, error) {
	pubkey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(raw))
	if err != nil {
		return nil, fmt.Errorf("could not parse key as a valid authorized key: %w", err)
	}
	return &Key{
		raw:    raw,
		name:   name,
		pubkey: pubkey,
	}, nil
}

// Name returns the name of the key as determined by the key comment.
func (k Key) Name() string {
	return k.name
}

// Raw returns the key in a useable authorized_key format.
func (k Key) Raw() string {
	return k.raw
}

// Type returns the type of the SSH key (e.g. ssh-rsa).
func (k Key) Type() string {
	return k.pubkey.Type()
}

// Fingerprint returns the SHA256 fingerprint of the key.
func (k Key) Fingerprint() string {
	return ssh.FingerprintSHA256(k.pubkey)
}

// CreateAuthorizedKeys returns a single string built of a set of keys all
// separated with a newline (\n). This can be used in SSH authorized_key files.
func CreateAuthorizedKeys(keys []*Key) string {
	authorizedKeys := ""
	for _, k := range keys {
		authorizedKeys += k.Raw() + "\n"
	}
	return authorizedKeys
}
