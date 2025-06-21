//go:build cgo

package p256k_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/mleku/manifold/chk"
	"github.com/mleku/manifold/log"
	"github.com/mleku/manifold/p256k"
	"github.com/mleku/manifold/sha256"
	realy "github.com/mleku/manifold/signer"
)

func TestSigner_Generate(t *testing.T) {
	for _ = range 10000 {
		var err error
		signer := &p256k.Signer{}
		var skb []byte
		if err = signer.Generate(); chk.E(err) {
			t.Fatal(err)
		}
		skb = signer.Sec()
		if err = signer.InitSec(skb); chk.E(err) {
			t.Fatal(err)
		}
	}
}

func TestSignerVerify(t *testing.T) {
	// Initialize a new signer
	signer := &p256k.Signer{}
	err := signer.Generate()
	if chk.E(err) {
		t.Fatalf("Failed to generate signer key pair: %v", err)
	}

	// Sample message to sign
	message := sha256.Sum256Bytes([]byte("Hello, world!"))
	// Sign the message
	signature, err := signer.Sign(message)
	if chk.E(err) {
		t.Fatalf("Failed to sign message: %v", err)
	}

	// Verify the signature
	valid, err := signer.Verify(message, signature)
	if chk.E(err) {
		t.Fatalf("Error verifying signature: %v", err)
	}

	// Check if the signature is valid
	if !valid {
		t.Error("Valid signature was rejected")
	}

	// Modify the message and verify again
	tamperedMessage := sha256.Sum256Bytes([]byte("Hello, tampered world!"))
	valid, err = signer.Verify(tamperedMessage, signature)
	if err == nil {
		t.Fatalf("Error verifying tampered message: %v", err)
	}

	// Expect the verification to fail
	if valid {
		t.Error("Invalid signature was accepted")
	}
}

func TestECDH(t *testing.T) {
	n := time.Now()
	var err error
	var s1, s2 realy.I
	var counter int
	const total = 100
	for _ = range total {
		s1, s2 = &p256k.Signer{}, &p256k.Signer{}
		if err = s1.GenerateForECDH(); chk.E(err) {
			t.Fatal(err)
		}
		for _ = range total {
			if err = s2.GenerateForECDH(); chk.E(err) {
				t.Fatal(err)
			}
			var secret1, secret2 []byte
			if secret1, err = s1.ECDH(s2.Pub()); chk.E(err) {
				t.Fatal(err)
			}
			if secret2, err = s2.ECDH(s1.Pub()); chk.E(err) {
				t.Fatal(err)
			}
			if !bytes.Equal(secret1, secret2) {
				counter++
				t.Errorf("ECDH generation failed to work in both directions, %x %x", secret1,
					secret2)
			}
		}
	}
	a := time.Now()
	duration := a.Sub(n)
	log.I.Ln("errors", counter, "total", total*total, "time", duration, "time/op",
		duration/total/total, "ops/sec", float64(time.Second)/float64(duration/total/total))
}
