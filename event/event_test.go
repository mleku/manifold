package event

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/mleku/manifold/log"
	"lukechampine.com/frand"
)

func TestSplit(t *testing.T) {
	b := frand.Bytes(64)
	bb := make([]byte, 86)
	base64.RawStdEncoding.Encode(bb, b)
	log.I.F("%s", bb)
	data := []byte(
		"PUBKEY:c29tZV9yYW5kb21fZGF0YV9lbmNvZGVkX3NjZW5hcml\n" +
			"TIMESTAMP:1672531200\n" +
			"CONTENT:example content\nwith linebreak\n" +
			"TAG:key1:value1\n" +
			"TAG:key2:value2\n" +
			"TAG:key3:value3\n" +
			"SIGNATURE:DCw0ytbRPs3Q1LrQ7HhFYDE2Au8hWj4TZFVu/MJDdjJjirmznTFMSsyL1UM39KW18zgmHeQ5qjAqCo70fa23kQ",
	)
	fmt.Printf("%s\n", data)
	e, err := Split(data)
	if err != nil {
		t.Fatalf("Error: %v", err)
	} else {
		t.Logf("Parsed Event:\n%s\n", e)
	}
}
