package config

import "testing"

func TestParseExampleConfig(t *testing.T) {
	cfg, err := ParseConfig("../../../cfg/stevenbooru.cfg")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Site.Name != "Stevenbooru" {
		t.Fatalf("Site name should be Stevenbooru, not %s", cfg.Site.Name)
	}

	if !cfg.Site.Testing {
		t.Fatalf("Site should be in testing mode, and it is not.")
	}
}
