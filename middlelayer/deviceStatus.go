// middlelayer
package middlelayer

import (
	"sync"
	"time"

	"github.com/astaxie/beego/logs"

	"eme/models"
)

type ChannelIpcStatus struct {
	Channel_id     int
	Channel_status int
	Time           time.Time
}

var ChannelStatusMapLock sync.Mutex
var ChannelStatusMap map[int]*ChannelIpcStatus = make(map[int]*ChannelIpcStatus)

func SyncIpcStatus(channel_id int) {
	ChannelStatusMapLock.Lock()
	if val, ok := ChannelStatusMap[channel_id]; ok {
		val.Time = time.Now()
	} else {
		var Channels ChannelIpcStatus
		Channels.Channel_id = channel_id
		Channels.Time = time.Now()
		Channels.Channel_status = 0
		ChannelStatusMap[channel_id] = &Channels
	}
	ChannelStatusMapLock.Unlock()
}

func UpdateIpcDeviceStatus() {
	for {
		/// get通道
		channels := models.Get_Enabled_Channels()
		for i := 0; i < len(channels) && i < models.EngineConfig.Channel_max; i++ {
			ChannelStatusMapLock.Lock()
			if val, ok := ChannelStatusMap[channels[i].Id]; ok {
				timeNow := time.Now()
				if int(timeNow.Sub(val.Time).Seconds()) > 20 {
					if channels[i].Reject_threshold_mode != 1 {
						/// updata 1
						err := models.Change_Channel_Ipc_status(channels[i].Id, 1)
						if err != nil {
							logs.Info("change ipc status failed", err.Error())
							models.Change_Channel_Ipc_status(channels[i].Id, 1)
						}
					}
				} else {
					if channels[i].Reject_threshold_mode != 0 {
						/// updata 0
						err := models.Change_Channel_Ipc_status(channels[i].Id, 0)
						if err != nil {

							models.Change_Channel_Ipc_status(channels[i].Id, 0)
						}
					}
				}
			} else {
				var Channels ChannelIpcStatus
				Channels.Channel_id = channels[i].Id
				Channels.Time = time.Now()
				Channels.Channel_status = 0
				ChannelStatusMap[channels[i].Id] = &Channels
			}
			ChannelStatusMapLock.Unlock()
		}

		time.Sleep(time.Second * 8)
	}
}
