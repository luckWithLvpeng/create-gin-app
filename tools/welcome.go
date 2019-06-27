package tools

import (
	"container/list"
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var queueTimeout, _ = beego.AppConfig.Float("QueueTimeout")
var welcomeMap = make(map[int]*WelcomeOper)
var welcomeFlag, _ = beego.AppConfig.Bool("WelcomeFlag")
var hmqtt = NewHMqtt()

type WelcomeMQTT struct {
	Image_url      string    `json:"image_url"`
	Channel        string    `json:"channel"`
	Name           string    `json:"name"`
	Music_url      string    `json:"music_url"`
	Background_url string    `json:"background_url"`
	Title          string    `json:"title"`
	Create_time    time.Time `json:"-"`
}

func (person *WelcomeMQTT) IsTimeout() bool {
	return time.Now().Sub(person.Create_time).Seconds() > queueTimeout
}

type VisitorMQTT struct {
	Capture_count int `json:"capture_count"`
	Hit_count     int `json:"hit_count"`
}

type WelcomeOper struct {
	personList *list.List
	ChannelId  int
}

func init() {
	if !welcomeFlag {
		return
	}
	hmqtt.SubscribeMqtt(topicCommon, MqttCallback)
}

func MqttCallback(client MQTT.Client, message MQTT.Message) {
	//beego.Info("enter welcome mqttcallback")
	var matchResult MatchResultMQTT
	err := json.Unmarshal(message.Payload(), &matchResult)
	if err != nil {
		beego.Error(err)
		return
	}
	// now : subscribe to mqtt_matchresult
	getWelcomeOper(matchResult.Channel_id).pushVisitorCount(VisitorMQTT{Hit_count: matchResult.Hit_count, Capture_count: matchResult.Capture_count})

	// hit
	if matchResult.Hit_flag {
		getWelcomeOper(matchResult.Channel_id).pushWelcome(WelcomeMQTT{
			Image_url: matchResult.Img_v,
			Channel:   matchResult.Label,
			Name:      matchResult.Name})
	}
}

func getWelcomeOper(channelId int) *WelcomeOper {
	_, ok := welcomeMap[channelId]
	if !ok {
		welcomeMap[channelId] = &WelcomeOper{ChannelId: channelId, personList: list.New()}
		welcomeMap[channelId].Begin()
	}

	return (welcomeMap[channelId])
}

func (this *WelcomeOper) Begin() {
	go func() {
		for {
			if this.personList.Len() == 0 {
				time.Sleep(50 * time.Millisecond)
				continue
			} else {
				if this.personList.Len() < 3 {
					time.Sleep(50 * time.Millisecond)
				}
				//dequeue
				this.publishWelcome()
				time.Sleep(3 * time.Second)
			}
		}
	}()
}

func (this *WelcomeOper) dequeue() interface{} {
	if this.personList.Len() == 0 {
		return nil
	}
	e := this.personList.Front()
	result := this.personList.Remove(e)
	return result
}

func (this *WelcomeOper) publishWelcome() {
	var result []WelcomeMQTT
	for len(result) < 3 {
		if this.personList.Len() == 0 {
			break
		}

		ele := this.dequeue().(WelcomeMQTT)
		if ele.IsTimeout() {
			continue
		} else {
			result = append(result, ele)
		}
	}

	if len(result) == 0 {
	} else {
		bs, _ := json.Marshal(result)
		hmqtt.PublishMqtt("welcome/"+strconv.Itoa(this.ChannelId), string(bs))
	}
}

func (this *WelcomeOper) pushWelcome(param WelcomeMQTT) {
	// process info
	param.Create_time = time.Now()
	param.Title = "终于等到你"
	param.Background_url = "/welcomeAsset/assets/normal.png"
	param.Music_url = "/welcomeAsset/assets/welcome/welcome-1.wav"
	//beego.Info("publish welcome ", param)
	this.personList.PushBack(param)
}

func (this *WelcomeOper) pushVisitorCount(param VisitorMQTT) {
	bs, _ := json.Marshal(param)
	//beego.Info("publish welcome visitor count ", param)
	hmqtt.PublishMqtt("welcome/visitor/"+strconv.Itoa(this.ChannelId), string(bs))
}
