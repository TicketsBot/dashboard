package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type (
	Config struct {
		Admins          []uint64
		ForceWhitelabel []uint64
		Server          Server
		Oauth           Oauth
		Database        Database
		Bot             Bot
		Redis           Redis
		Cache           Cache
	}

	Server struct {
		Host      string
		BaseUrl   string
		MainSite  string
		Ratelimit Ratelimit
		Session   Session
		Secret    string
	}

	Ratelimit struct {
		Window int
		Max    int
	}

	Session struct {
		Threads int
		Secret  string
	}

	Oauth struct {
		Id          int64
		Secret      string
		RedirectUri string
	}

	Database struct {
		Uri string
	}

	Bot struct {
		Token                 string
		PremiumLookupProxyUrl string `toml:"premium-lookup-proxy-url"`
		PremiumLookupProxyKey string `toml:"premium-lookup-proxy-key"`
		ObjectStore           string
		AesKey                string `toml:"aes-key"`
	}

	Redis struct {
		Host     string
		Port     int
		Password string
		Threads  int
	}

	Cache struct {
		Uri string
	}
)

var (
	Conf Config
)

func LoadConfig() {
	raw, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	_, err = toml.Decode(string(raw), &Conf)
	if err != nil {
		panic(err)
	}
}
