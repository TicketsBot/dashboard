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
		SentryDsn       string
		Server          Server
		Oauth           Oauth
		Database        Database
		Bot             Bot
		Redis           Redis
		Cache           Cache
		SecureProxyUrl  string
	}

	Server struct {
		Host           string
		MetricHost     string
		BaseUrl        string
		MainSite       string
		Ratelimit      Ratelimit
		Session        Session
		Secret         string
		RealIpHeaders  []string
		TrustedProxies []string
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
		Id                                   uint64
		Token                                string
		PremiumLookupProxyUrl                string `toml:"premium-lookup-proxy-url"`
		PremiumLookupProxyKey                string `toml:"premium-lookup-proxy-key"`
		ObjectStore                          string
		AesKey                               string `toml:"aes-key"`
		ProxyUrl                             string `toml:"discord-proxy-url"`
		RenderServiceUrl                     string `toml:"render-service-url"`
		ImageProxySecret                     string `toml:"image-proxy-secret"`
		PublicIntegrationRequestWebhookId    uint64 `toml:"public-integration-request-webhook-id"`
		PublicIntegrationRequestWebhookToken string `toml:"public-integration-request-webhook-token"`
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

// TODO: Proper env package
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
	botId, _ := strconv.ParseUint(os.Getenv("BOT_ID"), 10, 64)
	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redisThreads, _ := strconv.Atoi(os.Getenv("REDIS_THREADS"))
	publicIntegrationRequestWebhookId, _ := strconv.ParseUint(os.Getenv("PUBLIC_INTEGRATION_REQUEST_WEBHOOK_ID"), 10, 64)

	Conf = Config{
		Admins:          admins,
		ForceWhitelabel: forcedWhitelabel,
		Debug:           os.Getenv("DEBUG") != "",
		SentryDsn:       os.Getenv("SENTRY_DSN"),
		Server: Server{
			Host:       os.Getenv("SERVER_ADDR"),
			MetricHost: os.Getenv("METRIC_SERVER_ADDR"),
			BaseUrl:    os.Getenv("BASE_URL"),
			MainSite:   os.Getenv("MAIN_SITE"),
			Ratelimit: Ratelimit{
				Window: rateLimitWindow,
				Max:    rateLimitMax,
			},
			Session: Session{
				Threads: sessionThreads,
				Secret:  os.Getenv("SESSION_SECRET"),
			},
			Secret:         os.Getenv("JWT_SECRET"),
			TrustedProxies: strings.Split(os.Getenv("TRUSTED_PROXIES"), ","),
			RealIpHeaders:  strings.Split(os.Getenv("REAL_IP_HEADERS"), ","),
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
			Id:                                   botId,
			Token:                                os.Getenv("BOT_TOKEN"),
			PremiumLookupProxyUrl:                os.Getenv("PREMIUM_PROXY_URL"),
			PremiumLookupProxyKey:                os.Getenv("PREMIUM_PROXY_KEY"),
			ObjectStore:                          os.Getenv("LOG_ARCHIVER_URL"),
			AesKey:                               os.Getenv("LOG_AES_KEY"),
			ProxyUrl:                             os.Getenv("DISCORD_PROXY_URL"),
			RenderServiceUrl:                     os.Getenv("RENDER_SERVICE_URL"),
			ImageProxySecret:                     os.Getenv("IMAGE_PROXY_SECRET"),
			PublicIntegrationRequestWebhookId:    publicIntegrationRequestWebhookId,
			PublicIntegrationRequestWebhookToken: os.Getenv("PUBLIC_INTEGRATION_REQUEST_WEBHOOK_TOKEN"),
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
		SecureProxyUrl: os.Getenv("SECURE_PROXY_URL"),
	}
}
