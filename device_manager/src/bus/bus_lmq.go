package bus

import (
	"time"
	"fmt"
	"errors"
	"encoding/json"

	"io/ioutil"
	"crypto/tls"
	"crypto/x509"


	log "github.com/sirupsen/logrus"

	"github.com/eclipse/paho.mqtt.golang"

	"messages"
	"data"
)



// Connection is the high-level message bus interface
type Connection interface {
	Connect(options *ConnectionOptions) error
	Disconnect() error
	Publish(topic string, payload []byte) error
	Subscribe(topic string, handler SubscriptionHandler) error
}


var errorBadTLSCert = errors.New("Bad TLS certificate")

// Event describes different events which can happen over the
// life of a connection. Currently only BusConnection is supported.
type Event int

const (
	// ConnectedEvent indicates a working bus connection has been
	// established
	ConnectedEvent Event = iota
)


// EventHandler is called when a bus event occurs.
type EventHandler func(conn Connection, event Event)


type SubscriptionHandler func(conn Connection, topic string, message []byte)

// DisconnectMessage is sent when the connection is broken
type DisconnectMessage struct {
	Topic string
	Body  string
}



// ConnectionOptions describe how to configure a bus.Connection
type ConnectionOptions struct {
	ClientId	  string
	Username      string
	Password      string
	Host          string
	Port          int
	Topic 		  string
	SSLEnabled    bool
	SSLCertPath   string
	EventsHandler EventHandler
	AutoReconnect bool
	OnDisconnect  *DisconnectMessage
}



type MQTTConnection struct {
	options *ConnectionOptions
	conn mqtt.Client
	backoff *Backoff
}

func (mqc *MQTTConnection) Connect(options *ConnectionOptions) error {
	mqttOpts := mqc.buildMQTTOptions(options)
	if err := configureSSL(options, mqttOpts); err != nil {
		return err
	}

	if options.OnDisconnect != nil {
		//compressed := snappy.Encode(nil, []byte(options.OnDisconnect.Body))
		mqttOpts.SetWill(options.OnDisconnect.Topic, string(options.OnDisconnect.Body), 1, false)
	}

	mqc.backoff = NewBackoff()
	mqc.conn = mqtt.NewClient(mqttOpts)
	for {
		log.Errorf("prepare connect to %s", brokerURL(options))

		if token := mqc.conn.Connect(); token.Wait() && token.Error() != nil {
			log.Errorf("Error connecting to %s: %s", brokerURL(options), token.Error())
			mqc.backoff.Wait()
		} else {
			mqc.backoff.Reset()
			break
		}
	}

	log.Println("connected")

	mqc.options = options
	if mqc.options.EventsHandler != nil {
		mqc.options.EventsHandler(mqc, ConnectedEvent)
	}
	return nil
}


// Disconnect is required by the bus.Connection interface
func (mqc *MQTTConnection) Disconnect() error {
	mqc.conn.Disconnect(1000)
	return nil
}

// Publish is required by the bus.Connection interface
func (mqc *MQTTConnection) Publish(topic string, payload []byte) error {
	//compressed := snappy.Encode(nil, payload)
	token := mqc.conn.Publish(topic, 1, false, payload)
	token.Wait()
	return token.Error()
}

// Subscribe is required by the bus.Connection interface
func (mqc *MQTTConnection) Subscribe(topic string, handler SubscriptionHandler) error {
	mqttHandler := func(client mqtt.Client, message mqtt.Message) {
		payload := message.Payload()
		/*
		payload, err := snappy.Decode(nil, compressed)
		if err != nil {
			log.Errorf("Decompressing MQTT payload failed: %s", err)
			return
		}
		*/
		handler(mqc, message.Topic(), payload)
	}

	token := mqc.conn.Subscribe(topic, 1, mqttHandler)
	token.Wait()
	return token.Error()
}

func (mqc *MQTTConnection) disconnected(cilent mqtt.Client, err error) {
	log.Errorf("MQTT connection failed: %s.", err)
	for {
		if token := mqc.conn.Connect(); token.Wait() && token.Error() != nil {
			log.Errorf("Error connecting to %s: %s", brokerURL(mqc.options), token.Error())
			mqc.backoff.Wait()
		} else {
			mqc.backoff.Reset()
			break
		}
	}
	if mqc.options.EventsHandler != nil {
		mqc.options.EventsHandler(mqc, ConnectedEvent)
	}
}


func (mqc *MQTTConnection) buildMQTTOptions(options *ConnectionOptions) *mqtt.ClientOptions {
	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.SetAutoReconnect(options.AutoReconnect)
	mqttOpts.SetKeepAlive(time.Duration(60) * time.Second)
	mqttOpts.SetPingTimeout(time.Duration(15) * time.Second)
	mqttOpts.SetClientID(options.ClientId)
	mqttOpts.SetUsername(options.Username)
	mqttOpts.SetPassword(options.Password)
	mqttOpts.SetCleanSession(true)
	brokerURL := brokerURL(options)
	mqttOpts.AddBroker(brokerURL)

	if !options.AutoReconnect {
		mqttOpts.SetConnectionLostHandler(mqc.disconnected)
	}
	return mqttOpts
}

