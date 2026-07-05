package config

type Config struct {
	Environment EnvironmentConfig `yaml:"env"`
	Database    DatabaseConfig    `yaml:"database"`
	Storage     StorageConfig     `yaml:"storage"`
	Server      ServerConfig      `yaml:"server"`
}

type EnvironmentConfig struct {
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type StorageConfig struct {
	AssetsPath   string `yaml:"assets-path"`
	ContentsPath string `yaml:"contents-path"`
}

type ServerConfig struct {
	Port         string             `yaml:"port"`
	SessionStore SessionStoreConfig `yaml:"session-store"`
}

type SessionStoreConfig struct {
	PrivateKey string `yaml:"private-key"`
}
