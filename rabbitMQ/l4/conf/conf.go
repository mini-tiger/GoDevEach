package config

const (
	RMQADDR      = "amqp://test:123456@127.0.0.1:5672/testVH"
	EXCHANGENAME = "syslog_direct"
	CONSUMERCNT  = 4
)

var (
	RoutingKeys [4]string = [4]string{"info", "debug", "warn", "error"}
)
