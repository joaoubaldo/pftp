package pftp

import (
	"crypto/tls"

	"github.com/BurntSushi/toml"
)

type portRange struct {
	Start int
	End   int
}

type config struct {
	ListenAddr      string      `toml:"listen_addr"`
	RemoteAddr      string      `toml:"remote_addr"`
	IdleTimeout     int         `toml:"idle_timeout"`
	ProxyTimeout    int         `toml:"proxy_timeout"`
	TransferTimeout int         `toml:"transfer_timeout"`
	MaxConnections  int32       `toml:"max_connections"`
	TLS             *tlsPair    `toml:"tls"`
	TLSConfig       *tls.Config `toml:"-"`
}

type tlsPair struct {
	Cert string `toml:"cert"`
	Key  string `toml:"key"`
}

func loadConfig(path string) (*config, error) {
	var c config
	defaultConfig(&c)

	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		return nil, err
	}

	if c.TLS != nil {
		if cert, err := tls.LoadX509KeyPair(c.TLS.Cert, c.TLS.Key); err == nil {
			c.TLSConfig = &tls.Config{
				NextProtos:   []string{"ftp"},
				Certificates: []tls.Certificate{cert},
			}
		} else {
			return nil, err
		}
	}

	return &c, nil
}

func defaultConfig(config *config) {
	config.ListenAddr = "0.0.0.0:2121"
	config.IdleTimeout = 900
	config.ProxyTimeout = 900
	config.TransferTimeout = 900
}
