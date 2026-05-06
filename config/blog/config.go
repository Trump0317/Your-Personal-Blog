package config

type (
	Config struct {
		Log   Log   `yaml:"log"`
		Admin Admin `yaml:"admin"`
	}

	Log struct {
		Level string `yaml:"level"`
	}

	Admin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		APIKey   string `yaml:"api_key"`
	}
)
