package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database       DB
	Redis          Redis
	Spotify        Spotify
	Env            string   `yaml:"env" env:"ENV" envDefault:"development"`
	MigrationsPath string   `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
	ListenAddr     string   `yaml:"listen_addr" env:"LISTEN_ADDR" env-default:":8080"`
	JWTKey         string   `yaml:"jwt_key" env:"JWT_KEY" env-required:"true"`
	AllowedOrigins []string `yaml:"allowed_origins" env:"ALLOWED_ORIGINS" env-default:"http://localhost:5173"`
	FrontendURL    string   `yaml:"frontend_url" env:"FRONTEND_URL" env-default:"http://localhost:5173"`
}

type Spotify struct {
	ClientID     string `yaml:"client_id" env:"SPOTIFY_CLIENT_ID" env-required:"true"`
	ClientSecret string `yaml:"client_secret" env:"SPOTIFY_CLIENT_SECRET" env-required:"true"`
	RedirectURL  string `yaml:"redirect_url" env:"SPOTIFY_REDIRECT_URL" env-required:"true"`
}

type Redis struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR" env-required:"true"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
}

type DB struct {
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"name" env:"DB_NAME" env-required:"true"`
	Driver   string `yaml:"driver" env:"DB_DRIVER" env-default:"postgres"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

func (cfg DB) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) IsProd() bool {
	return cfg.Env == "production"
}
