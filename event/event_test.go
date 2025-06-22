package event

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/mleku/manifold/chk"
	"github.com/mleku/manifold/p256k"
)

func TestUnmarshal_Marshal(t *testing.T) {
	data := []byte(`PUBKEY:Opp5HiOfR7wC8OLN-HudU49y1k-ZzsUt4fUJSHt7LXo
TIMESTAMP:1672531200
CONTENT:b64:XFbfJ2Pvwqjim3MSfE0pD0O1g_TssKOr9Y7tSgOjsdsA
TAG:key1:value1
TAG:key2:value2
TAG:hashtag:winning
TAG:mention:b64:XFbfJ2Pvwqjim3MSfE0pD0O1g_TssKOr9Y7tSgOjsdsA
SIGNATURE:9DOOTXtcIZqcO7LaRaNAD8s9BjMyf46qp75NNJb_T-5piA57L4EjGYIx3Fok8L3pSIH7hB1XNeJwAbaLiCWgjA`,
	)
	fmt.Printf("original raw event:\n%s\n", data)
	e := new(E)
	var err error
	if err = e.Unmarshal(data); chk.E(err) {
		t.Fatalf("Error: %v", err)
	} else {
		fmt.Printf("\nUnmarshalled Event:\n%s\n", spew.Sdump(e))
	}
	var b []byte
	if b, err = e.Marshal(); chk.E(err) {
		t.Fatalf("Error: %v", err)
	}
	fmt.Printf("\nMarshalled Event:\n%s\n", b)
	var valid bool
	if valid, err = e.Verify(); chk.E(err) {
		t.Fatalf("failed to verify event: %v", err)
	}
	if !valid {
		t.Fatalf("event signature is invalid")
	}
	t.Log("event signature is valid")
	var c []byte
	e.Signature = nil
	if c, err = e.Marshal(); chk.E(err) {
		t.Fatalf("Error: %v", err)
	}
	fmt.Printf("\nMarshalled canonical Event:\n%s\n", c)
	var id []byte
	if id, err = e.Id(); chk.E(err) {
		t.Fatalf("failed to get event Id %v", err)
	}
	id64 := make([]byte, 43)
	base64.RawStdEncoding.Encode(id64, id)
	fmt.Printf("\nEvent Id: %s\n", id64)
	sign := new(p256k.Signer)
	if err = sign.Generate(); chk.E(err) {
		t.Fatalf("failed to generate key pair: %v", err)
	}
	e.Pubkey = sign.Pub()
	if err = e.Sign(sign); chk.E(err) {
		t.Fatalf("failed to sign event: %v", err)
	}
	if b, err = e.Marshal(); chk.E(err) {
		t.Fatalf("failed to marshal signed event: %v", err)
	}
	fmt.Printf("\nSigned Event:\n%s\n", b)
}
