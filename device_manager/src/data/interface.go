package data



type MQTTConfig struct {
	ClientId string
	Username string
	Token string
	Host string
	Port int
	Topic string
	SSLEnabled bool
	SSLCertPath string
}



type Config struct {
	MQTT MQTTConfig
}

