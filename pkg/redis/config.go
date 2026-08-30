package redis

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// TLSVersion represents a TLS version as a string for configuration purposes.
// Valid values: "1.2", "1.3". Defaults to "1.2" when empty or unset.
type TLSVersion string

const (
	TLSVersion12 TLSVersion = "1.2"
	TLSVersion13 TLSVersion = "1.3"
)

// TLSVersionID returns the crypto/tls constant for the configured version.
// Defaults to tls.VersionTLS12 for unrecognized or empty values.
func (v TLSVersion) TLSVersionID() uint16 {
	switch v {
	case TLSVersion13:
		return tls.VersionTLS13
	default:
		return tls.VersionTLS12
	}
}

// Config stores the user provided configuration parameters
type Config struct {
	Address  string
	Database int
	Username string
	Password string
	Sentinel struct {
		Enabled    bool
		MasterName string
	}
	TLS struct {
		Enabled            bool
		InsecureSkipVerify bool
		MinVersion         TLSVersion
	}
}

// Validate checks if the configuration is valid
func (c Config) Validate() error {
	var errs []error

	if c.Address == "" {
		errs = append(errs, errors.New("address is required"))
	}

	if c.Sentinel.Enabled {
		if c.Sentinel.MasterName == "" {
			errs = append(errs, errors.New("sentinel: master name is required"))
		}
	}

	return errors.Join(errs...)
}

// New client returns a redis client returns a redis client from the provided configuration
func (c Config) NewClient() (*redis.Client, error) {
	return NewClient(Options{Config: c})
}

// Configure sets the default values for the Redis configuration
func Configure(v *viper.Viper, prefix string) {
	// Redis driver
	v.SetDefault(fmt.Sprintf("%s.address", prefix), "127.0.0.1:6379")
	v.SetDefault(fmt.Sprintf("%s.database", prefix), 0)
	v.SetDefault(fmt.Sprintf("%s.username", prefix), "")
	v.SetDefault(fmt.Sprintf("%s.password", prefix), "")
	v.SetDefault(fmt.Sprintf("%s.expiration", prefix), "24h")
	v.SetDefault(fmt.Sprintf("%s.sentinel.enabled", prefix), false)
	v.SetDefault(fmt.Sprintf("%s.sentinel.masterName", prefix), "")
	v.SetDefault(fmt.Sprintf("%s.tls.enabled", prefix), false)
	v.SetDefault(fmt.Sprintf("%s.tls.insecureSkipVerify", prefix), false)
	v.SetDefault(fmt.Sprintf("%s.tls.minVersion", prefix), string(TLSVersion12))
}
