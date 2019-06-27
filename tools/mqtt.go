// mqtt
package tools

import (
	"fmt"
	"sync"
	"time"
)

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	UUID "github.com/satori/go.uuid"
)

type MatchResultMQTT struct {
	Log_id              int       `json:"id"`
	Img_f               string    `json:"img_f"`
	Img_v               string    `json:"img_v"`
	Time                time.Time `json:"time"`
	Score               string    `json:"score"`
	Label               string    `json:"label"`
	Name                string    `json:"name"`
	Type                int       `json:"type"`
	Channel_id          int       `json:"channel_id"`
	Capture_count       int       `json:"capture_count"` // with only one channel_id
	Hit_count           int       `json:"hit_count"`     //with only one channel_id
	Hit_flag            bool      `json:"hit_flag"`
	Sublib_id           int       `json:"sublib_id"`
	Sublib_name         string    `json:"sublib_name"`
	Capture_count_today int64     `json:"capture_count_today"`
	Hit_count_today     int64     `json:"hit_count_today"`
	Capture_count_all   int64     `json:"capture_count_all"`
	Hit_count_all       int64     `json:"hit_count_all"`
	Other_matches       string    `json:"other_matches"`
	Face_id             int       `json:"face_id"`
	Face_left           int       `json:"face_left"`
	Face_top            int       `json:"face_top"`
	Face_right          int       `json:"face_right"`
	Face_bottom         int       `json:"face_bottom"`
}

var (
	hmqttCommon *hMqtt
	topicCommon = beego.AppConfig.String("MQTT_TOPIC")
	qos, _      = beego.AppConfig.Int("MQTT_QOS")
	wait_ms, _  = beego.AppConfig.Int("MQTT_WAITMS")
	retained    bool
)

type hMqtt struct {
	sync.Mutex
	client MQTT.Client
	opts   *MQTT.ClientOptions
}

func init() {
	hmqttCommon = NewHMqtt()
}

func Publish(payload string) {
	hmqttCommon.PublishMqtt(topicCommon, payload)
}

func PublishWarning(payload string) {
	hmqttCommon.PublishMqtt("/box/info", payload)
}

func Subscribe(callback MQTT.MessageHandler) {
	hmqttCommon.SubscribeMqtt(topicCommon, callback)
}
func PublishTopic(topic string, payload string) {
	hmqttCommon.PublishMqtt(topic, payload)
}
func NewHMqtt() *hMqtt {
	hmqtt := &hMqtt{}

	broker := beego.AppConfig.String("MQTT_BROKER")
	cleansess, _ := beego.AppConfig.Bool("MQTT_CLEANSESS")

	hmqtt.opts = MQTT.NewClientOptions()
	hmqtt.opts.AddBroker(broker)
	hmqtt.opts.SetCleanSession(cleansess)
	hmqtt.opts.SetWriteTimeout(time.Duration(3) * time.Second)
	hmqtt.opts.SetConnectTimeout(time.Duration(5) * time.Second)
	hmqtt.opts.SetKeepAlive(time.Duration(20) * time.Second)
	hmqtt.opts.SetPingTimeout(time.Duration(2) * time.Second)
	hmqtt.opts.SetAutoReconnect(true)

	hmqtt.connect()

	return hmqtt
}

func (this *hMqtt) connect() {
	this.opts.SetClientID(func() string {
		uuid := UUID.NewV1()
		sl_uuid := []byte(uuid.String())
		cliendId := string(sl_uuid[:20])
		return cliendId
	}())

	this.client = MQTT.NewClient(this.opts)
	if token := this.client.Connect(); token.Wait() && token.Error() != nil {
		logs.Error(fmt.Sprintf("Matchresult MQTT publisher start failed, error is %s\n", token.Error().Error()))
		return
	}

	logs.Info("Matchresult MQTT Publisher Started")
}

func (this *hMqtt) PublishMqtt(topic string, payload string) {
	//	logs.Debug("Publisher in")
	//	pub_mutex.Lock()
	//	defer pub_mutex.Unlock()
	//	if client.IsConnected() {
	//		token := client.Publish(topic, byte(qos), false, payload)
	//		token.Wait()
	//		if token.Error() != nil {
	//			logs.Error(token.Error().Error())
	//		} else {
	//			logs.Debug("Publisher success")
	//		}
	//	} else {
	//		logs.Error("MQTT reconnect...")
	//		connect()
	//	}

	//logs.Debug("Publisher in")
	right := make(chan bool)
	over := make(chan bool)
	this.Lock()
	defer this.Unlock()
	if this.client.IsConnected() {
		go func() {
			var stop = false
			for !stop {
				select {
				case <-right:
					over <- true
					stop = true
				case <-time.After(time.Duration(5) * time.Second):
					logs.Error("Publisher timeout, disconnect and reconnect.")
					if this.client.IsConnected() {
						this.client.Disconnect(uint(wait_ms))
					}
					this.connect()
					stop = true
					over <- true
				}
			}
		}()

		token := this.client.Publish(topic, byte(qos), false, payload)
		right <- true
		token.Wait()
		if token.Error() != nil {
			logs.Error(token.Error().Error())
			logs.Error("MQTT reconnect...")
			this.client.Disconnect(uint(wait_ms))
			this.connect()
		} else {
			logs.Debug("Publisher success.")
		}

		<-over
	} else {
		logs.Error("MQTT reconnect...")
		this.connect()
	}
}

func (this *hMqtt) SubscribeMqtt(topic string, callback MQTT.MessageHandler) {
	logs.Debug("Subscribe in")
	this.Lock()
	defer this.Unlock()
	sub := make(chan bool)
	defer close(sub)
	for true {
		if this.client == nil || !this.client.IsConnected() {
			logs.Error("MQTT reconnect...")
			this.connect()
			time.Sleep(time.Millisecond * 50)
		} else {
			break
		}
	}
	go func() {
		token := this.client.Subscribe(topic, byte(qos), callback)
		token.Wait()
		if token.Error() != nil {
			logs.Error(token.Error().Error())
			this.client.Disconnect(250)
		} else {
			logs.Debug("Subscribe success.")
		}
		sub <- true
	}()

	select {
	case <-time.After(time.Second * 5):
		logs.Error("Subscribe timeout.disconnect...")
		this.client.Disconnect(250)
	case <-sub:
	}
}
