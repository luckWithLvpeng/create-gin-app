// init
package middlelayer

import (
	"fmt"
	"eme/models"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/casbin/casbin"
)

var EngineLoaded = false
var GlobalCasbin *casbin.Enforcer
var stream *Stream

// 初始化用户权限配置
func InitCasbing() {
	GlobalCasbin = casbin.NewEnforcer("./rbac_model.conf", "./rbac_policy.csv")
}

// 初始化Cache
func Init_Cache() {
	var err error
	// 初始化 memory, check per 60s
	CR, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		panic(err)
	}
}

// 初始化mjpeg server
func Init_JPEG_Server() {
	stream = NewStream()
	mjpeg_httpport := beego.AppConfig.String("mjpeg_httpport")
	mjpeg_server := &http.Server{
		Addr:         fmt.Sprintf(":%s", mjpeg_httpport),
		Handler:      stream,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 21 * 60 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	//http.Handle("/camera", stream)
	err := mjpeg_server.ListenAndServe() //设置监听的端口
	if err != nil {
		panic(err)
	}
}

// 初始化加载所有分库、特征和通道信息，并添加到引擎中
func Init_Engine() {
	EngineLoaded = false
	var err error
	var count int64

	// 初始化引擎
	// 设置引擎回调函数
	SetMonitorCallback()
	SetTrackCallback()

	//BEGIN.............. - added by Jiefeng Lai
	if models.EngineConfig.Sdk_mode == models.MODE_RESTFUL { //Restful API mode
		SetIdentifyFromServerCallback()
		SetExtractFeatureFromServerCallback()
	}
	//END.//////////////////////////////////////

	// 启动引擎
	SDKPath := beego.AppConfig.String("SDK_PATH")
	batch_size, _ := beego.AppConfig.Int("BATCH_SIZE")
	err = CreateEngine(batch_size, []byte(SDKPath))
	if err != nil {
		panic(err)
	}

	UpdateMjpegSetting(models.EngineConfig.JPG_resize, models.EngineConfig.JPG_quality)
	SetFeatureMaxCount((models.EngineConfig.Feature_max << 9) & models.EngineConfig.Channel_max)

	models.FeatureLoadedSet(int32(models.Count_Sublib_Feature(0, 0)))
	var addedCount int = 0
	// 添加分库
	sublibs, count := models.Get_Sublib_List(0, 0)

	for i := 0; i < int(count); i++ {
		sublib_id := int(sublibs[i]["Id"].(int64))
		sublib_quality := int(sublibs[i]["Quality_score"].(int64))
		err = AddSublib(sublib_id, sublib_quality)
		if err != nil {
			panic(err)
		}

		if models.EngineConfig.Use_hobot {
		} else {
			if true {
				// 添加特征
				var feature_id_start int = 0
				if addedCount >= models.EngineConfig.Feature_max {
					break
				}
				for {
					features := models.Get_Sublib_Good_Features(sublib_id, feature_id_start)
					if len(features) != 0 {
						err = AddFeature(features, sublib_id)
						if err != nil {
							panic(err)
						}
						lens := len(features)
						addedCount += lens
						feature_id_start = features[lens-1].Id
						continue
					} else {
						break
					}
				}
			}
		}
	}

	// 添加通道
	channels := models.Get_Enabled_Channels()
	for i := 0; i < len(channels) && i < models.EngineConfig.Channel_max; i++ {
		if models.EngineConfig.Engine_id < 0 || (models.EngineConfig.Engine_id > 0 && models.EngineConfig.Engine_id != channels[i].Engine_id) {
			continue
		}
		AddChannel(channels[i], models.Get_Channel_Sublibs(channels[i].Id))
	}
	EngineLoaded = true

	go UpdateIpcDeviceStatus()
}

///  初始化重置数据库时 删除引擎所有特征
func Init_Syncdb_Engine() {
	// 启动引擎
	SDKPath := beego.AppConfig.String("SDK_PATH")
	batch_size, _ := beego.AppConfig.Int("BATCH_SIZE")
	err := CreateEngine(batch_size, []byte(SDKPath))
	if err != nil {
		panic(err)
	}

	// 添加分库
	sublibs, count := models.Get_Sublib_List(0, 0)

	for i := 0; i < int(count); i++ {
		sublib_id := int(sublibs[i]["Id"].(int64))
		sublib_quality := int(sublibs[i]["Quality_score"].(int64))
		err = AddSublib(sublib_id, sublib_quality)
		if err != nil {
			panic(err)
		}
	}

	/// 删除引擎所有特征
	err = RemoveAllFeature()
	if err != nil {
		panic(err)
	}
}