func configureSSL(options *ConnectionOptions, mqttOpts *mqtt.ClientOptions) error {
	if !options.SSLEnabled {
		return nil
	}
	log.Info("SSL enabled on MQTT connection to Cog")
	if options.SSLCertPath == "" {
		log.Warn("TLS certificate verification disabled.")
		mqttOpts.TLSConfig = tls.Config{
			InsecureSkipVerify: true,
		}
	} else {
		buf, err := ioutil.ReadFile(options.SSLCertPath)
		if err != nil {
			log.Errorf("Error reading TLS certificate file %s: %s.",
				options.SSLCertPath, err)
			return err
		}
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(buf)
		if !ok {
			log.Errorf("Failed to parse TLS certificate file %s.",
				options.SSLCertPath)
			return errorBadTLSCert
		}
		log.Info("TLS certificate verification enabled.")
		mqttOpts.TLSConfig = tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            roots,
		}
	}
	return nil
}

func brokerURL(options *ConnectionOptions) string {
	prefix := "tcp"
	if options.SSLEnabled {
		prefix = "ssl"
	}
	return fmt.Sprintf("%s://%s:%d", prefix, options.Host, options.Port)
}


type LMQBusManager struct {
	connOpts *ConnectionOptions
	conn *MQTTConnection
}

func NewLMQBusManager(config *data.Config) *LMQBusManager {
	connOpts := &ConnectionOptions{
		ClientId:	   config.MQTT.ClientId,
		Username:      config.MQTT.Username,
		Password:      config.MQTT.Token,
		Host:          config.MQTT.Host,
		Port:          config.MQTT.Port,
		Topic:		   config.MQTT.Topic,
		SSLEnabled:    config.MQTT.SSLEnabled,
		SSLCertPath:   config.MQTT.SSLCertPath,
	}

	bm := &LMQBusManager{
		connOpts : connOpts,
		conn : &MQTTConnection{},
	}

	connOpts.EventsHandler = bm.handleBusEvents

	connOpts.OnDisconnect = &DisconnectMessage{
		Topic: fmt.Sprintf("%s/%s", bm.connOpts.Topic, "discover"),
		Body:  newWill(config.MQTT.ClientId, fmt.Sprintf("bot/relays/%s/announcer", config.MQTT.ClientId)),
	}

	return bm
}

func (bm *LMQBusManager) Start() error {
	if err := bm.conn.Connect(bm.connOpts); err != nil {
		return err
	}

	return nil
}

func (bm *LMQBusManager) Subscribe(topic string, fun MQFunc) {
	bm.conn.Subscribe(topic, func(conn Connection, topic string, payload []byte) {
		fun(topic, payload)
	})
}

func (bm *LMQBusManager) UnSubscribe(topic string) {
}

func (bm *LMQBusManager) Publish(topic string, payload []byte) error {
	return bm.conn.Publish(topic, payload)
}








func (bm *LMQBusManager) handleBusEvents(conn Connection, event Event) {
	if event == ConnectedEvent {
		log.Infof("ConnectedEvent")

		topic := fmt.Sprintf("%s/%s", bm.connOpts.Topic, "wwj-1")

		log.Infof("topic: %s\n", topic)


		m := messages.BundleRef{
			Name: "Bunlde名称1",
			Version : "1.0",
		}

		data, _ := json.Marshal(m)
		err := bm.Publish(topic, data)
		if err != nil {
			log.Errorf("Publish Error: %v\n", err)
		}
		log.Infof("after send\n")


		/*
		bm.conn = conn

		if bm.announcer == nil {
			r.announcer = NewAnnouncer(r.config.ID, r.conn, r.catalog)
			if err := r.announcer.Run(); err != nil {
				log.Errorf("Failed to start announcer: %s.", err)
				panic(err)
			}
			if r.config.ManagedDynamicConfig == true {
				opts := r.makeConnOpts()
				r.dynConfigUpdater = NewDynamicConfigUpdater(r.config.ID, opts, r.config.DynamicConfigRoot,
					r.config.ManagedDynamicConfigRefreshDuration())
				if err := r.dynConfigUpdater.Run(); err != nil {
					log.Errorf("Failed to start bundle dynamic config updater: %s.", err)
					panic(err)
				}
			}
		} else {
			if err := r.announcer.SetSubscriptions(); err != nil {
				log.Fatalf("Failed to subscribe to required bundle announcement topics: %s.", err);
			}
			r.announcer.SendAnnouncement();
		}
		if err := r.setSubscriptions(); err != nil {
			log.Errorf("Failed to set Relay subscriptions: %s.", err)
			panic(err)
		}
		if r.catalog.Len() > 0 {
			r.catalog.Reconnected()
		}
		log.Info("Loading bundle catalog.")
		r.requestBundles()
		*/
	}
}



func newWill(id string, replyTo string) string {
	announcement := messages.NewOfflineAnnouncement(id, replyTo)
	data, _ := json.Marshal(announcement)
	return string(data)
}