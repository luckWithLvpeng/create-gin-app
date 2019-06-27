package tools

import (
	"github.com/astaxie/beego"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"fmt"
)

var (
	sublibAdd = "/v2/master/sublic/add"
	useArgos,_ = beego.AppConfig.Bool("argos")
	mqttClient = NewHMqtt()
	client *http.Client
)
func init() {
	if !useArgos {
		return
	}
	client = &http.Client{}
	mqttClient.SubscribeMqtt(sublibAdd, func(client MQTT.Client, message MQTT.Message) {
	beego.Info("enter welcome mqttcallback")
	result := message.Payload()
	fmt.Println(result)

	//r,err := DoBytesPost("POST","ww",result)

	})
}
