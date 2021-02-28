package config

// Logger writes to stdout
type Logger struct {
	Prod bool
}

// Config is struct which is filling by config from App path like /etc/app.yml
type Config struct {
	Confirmation string
	Socket       string `default:"/var/run/vkbotserver.sock"`
	Logger Logger
}
