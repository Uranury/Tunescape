package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBConfig
	MigrationsPath string   `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
	RedisAddr      string   `yaml:"redis_addr" env:"REDIS_ADDR" env-required:"true"`
	ListenAddr     string   `yaml:"listen_addr" env:"LISTEN_ADDR" env-default:":8080"`
	JWTKey         string   `yaml:"jwt_key" env:"JWT_KEY" env-required:"true"`
	AllowedOrigins []string `yaml:"allowed_origins" env:"ALLOWED_ORIGINS" env-default:"http://localhost:5173"`
}

type DBConfig struct {
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"name" env:"DB_NAME" env-required:"true"`
	Driver   string `yaml:"driver" env:"DB_DRIVER" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

func (cfg DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
