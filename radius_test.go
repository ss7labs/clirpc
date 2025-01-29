package clirpc

import (
	"testing"
)

func TestNewRadiusConfig(t *testing.T) {
	host := "192.168.3.1"
	port := 3799
	secret := "secret"

	config := NewRadiusConfig(host, port, secret)

	if config.Host != host {
		t.Errorf("expected host %s, got %s", host, config.Host)
	}
	if config.Port != port {
		t.Errorf("expected port %d, got %d", port, config.Port)
	}
	if config.Secret != secret {
		t.Errorf("expected secret %s, got %s", secret, config.Secret)
	}
	if config.Username != "0:0:130" {
		t.Errorf("expected username 0:0:130, got %s", config.Username)
	}
}
