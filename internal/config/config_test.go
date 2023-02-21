package config

import "testing"

func Test_config(t *testing.T) {
	wantedPort := "5000"

	cfg := GetConfig()
	if cfg.Listen.Port != wantedPort {
		t.Errorf("port = %s, want = %s", cfg.Listen.Port, wantedPort)
	}
}
