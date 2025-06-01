package main

type NetworkConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type DatabaseConfig struct {
	Engine   string `json:"engine"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type Config struct {
	Network  NetworkConfig  `json:"network"`
	Database DatabaseConfig `json:"database"`
}
