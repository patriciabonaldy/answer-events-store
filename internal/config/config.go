package config

import "time"

type Config struct {
	URI             string
	DatabaseName    string
	User            string
	Password        string
	ConnectTimeout  time.Duration
	MinPoolSize     uint64
	MaxPoolSize     uint64
	MaxConnIdleTime time.Duration
}
