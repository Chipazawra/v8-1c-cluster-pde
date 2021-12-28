package app

type Config struct {
	Host   string `env:"HOST_1C" envDefault:"localhost"`
	Port   string `env:"PORT_1C" envDefault:"1545"`
	User   string `env:"USER_1C"`
	Pass   string `env:"PASS_1C"`
	Expose string `env:"EXPOSE" envDefault:"9096"`
}
