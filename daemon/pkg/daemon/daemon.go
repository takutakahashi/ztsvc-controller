package daemon

type NetworkDaemon struct{}

type Config struct{}

func LoadConfig(path string) (Config, error) {
	return Config{}, nil
}

func NewDaemon(c Config) (NetworkDaemon, error) {
	return NetworkDaemon{}, nil
}

func (d NetworkDaemon) Start() error {
	return d.start()
}

func (d NetworkDaemon) start() error {
	return nil
}
