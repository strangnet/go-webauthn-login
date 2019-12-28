package domain

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
)

type User struct {
	id uuid.UUID
	wid uint64
	username string
	displayName string
	credentials []webauthn.Credential
}

type UserRepository interface {
	Create(user *User) error
	FindByUsername(username string) (*User, error)
	ListAllUsers() []*User
}

func NewUser(username string, displayName string) *User {
	return &User{
		id: uuid.New(),
		wid: randomUint64(),
		username: username,
		displayName: displayName,
	}
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

func (u *User) Username() string {
	return u.username
}

func (u *User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.wid))
	return buf

}

func (u *User) WebAuthnName() string {
	return u.username
}

func (u *User) WebAuthnDisplayName() string {
	return u.displayName
}

func (u *User) WebAuthnIcon() string {
	return ""
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func (u *User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
}

func (u *User) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
