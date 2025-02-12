package configs

type (
	Config struct {
		Service       Service  `mapstructure:"service"`
		Database      Database `mapstructure:"database"`
		SpotifyConfig SpotifyConfig
	}

	Service struct {
		Port      string `mapstructure:"port"`
		SecretJWT string `mapstructure:"secretJWT"`
	}

	Database struct {
		DatabaseSourceName string `mapstructure:"dataSourceName"`
	}

	SpotifyConfig struct {
		ClientID     string
		ClientSecret string
	}
)
