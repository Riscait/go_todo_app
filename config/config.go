package config

import "github.com/caarlos0/env/v6"

// Config は環境変数から作成する設定情報の構造体。
// env:{環境変数名}
type Config struct {
	Env  string `env:"TODO_ENV"`
	Port int    `env:"PORT"`
}

// New returns a new Config struct.
func New() (*Config, error) {
	cfg := &Config{}
	// デフォルト値の設定がなければ必須項目として扱う。
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.Parse(cfg, opts); err != nil {
		return nil, err
	}
	return cfg, nil
}
