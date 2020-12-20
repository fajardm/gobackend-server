package config

import "time"

type Config struct {
	App      App      `json:"app"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
	Redis    Redis    `json:"redis"`
}

type App struct {
	Name    string        `json:"name"`
	Version string        `json:"version"`
	Env     string        `json:"env"`
	Address string        `json:"address"`
	Timeout time.Duration `json:"timeout"`
}

func (c App) IsDevelopment() bool {
	return c.Env == "development"
}

type Redis struct {
	Addresses []string `json:"addresses"`
	Password  string   `json:"password"`
}

type SQL struct {
	Database      string `json:"database"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	MaxConnection int    `json:"maxConnection"`
}

type MongoDB struct {
	URI string `json:"uri"`
}

type Database struct {
	Type    string  `json:"type"`
	SQL     SQL     `json:"sql"`
	MongoDB MongoDB `json:"mongoDb"`
}

type Logger struct {
	Level string `json:"level"`
}
