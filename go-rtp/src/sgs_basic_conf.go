package rtp

type ConfigBasic struct {
	MaxCpus int
}

func (cfg *ConfigBasic) SetDefaultConf() {
	cfg.MaxCpus = 0
}

func (cfg *ConfigBasic) BasicConfigCheck() (err error) {
	return err
}
