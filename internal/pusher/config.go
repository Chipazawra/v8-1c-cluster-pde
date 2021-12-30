package pusher

type PusherConfig struct {
	PUSH_INTERVAL int    `env:"PUSH_INTERVAL" envDefault:"500"`
	PUSH_HOST     string `env:"PUSH_HOST" envDefault:"localhost"`
	PUSH_PORT     string `env:"PUSH_PORT" envDefault:"9091"`
}
