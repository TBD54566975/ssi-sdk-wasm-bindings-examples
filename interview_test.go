//go:build jwx_es256k

package main

import (
	"encoding/json"
	"testing"
)

func TestGenerateDIDs(t *testing.T) {
	println("DID: 0")
	prettyPrint(GetDID(0))
	println("DID: 1")
	prettyPrint(GetDID(1))
}

func TestGetSchema(t *testing.T) {
	prettyPrint(GetVCJSONSchema())
}

func TestGenerateCredentials(t *testing.T) {
	println("Credential: 0")
	prettyPrint(GetCredentialToken(0))
	println("Credential: 1")
	prettyPrint(GetCredentialToken(1))
}

func TestGenerateRevocations(t *testing.T) {
	println("Revocation: 0")
	prettyPrint(GetRevocation(0))
	println("Revocation: 1")
	prettyPrint(GetRevocation(1))
	println("Revocation: 2")
	prettyPrint(GetRevocation(2))
	println("Revocation: 3")
	prettyPrint(GetRevocation(3))
}

func prettyPrint(d interface{}) {
	prettyBytes, _ := json.MarshalIndent(d, "", " ")
	println(string(prettyBytes))
}
