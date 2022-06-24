package frontend_interview

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/TBD54566975/ssi-sdk/credential"
	"github.com/TBD54566975/ssi-sdk/credential/signing"
	"github.com/TBD54566975/ssi-sdk/credential/status"
	"github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/cryptosuite"
	"github.com/TBD54566975/ssi-sdk/did"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

var (
	knownSeed0 = []byte("12345678901234567890123456789012")
	knownSeed1 = []byte("01234567890121234567890123456789")
	knownDID0  = "did:key:z6MkhesMp8iSdumBExtuozsz3PYfapPpQUCarQA5uLcRee4d"
	knownDID1  = "did:key:z6Mko2Svq9tXt1WVmgAvHbkQE2RtCTUzfbgmjQMB7wcQa5SA"
)

func GetDID(id int) did.DIDDocument {
	var seed []byte
	if id == 0 {
		seed = knownSeed0
	} else if id == 1 {
		seed = knownSeed1
	} else {
		panic("id can only be 0 or 1")
	}
	privateKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privateKey.Public()
	didKey, _ := did.CreateDIDKey(crypto.Ed25519, pubKey.(ed25519.PublicKey))
	expanded, _ := didKey.Expand()
	return *expanded
}

func GetVCJSONSchema() string {
	return `{
  "type": "https://w3c-ccg.github.io/vc-json-schemas/schema/2.0/schema.json",
  "version": "1.0",
  "id": "did:key:z6MkhesMp8iSdumBExtuozsz3PYfapPpQUCarQA5uLcRee4d;id=06e126d1-fa44-4882-a243-1e326fbe21db;version=1.0",
  "name": "Email",
  "author": "did:key:z6MkhesMp8iSdumBExtuozsz3PYfapPpQUCarQA5uLcRee4d",
  "authored": "2025-05-05T00:00:00+00:00",
  "schema": {
    "$id": "email-schema-1.0",
    "$schema": "https://json-schema.org/draft/2019-09/schema",
    "description": "Email",
    "type": "object",
    "properties": {
      "emailAddress": {
        "type": "string",
        "format": "email"
      }
    },
    "required": [
      "emailAddress"
    ],
    "additionalProperties": false
  }
}`
}

func GetCredentialVC(id int) credential.VerifiableCredential {
	credJWT := getCredentialJWT(id)
	cred, _ := signing.ParseVerifiableCredentialFromJWT(string(credJWT))
	return *cred
}

func GetCredentialToken(id int) jwt.Token {
	parsed, _ := jwt.Parse(getCredentialJWT(id))
	return parsed
}

func getCredentialJWT(id int) []byte {
	var issuer, subject = "", ""
	if id == 0 {
		issuer = knownDID0
		subject = knownDID1
	} else if id == 1 {
		issuer = knownDID1
		subject = knownDID0
	} else {
		panic("id can only be 0 or 1")
	}
	testCredential := credential.VerifiableCredential{
		Context:      []interface{}{"https://www.w3.org/2018/credentials/v1", "https://w3id.org/security/suites/jws-2020/v1"},
		ID:           fmt.Sprintf("credential-%d", id),
		Type:         []string{"VerifiableCredential"},
		Issuer:       issuer,
		IssuanceDate: "2025-05-05T05:05:05Z",
		CredentialStatus: status.StatusList2021Entry{
			ID:                   fmt.Sprintf("revocation-id-%d", id),
			Type:                 status.StatusList2021EntryType,
			StatusPurpose:        status.StatusRevocation,
			StatusListIndex:      strconv.Itoa(id),
			StatusListCredential: "status-list-credential",
		},
		CredentialSubject: map[string]interface{}{
			"id":           subject,
			"emailAddress": "interview@tbd.email",
		},
	}
	signer, _ := cryptosuite.NewJSONWebKeySigner(issuer, getKnownJWK(id).PrivateKeyJWK, cryptosuite.AssertionMethod)
	signed, _ := signing.SignVerifiableCredentialJWT(*signer, testCredential)
	return signed
}

func GetRevocation(id int) credential.VerifiableCredential {
	var issuer string
	var creds []credential.VerifiableCredential
	if id == 0 {
		issuer = knownDID0
	} else if id == 1 {
		issuer = knownDID0
		creds = []credential.VerifiableCredential{GetCredentialVC(0)}
	} else if id == 2 {
		issuer = knownDID1
	} else if id == 3 {
		issuer = knownDID1
		creds = []credential.VerifiableCredential{GetCredentialVC(1)}
	} else {
		panic("id must be 0, 1, 2, or 3")
	}
	statusCred, _ := status.GenerateStatusList2021Credential("status-list-credential", issuer, status.StatusRevocation, creds)
	return *statusCred
}

func getKnownJWK(id int) cryptosuite.JSONWebKey2020 {
	var privateKey ed25519.PrivateKey
	if id == 0 {
		privateKey = ed25519.NewKeyFromSeed(knownSeed0)
	} else if id == 1 {
		privateKey = ed25519.NewKeyFromSeed(knownSeed1)
	} else {
		panic("id can only be 0 or 1")
	}
	ed25519JWK := jwk.NewOKPPrivateKey()
	_ = ed25519JWK.FromRaw(privateKey)

	kty := ed25519JWK.KeyType().String()
	crv := ed25519JWK.Crv().String()
	x := base64.RawURLEncoding.EncodeToString(ed25519JWK.X())
	return cryptosuite.JSONWebKey2020{
		Type: cryptosuite.JsonWebKey2020,
		PrivateKeyJWK: cryptosuite.PrivateKeyJWK{
			KTY: kty,
			CRV: crv,
			X:   x,
			D:   base64.RawURLEncoding.EncodeToString(ed25519JWK.D()),
		},
		PublicKeyJWK: cryptosuite.PublicKeyJWK{
			KTY: kty,
			CRV: crv,
			X:   x,
		},
	}
}
