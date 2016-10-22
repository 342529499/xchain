package util

import (
	"testing"
)

func TestSetLocalIP(t *testing.T) {
	ip := GetLocalIP()

	err := SetLocalIP(ip)
	if err != nil {
		t.Fatalf(err)
	}

	ip2 := GetLocalIP()
	if ip != ip2 {
		t.Fatalf("set local ip fail")
	}
}
