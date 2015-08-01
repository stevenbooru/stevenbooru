package config

import "github.com/scalingdata/gcfg"

// ParseConfig loads and parses a configuration file for Stevenbooru.
func ParseConfig(fname string) (Config, error) {
	cfg := Config{}

	return cfg, gcfg.ReadFileInto(&cfg, fname)
}
