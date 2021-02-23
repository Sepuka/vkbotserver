package config

type Logger struct {
	Prod bool
}

type Config struct {
	Confirmation string
	Socket       string `default:"/var/run/vkbotserver.sock"`
	Logger Logger
}
