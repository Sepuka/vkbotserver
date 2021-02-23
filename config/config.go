package config

type Config struct {
	Confirmation string
	Socket       string `default:"/var/run/vkbotserver.sock"`
}
