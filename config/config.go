package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type(
	Config struct {
		Server Server
		Oauth Oauth
		Bot Bot
		Redis Redis
	}

	Server struct {
		Host string
		BaseUrl string
		MainSite string
		Ratelimit Ratelimit
		Session Session
	}

	Ratelimit struct {
		Window int
		Max int
	}

	Session struct {
		Database string
	}

	Oauth struct {
		Id int64
		Secret string
	}

	MariaDB struct {
		Host string
		Username string
		Password string
		Database string
		Threads int
	}

	Bot struct {
		Key string
		HttpServer []string
	}

	Redis struct {
		Uri string
	}
)

var(
	Conf Config
)

func LoadConfig() {
	raw, err := ioutil.ReadFile("config.toml"); if err != nil {
		panic(err)
	}

	_, err = toml.Decode(string(raw), &Conf); if err != nil {
		panic(err)
	}
}
