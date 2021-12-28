package app

type Config struct {
	Host   string `env:"RAS_HOST" envDefault:"localhost"`
	Port   string `env:"RAS_PORT" envDefault:"1545"`
	User   string `env:"CLS_USER"`
	Pass   string `env:"CLS_PASS"`
	Expose string `env:"EXPOSE" envDefault:"9096"`
}
