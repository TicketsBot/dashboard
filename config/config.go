package config

import (
	"github.com/BurntSushi/toml"
	"github.com/TicketsBot/common/sentry"
	"os"
	"strconv"
	"strings"
)

type (
	Config struct {
		Admins          []uint64
		ForceWhitelabel []uint64
		Debug           bool
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
		Id          uint64
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
		ProxyUrl              string `toml:"discord-proxy-url"`
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

	Referral struct {
		Show bool
		Link string
	}
)

var (
	Conf Config
)

func LoadConfig() {
	if _, err := os.Stat("config.toml"); err == nil {
		fromToml()
	} else {
		fromEnvvar()
	}
}

func fromToml() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		panic(err)
	}
}

func fromEnvvar() {
	var admins []uint64
	for _, id := range strings.Split(os.Getenv("ADMINS"), ",") {
		if parsed, err := strconv.ParseUint(id, 10, 64); err == nil {
			admins = append(admins, parsed)
		} else {
			sentry.Error(err)
		}
	}

	var forcedWhitelabel []uint64
	for _, id := range strings.Split(os.Getenv("FORCED_WHITELABEL"), ",") {
		if parsed, err := strconv.ParseUint(id, 10, 64); err == nil {
			forcedWhitelabel = append(forcedWhitelabel, parsed)
		} else {
			sentry.Error(err)
		}
	}

	rateLimitWindow, _ := strconv.Atoi(os.Getenv("RATELIMIT_WINDOW"))
	rateLimitMax, _ := strconv.Atoi(os.Getenv("RATELIMIT_MAX"))
	sessionThreads, _ := strconv.Atoi(os.Getenv("SESSION_DB_THREADS"))
	oauthId, _ := strconv.ParseUint(os.Getenv("OAUTH_ID"), 10, 64)
	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redisThreads, _ := strconv.Atoi(os.Getenv("REDIS_THREADS"))

	Conf = Config{
		Admins:          admins,
		ForceWhitelabel: forcedWhitelabel,
		Debug:           os.Getenv("DEBUG") != "",
		Server: Server{
			Host:     os.Getenv("SERVER_ADDR"),
			BaseUrl:  os.Getenv("BASE_URL"),
			MainSite: os.Getenv("MAIN_SITE"),
			Ratelimit: Ratelimit{
				Window: rateLimitWindow,
				Max:    rateLimitMax,
			},
			Session: Session{
				Threads: sessionThreads,
				Secret:  os.Getenv("SESSION_SECRET"),
			},
			Secret: os.Getenv("JWT_SECRET"),
		},
		Oauth: Oauth{
			Id:          oauthId,
			Secret:      os.Getenv("OAUTH_SECRET"),
			RedirectUri: os.Getenv("OAUTH_REDIRECT_URI"),
		},
		Database: Database{
			Uri: os.Getenv("DATABASE_URI"),
		},
		Bot: Bot{
			Token:                 os.Getenv("BOT_TOKEN"),
			PremiumLookupProxyUrl: os.Getenv("PREMIUM_PROXY_URL"),
			PremiumLookupProxyKey: os.Getenv("PREMIUM_PROXY_KEY"),
			ObjectStore:           os.Getenv("LOG_ARCHIVER_URL"),
			AesKey:                os.Getenv("LOG_AES_KEY"),
			ProxyUrl:              os.Getenv("DISCORD_PROXY_URL"),
		},
		Redis: Redis{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     redisPort,
			Password: os.Getenv("REDIS_PASSWORD"),
			Threads:  redisThreads,
		},
		Cache: Cache{
			Uri: os.Getenv("CACHE_URI"),
		},
	}
}
