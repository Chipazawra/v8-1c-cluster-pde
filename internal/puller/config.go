package puller

type PullerConfig struct {
	PULL_EXPOSE string `env:"PULL_EXPOSE" envDefault:"9096"`
}
