package zerotier

import (
	"os"
	"testing"
)

func TestRequest(t *testing.T) {

	zt, err := NewClient(os.Getenv("ZT_TOKEN"))
	if err != nil {
		t.Fatal("error", err)
	}
	n, err := zt.getMembers(os.Getenv("NETWORK_ID"))
	if err != nil {
		t.Fatal("error", err)
	}
	t.Log("members:", n)
	t.Fatal("args ...interface{}")
}
