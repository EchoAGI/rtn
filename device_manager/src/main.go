package main


import (
	"os"
	"os/signal"
	"syscall"

	"logic"
	"server"
	"api"
	"bus"
	"data"

	log "github.com/sirupsen/logrus"
)




func main() {
	appid := "LTAI82fei8OjVVIU"
	appsecret := "ow2HnGumKHdbtwmGVoGliE4D6peyNJ"
	groupId := "GID_wwj_dvc"
	clientId := groupId + "@@@" + "edison"

	token := bus.GenToken(groupId, appsecret)
	log.Infof("token: %s\n", token)

	topic := "device-control"

	config := &data.Config{
		MQTT: data.MQTTConfig{
			ClientId : clientId,
			Username: appid,
			Token: token,
			Host: "mqtt-cn-4590dvwb801.mqtt.aliyuncs.com",
			Port: 1883,
			Topic: topic,
			SSLEnabled: false,
			SSLCertPath: "",
		},
	}


	api := api.NewAPIImpl()
	busManager := bus.NewLMQBusManager(config)
	err := busManager.Start()
	if err != nil {
		log.Error(err)
	}

	hubManager := logic.NewHubManager()
	codec := logic.NewJsonCodec()

	s := server.NewMQServer(api, busManager, hubManager, codec)

	s.Serve()


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("wait for running")

	// Wait for termination signal
	select {
	case <-sigChan:
		log.Println("Received SIGTERM, the service is closing.")
	}
}