// middlelayer
package middlelayer

/*
#include "middlelayerc.h"
*/
import "C"

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"

	"eme/models"
	"eme/restapi"
	"eme/tools"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// 200 张图片测试一下硬盘空间
var stat syscall.Statfs_t
var fileSize = 500
var availMb uint64 = 0
var reservedSpace, _ = beego.AppConfig.Int("Reserved_space")
var warningTime = time.Now()
var (
	extract_feature_mutex  sync.Mutex
	monitor_callback_mutex sync.Mutex
	CR                     cache.Cache
)
var ChannelsLock sync.RWMutex

var Save_DB_Data chan models.Match_result
var Send_Match_Data chan tools.MatchResultMQTT

var SendExternalResultFlag, _ = beego.AppConfig.Bool("SendResultFlag")
var EngineErrorMap = map[int32]string{
	-1: "UNKNOWN ERROR, 未知错误 ",
	0:  "NO ERROR, 正确",
	1:  "INVALID PARAMETER, 无效参数",
	2:  "INVALID MODE, 无效模式 ",
	3:  "INVALID FEATURE SIZE, 无效的特征大小失败",
	4:  "ENGINE NOT CREATED, 引擎未创建",
	5:  "ENGINE ALREADY CREATED, 引擎未创建",
	6:  "MODULE NOT INIT, 模块未初始化",
	7:  "MODULE ALREADY INITED, 模块已初始化",
	8:  "CREATE LOGGER FAILED, 创建日志模块失败",
	9:  "FILE ERROR, 文件错误",
	10: "MEMORY ERROR, 内存错误",
	11: "SDK ERROR, 算法错误",
	12: "THREAD ERROR, 线程错误",
	13: "THREAD QUEUE FULL, 线程队列已满",
	14: "THREAD QUEUE DISABLED, 线程队列已禁用",
	15: "CHANNEL NOT EXIST, 无效的通道",
	16: "CHANNEL ALREADY EXIST, 通道已存在",
	17: "SUBLIB NOT EXIST, 分库不存在",
	18: "TEMPLATE NOT EXIST, 模板不存在",
	19: "TEMPLATE POOL IS FULL, 模板池已满",
	20: "NO FEATURE TO VERIFY, 核验输入不存在有效特征",
	21: "NO FACE IN IMG, 图片中无人脸",
	22: "EXCEED MAX CHANNEL COUNT, 超出通道最大值",
	23: "MULTI FACES IN IMG, 图片中有多张人脸",
	24: "FILE TOO SMALL, 图片太小建议大于4KB",
	25: "LOW QUALITY IMG, 图片质量过低",
	26: "IMG EXCEED MAX RESOLUTION, 图片分辨率超出3840 * 2160",
}
var DecoderErrorMap = map[int32]string{
	-1: "UNKNOWN ERROR, 未知错误",
	0:  "NO ERROR, 正确",
	1:  "INVALID PARAMETER, 无效参数",
	2:  "NOT INITED, 未初始化",
	3:  "ALREADY INITED, 已初始化",
	4:  "NV SDK ERROR, 英伟达解码器错误",
	5:  "THREAD ERROR, 线程异常",
	6:  "EXCEPTION, 异常",
	7:  "CONVERT IMG FAILED, 图像转换失败",
	8:  "MAP NVXIO IMG FAILED, NVXIO转换失败",
	9:  "OPEN FAILED, 打开失败",
	10: "CAPTURE FAILED, 获取失败",
	11: "FORMAT NOT SUPPORTED, 格式不支持",
}

//BEGIN......added by Jiefeng Lai / 2018.07.23
func SetExtractFeatureFromServerCallback() {
	logs.Debug(" ****** SetCallBackIdentifyFromServer in")
	C.SetCallBackFeatureExtraction()
	logs.Debug(" ****** SetCallBackIdentifyFromServer out")
}

//export ExtractFeatureFromServerWithGo
func ExtractFeatureFromServerWithGo(pImageBuffer *C.HSFE_Buffer, imgNum C.int, pFeatureBuffer *C.HSFE_Buffer) {
	var num = int(imgNum)
	var pFea = unsafe.Pointer(pFeatureBuffer)
	feaSize := unsafe.Sizeof(C.HSFE_Buffer{})

	for i := 0; i < num; i++ {
		p := (*C.HSFE_Buffer)(
			unsafe.Pointer(
				(uintptr(pFea) +
					feaSize*uintptr(i))))
		p.size = C.int(0)
	}
}

func SetIdentifyFromServerCallback() {
	logs.Debug(" ****** SetCallBackIdentifyFromServer in")
	C.SetCallBackIdentifyFromServer()
	logs.Debug(" ****** SetCallBackIdentifyFromServer out")
}

//export IdentifyFromServerWithGo
func IdentifyFromServerWithGo(subid C.int, pImage *C.char, imglen C.int, pCands *C.HSFE_PersonCandidate, pCandsCount *C.int) {
	fmt.Println("IdentifyFromServerWithGo...................................")
	var img = base64.StdEncoding.EncodeToString(
		[]byte(C.GoStringN((*C.char)(pImage),
			imglen)))

	var err_search error
	var resp *restapi.RecogResp
	if resp, err_search = restapi.HRZService.RecogInHrz(img, strconv.Itoa(int(subid)), int(*pCandsCount)); err_search != nil {
		fmt.Println("restapi.HRZService.RecogInHrz failed: " + err_search.Error())
		*pCandsCount = C.int(0)
		return
	}

	var num = 0
	var pcand = unsafe.Pointer(pCands)
	//CandSize := unsafe.Sizeof(*pCands)
	CandSize := unsafe.Sizeof(C.HSFE_PersonCandidate{})
	for i := 0; i < len(resp.Details); i++ {
		if _id, err := strconv.Atoi(resp.Details[i].Id); _id > 0 && err == nil {
			p := (*C.HSFE_PersonCandidate)(
				unsafe.Pointer(
					(uintptr(pcand) +
						CandSize*uintptr(i))))
			p.template_id = C.int(_id)
			//To be consistent with the score returned by interface under ARM64, get score = 1000-score*1000
			//p.score = C.float(resp.Details[i].Score * 1000)
			p.score = 1000.0 - C.float(resp.Details[i].Score*1000)
			num++
		}
	}

	*pCandsCount = C.int(num)
}

//END.//////////////////////////////////////////////////////////////////

func Init_ResultQueue() {
	Save_DB_Data = make(chan models.Match_result, models.EngineConfig.Save_DB_list_size)
	Send_Match_Data = make(chan tools.MatchResultMQTT, models.EngineConfig.Send_Match_buffer_size)

	go saveMatchResult()
	go sendMatchResult()
}

func saveMatchResult() {
	var count int64 = -1
	for {
		result := <-Save_DB_Data
		count++
		if count >= models.EngineConfig.DB_check_count || count == 0 {
			models.Del_MatchResult_()
			count = 0
		}
		id, _ := models.Add_MatchResult(&result)

		bytes, _ := base64.StdEncoding.DecodeString(result.CropFrame_data)
		// save to memory cache
		CR.Put("log"+strconv.Itoa(int(id)), bytes, time.Minute*2)
		incrCache("capture_count", int(result.Channel_id))
		incrCache("capture_count_today", -1)
		incrCache("capture_count_all", -1)

		// MQTT pub
		var score string
		var img_f string

		if int(result.Hit_count) == 0 { // 未比中

		} else if int(result.Hit_count) == 1 { // 比中
			img_f = fmt.Sprintf("/v1/feature/img/%d", int(result.Feature_id))
			score = fmt.Sprintf("%d 分", int(result.Score))
			incrCache("hit_count_all", -1)
			incrCache("hit_count_today", -1)
		}

		_, subType, _ := models.Get_Sublib_NameAndType(int(result.Sublib_id))
		match_result :=
			tools.MatchResultMQTT{
				Log_id:              int(id),
				Img_f:               img_f,
				Img_v:               fmt.Sprintf("/v1/log/img/%s", strconv.Itoa(int(id))),
				Time:                time.Now(),
				Score:               score,
				Label:               result.Channel_name,
				Name:                result.Feature_name,
				Type:                subType,
				Channel_id:          result.Channel_id,
				Capture_count:       int(getCache("capture_count", result.Channel_id)),
				Hit_count:           int(getCache("hit_count", result.Channel_id)),
				Hit_flag:            result.Hit_count == 1,
				Sublib_id:           result.Sublib_id,
				Sublib_name:         result.Sublib_name,
				Capture_count_today: getCache("capture_count_today", -1),
				Hit_count_today:     getCache("hit_count_today", -1),
				Capture_count_all:   getCache("capture_count_all", -1),
				Hit_count_all:       getCache("hit_count_all", -1),
				Face_left:           result.Face_left,
				Face_top:            result.Face_top,
				Face_right:          result.Face_right,
				Face_bottom:         result.Face_bottom,
				Face_id:             result.Face_id,
				Other_matches:       result.Other_matches}

		select {
		case <-time.After(time.Second * 2):
			logs.Error("Result sending timeout")
			for len(Send_Match_Data) > 0 {
				<-Send_Match_Data
			}
		case Send_Match_Data <- match_result:
		}
	}
}

func sendMatchResult() {
	var count int64 = -1
	for {
		result := <-Send_Match_Data
		send_json, _ := json.Marshal(result)
		tools.Publish(string(send_json))
		count++
	}
}

func incrCache(key string, channelId int) {
	if channelId > 0 {
		key = key + "_" + strconv.Itoa(channelId)
	}

	if !CR.IsExist(key) {
		now := time.Now()
		beginTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		endTime := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
		// get count
		total, hit_count := models.Count_MatchResult(channelId, beginTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
		// set memory cache
		timeout := endTime.Sub(time.Now())
		CR.Put("capture_count_"+strconv.Itoa(channelId), total, timeout)
		CR.Put("hit_count_"+strconv.Itoa(channelId), hit_count, timeout)

		total_today, hit_count_today := models.Count_MatchResult(-1, beginTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
		// set memory cache
		CR.Put("capture_count_today", total_today, timeout)
		CR.Put("hit_count_today", hit_count_today, timeout)

		startTime := time.Date(now.Year()-10, now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		total_all, hit_count_all := models.Count_MatchResult(-1, startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
		// set memory cache
		CR.Put("capture_count_all", total_all, timeout)
		CR.Put("hit_count_all", hit_count_all, timeout)

		//beego.Info(key, beginTime, endTime, total, hit_count, timeout)
	}
	CR.Incr(key)
}

func getCache(key string, channelId int) int64 {
	if channelId > 0 {
		key = key + "_" + strconv.Itoa(channelId)
	}
	if CR.IsExist(key) {
		i := CR.Get(key).(int64)
		return i
	} else {
		return 0
	}
}

func ClearCatchLogCount() {

	CR.Delete("capture_count_today")
	CR.Delete("hit_count_today")
	CR.Delete("capture_count_all")
	CR.Delete("hit_count_all")
}

func write_weed_query() (map[string]interface{}, error) {
	write_query_url := fmt.Sprintf("http://%s:%s/dir/assign", models.EngineConfig.Weed_host, models.EngineConfig.Weed_port)
	// resp, err := http.Get(write_query_url)
	resp, err := http.PostForm(write_query_url, url.Values{"ttl": {models.EngineConfig.Weed_ttl}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{}, 0)
	json.Unmarshal(result, &ret)
	if _, ok := ret["error"]; ok {
		return nil, errors.New(ret["error"].(string))
	}
	return ret, nil
}

func write_weed(addr, fid string, file io.Reader) (map[string]interface{}, error) {
	buf := &bytes.Buffer{}
	body_writer := multipart.NewWriter(buf)
	defer body_writer.Close()

	file_writer, err := body_writer.CreateFormFile("uploadfile", fid+".jpg")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(file_writer, file)
	if err != nil {
		return nil, err
	}

	content_type := body_writer.FormDataContentType()
	err = body_writer.Close()
	if err != nil {
		return nil, err
	}

	write_weed_url := fmt.Sprintf("http://%s/%s", addr, fid)
	write_resp, err := http.Post(write_weed_url, content_type, buf)
	if err != nil {
		return nil, fmt.Errorf("create new post request error: %s", write_weed_url)
	}
	defer write_resp.Body.Close()
	result, err := ioutil.ReadAll(write_resp.Body)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{}, 0)
	json.Unmarshal(result, &ret)
	if _, ok := ret["error"]; ok {
		return nil, errors.New("Error from server: " + ret["error"].(string))
	}
	return ret, nil
}

//export MonitorCallbackFromGo
func MonitorCallbackFromGo(result *C.MonitorResultFromC) {
	if int(result.Score) < models.EngineConfig.Monitor_quality {
		return
	}
	// 200 张图片测试一下硬盘空间
	if fileSize == 500 {
		syscall.Statfs("/", &stat)
		availMb = stat.Bavail * uint64(stat.Bsize) / MB
		fileSize = 0
	} else {
		fileSize++
	}
	// 小于700M 每次都检测
	if int(availMb) <= reservedSpace {
		syscall.Statfs("/", &stat)
		availMb = stat.Bavail * uint64(stat.Bsize) / MB
		if int(availMb) <= reservedSpace {
			// 30m  提示一次空间不足
			if time.Now().Sub(warningTime).Seconds() > 30 {
				send_json, _ := json.Marshal(&map[string]string{"msg": "硬盘空间不足" + strconv.Itoa(reservedSpace) + "M，请清理空间:Hard disk space less than " + strconv.Itoa(reservedSpace) + "M, clearing up more space"})
				tools.PublishWarning(string(send_json))
				warningTime = time.Now()
			}
			return
		}
	}

	logs.Debug(" ****** MonitorCallbackFromGo in %d", models.EngineConfig.Save_scene_hit)
	monitor_callback_mutex.Lock()
	defer monitor_callback_mutex.Unlock()

	if int(result.CropRet) != 0 {
		logs.Error(fmt.Sprintf("Crop image failed, ret is %d", int(result.CropRet)))
		logs.Debug(" ****** MonitorCallbackFromGo out")
		//return
	}

	var cropBytes []byte = C.GoBytes(result.CropFrame, result.CropFrameSize)
	var featureBytes []byte = C.GoBytes(result.CropFrameFeature, result.FeatureSize)

	var channel_name string
	var sublib_name string

	name := models.Get_Feature_Name(int(result.FeatureId))

	if int(result.HitFlag) == 0 { // 未比中

	} else if int(result.HitFlag) == 1 { // 比中
		incrCache("hit_count", int(result.ChannelId))
	}

	channel_name = models.Get_Channel_Name(int(result.ChannelId))
	sublib_name, _, _ = models.Get_Sublib_NameAndType(int(result.SublibId))

	other_matches := make([]models.Other_Match_, int(result.OtherCandidatesCount))

	for i, _ := range other_matches {
		other_matches[i].Feature_id = int(result.OtherFeatureIds[i])
		other_matches[i].Sublib_id = int(result.OtherSublibIds[i])
		other_matches[i].Score = int(result.OtherScores[i])
		other_matches[i].Hit_flag = (int(result.OtherHitFlags[i]) == 1)
		other_matches[i].Url = fmt.Sprintf("/v1/feature/img/%s", fmt.Sprintf("%d", int(result.OtherFeatureIds[i])))
	}
	other_matches_bytes, _ := json.Marshal(other_matches)

	// save db
	var frame_data_tmp string
	var frame_fid, frame_url string
	if models.EngineConfig.Save_scene_hit && int(result.HitFlag) == 1 || models.EngineConfig.Save_scene_not_hit && int(result.HitFlag) == 0 {
		// 现场原图 rgb/yuv->jpg
		var jpgBuf C.HSFE_Buffer

		var srcImage C.HSFE_Image
		srcImage.width = result.FrameWidth
		srcImage.height = result.FrameHeight
		srcImage.format = result.FrameFormat
		srcImage.p_buf = unsafe.Pointer(result.Frame)

		sceneJPGRet := C.MyCreateJPG(&srcImage, C.int(models.EngineConfig.Scene_reduce_percent), C.int(models.EngineConfig.JPG_quality), &jpgBuf, (*C.struct__HSFE_Rect)(nil))

		if jpgBuf.p_buf != (unsafe.Pointer)(nil) && sceneJPGRet == C.HSFE_EC_NO_ERROR {
			frame_data_tmp = base64.StdEncoding.EncodeToString(C.GoBytes(jpgBuf.p_buf, jpgBuf.size))
			C.MyFreeJPGBuf(jpgBuf.p_buf)
		} else {
			logs.Error(fmt.Sprintf("Create jpg failed, ret is %d", int(sceneJPGRet)))
		}
	}

	var crop_data_tmp string
	var crop_frame_fid, crop_frame_url string
	if models.EngineConfig.Use_weed {
		crop_data_reader := bytes.NewReader(cropBytes)

		ret, err := write_weed_query()
		if err != nil {
			logs.Error(err.Error())
			return
		}
		// map[fid:7,0137cf8972 url:127.0.0.1:8888 publicUrl:127.0.0.1:8888 count:1]
		crop_frame_fid = ret["fid"].(string)
		crop_frame_url = ret["url"].(string)

		_, err = write_weed(crop_frame_url, crop_frame_fid, crop_data_reader)
		if err != nil {
			logs.Error(err.Error())
			return
		}
	} else {
		crop_data_tmp = base64.StdEncoding.EncodeToString(cropBytes)
	}

	result_log := models.Match_result{
		Feature_id:        int(result.FeatureId),
		Sublib_id:         int(result.SublibId),
		Score:             int(result.Score),
		Hit_count:         int(result.HitFlag),
		Face_id:           int(result.FaceId),
		Face_left:         int(result.FaceLeft),
		Face_top:          int(result.FaceTop),
		Face_right:        int(result.FaceRight),
		Face_bottom:       int(result.FaceBottom),
		Channel_id:        int(result.ChannelId),
		Frame_id:          int(result.FrameId),
		Frame_data:        frame_data_tmp,
		Frame_fid:         frame_fid,
		Frame_url:         frame_url,
		Frame_size:        len(frame_data_tmp),
		Feature_name:      name,
		Sublib_name:       sublib_name,
		Channel_name:      channel_name,
		CropFrame_data:    crop_data_tmp,
		CropFrame_fid:     crop_frame_fid,
		CropFrame_url:     crop_frame_url,
		CropFrame_size:    len(crop_data_tmp),
		CropFrame_feature: base64.StdEncoding.EncodeToString(featureBytes),
		Other_matches:     string(other_matches_bytes)}

	Save_DB_Data <- result_log
	logs.Debug(" ****** MonitorCallbackFromGo out")
}

//export TrackResultArriveFromGo
func TrackResultArriveFromGo(jpgBuf C.HSFE_Buffer, channelId C.int) {
	//logs.Debug(" ****** TrackResultArriveFromGo in , channel id is %d", int(channelId))
	var jpeg []byte = C.GoBytes(jpgBuf.p_buf, jpgBuf.size)
	stream.UpdateJPEG(jpeg, int(channelId))
	//logs.Info(" ****** TrackResultArriveFromGo out")
}

func GetEngineErrMessage(ret int32) (err_msg string) {
	err_msg = "Unknown error."
	if ret >= 0 && int(ret) < len(EngineErrorMap) {
		err_msg = fmt.Sprintf("Engine error, error code is %d, %s", ret, EngineErrorMap[ret])
	}
	return
}

func GetDecoderErrMessage(ret int32) (err_msg string) {
	err_msg = "Unknown error."
	if ret >= 0 && int(ret) < len(DecoderErrorMap) {
		err_msg = fmt.Sprintf("Decoder error, error code is %d, %s", ret, DecoderErrorMap[ret])
	}
	return
}

// 获取wrapper错误码
func detectWrapperErr(ret C.int) error {
	if ret != C.HSFE_WRAPPER_EC_NO_ERROR {
		var err error
		var err_msg string
		if ret == C.HSFE_WRAPPER_EC_ENGINE_ERROR { // 引擎错误
			err_msg = GetEngineErrMessage(int32(C.MyGetEngineErrorCode()))
		} else if ret == C.HSFE_WRAPPER_EC_DECODER_ERROR { // 解码器错误
			err_msg = GetDecoderErrMessage(int32(C.MyGetDecoderErrorCode()))
		}
		err = errors.New(err_msg)
		logs.Error(err.Error())
		return err
	} else {
		return nil
	}
}

// 获取错误码
func detectErr(ret C.int) error {
	if ret != C.HSFE_EC_NO_ERROR {
		var err error
		var err_msg string
		if ret == C.HSFE_EC_SDK_ERROR { // sdk错误
			ret = C.MyGetSDKErrorCode()
			err_msg = fmt.Sprintf("SDK error, error code is %d", int32(ret))
		} else { // 引擎错误
			err_msg = GetEngineErrMessage(int32(ret))
		}
		err = errors.New(err_msg)
		logs.Error(err.Error())
		return err
	} else {
		return nil
	}
}

// 设置监控结果返回回调函数
func SetMonitorCallback() {
	logs.Debug(" ****** SetMonitorCallback in")
	C.SetCallBackMonitorResult()
	logs.Debug(" ****** SetMonitorCallback out")
}

// 设置跟踪回调函数
func SetTrackCallback() {
	logs.Debug(" ****** SetTrackCallback in")
	C.SetCallBackTrackResult()
	logs.Debug(" ****** SetTrackCallback out")
}

// 创建引擎
func CreateEngine(batch_size int, SDKPath []byte) error {
	logs.Info(" ****** CreateEngine in")
	var engine_parameters C.HSFE_EngineParameters
	engine_parameters.batch_size = (C.int)(batch_size)
	engine_parameters.collect_time_out_us = 0
	engine_parameters.chrominance_thresh = (C.int)(models.EngineConfig.Chrominance_thresh)
	engine_parameters.small_face_score_decay = (C.int)(models.EngineConfig.Small_face_score_decay)
	engine_parameters.min_pic_size = (C.int)(models.EngineConfig.Min_pic_size)
	var ret C.int
	ret = C.MyEngineCreate(&engine_parameters)
	logs.Info(" ****** CreateEngine out")
	return detectWrapperErr(ret)
}

func UpdateMjpegSetting(resize int, quality int) error {
	ret := C.MyUpdateMjpegSetting(C.int(resize), C.int(quality))
	return detectErr(ret)
}

// 设置引擎特征最大数量
func SetFeatureMaxCount(max_count int) error {
	logs.Debug(" ****** SetFeatureMaxCount in")
	var ret C.int
	ret = C.MySetFeatureMaxCount((C.int)(max_count))
	logs.Debug(" ****** SetFeatureMaxCount out")
	return detectErr(ret)
}

// 销毁引擎
func DestroyEngine() {
	C.MyEngineDestory()
}

// 添加通道
func AddChannel(channel *models.Channel, sublibs []int) error {
	ChannelsLock.Lock()
	defer ChannelsLock.Unlock()
	logs.Info(" ****** AddChannel in")
	logs.Info(fmt.Sprintf(" ****** AddChannel's id is %d", channel.Id))

	param := models.Param_{
		DetectRectLeft:   0,
		DetectRectTop:    0,
		DetectRectRight:  0,
		DetectRectBottom: 0,
		Liveness_thresh:  0,
		Loop:             false,
		Record:           0,
	}
	json.Unmarshal([]byte(channel.Param), &param)

	var decoder_parameters C.HSVD_DecoderParameters
	decoder_parameters.video_type = C.HSVD_VIDEO_TYPE_FILE
	if param.Loop {
		decoder_parameters.loop = (C.int)(1)
	} else {
		decoder_parameters.loop = (C.int)(0)
	}
	decoder_parameters.decoder_id = (C.int)(channel.Id)
	if models.EngineConfig.Use_Yuv {
		decoder_parameters.decode_format = C.HSVD_DECODE_TYPE_I420
	} else {
		decoder_parameters.decode_format = C.HSVD_DECODE_TYPE_BGR
	}
	decoder_parameters.frame_skip_num = (C.int)(1)

	tmp := C.CString(param.Localpath)
	defer C.MyFree(unsafe.Pointer(tmp))
	decoder_parameters.p_video_url = tmp

	var channel_param C.HSFE_ChannelParameters
	channel_param.channel_id = (C.int)(channel.Id)
	channel_param.reduce_ratio = (C.int)(channel.Resize_multiple)
	channel_param.frame_width = (C.int)(channel.Image_width)
	channel_param.frame_height = (C.int)(channel.Image_height)
	if models.EngineConfig.Use_Yuv {
		channel_param.frame_format = C.HSFE_IMG_FORMAT_I420
	} else {
		channel_param.frame_format = C.HSFE_IMG_FORMAT_COLOR
	}
	channel_param.color_channel_sequence = C.HSFE_CCS_RGB
	channel_param.frame_keep_flag = 1
	channel_param.face_filte_count = (C.int)(models.EngineConfig.Face_filt_count)
	channel_param.face_min_size = (C.int)(channel.Face_min_size)
	channel_param.face_max_size = (C.int)(channel.Face_max_size)
	channel_param.face_confidence = 75
	channel_param.face_pose_estimate_flag = (C.int)(channel.Face_yaw_enable)
	channel_param.face_yaw_left = (C.int)(channel.Face_yaw_left)
	channel_param.face_yaw_right = (C.int)(channel.Face_yaw_right)
	channel_param.prefeature_flag = 0
	channel_param.recog_top_threshold = (C.int)(channel.Accept_max_threshold)
	channel_param.recog_bottom_threshold = channel_param.recog_top_threshold - 100
	channel_param.reject_flag = 0
	channel_param.reject_top_threshold = 0
	channel_param.reject_bottom_threshold = 0
	channel_param.channel_mode = (C.int)(models.EngineConfig.Channel_mode)
	channel_param.max_result_count = (C.int)(channel.Result_num)
	channel_param.merge_time_out_ms = (C.int)(channel.Face_collection_time)
	channel_param.liveness_thresh = (C.int)(param.Liveness_thresh)

	roi_width := 1920
	roi_height := 1080
	if !strings.HasPrefix(param.Localpath, "hobotipc://") {
		roi_width = channel.Image_width
		roi_height = channel.Image_height
	}
	/// detect rect edit by wangweiwei
	if param.DetectRectLeft < 0.0001 && param.DetectRectTop < 0.0001 &&
		param.DetectRectRight < 0.0001 && param.DetectRectBottom < 0.0001 {
		channel_param.detect_roi.left = (C.int)(10)
		channel_param.detect_roi.top = (C.int)(10)
		channel_param.detect_roi.right = (C.int)(roi_width - 10)
		channel_param.detect_roi.bottom = (C.int)(roi_height - 10)
	} else {

		left := int(float32(roi_width) * param.DetectRectLeft)
		top := int(float32(roi_height) * param.DetectRectTop)
		right := int(float32(roi_width) * param.DetectRectRight)
		bottom := int(float32(roi_height) * param.DetectRectBottom)

		if left > 0 && left < roi_width {
			channel_param.detect_roi.left = (C.int)(left)
		} else {
			channel_param.detect_roi.left = (C.int)(10)
		}

		if top > 0 && top < roi_height {
			channel_param.detect_roi.top = (C.int)(top)
		} else {
			channel_param.detect_roi.top = (C.int)(10)
		}

		if right > 0 && right < roi_width && right > left {
			channel_param.detect_roi.right = (C.int)(right)
		} else {
			channel_param.detect_roi.right = (C.int)(roi_width - 10)
		}

		if bottom > 0 && bottom < roi_height && bottom > top {
			channel_param.detect_roi.bottom = (C.int)(bottom)
		} else {
			channel_param.detect_roi.bottom = (C.int)(roi_height - 10)
		}
	}

	channel_param.sublib_count = (C.int)(len(sublibs))
	if len(sublibs) > 0 {
		var sublib_id []int32
		for i := 0; i < len(sublibs); i++ {
			sublib_id = append(sublib_id, int32(sublibs[i]))
		}
		channel_param.p_sublib_for_monitor = C.MyMallocIntArry(channel_param.sublib_count)
		C.MyMemcpyIntArry(channel_param.p_sublib_for_monitor, (*C.int)(unsafe.Pointer(&sublib_id[0])), channel_param.sublib_count)
		defer C.MyFree(unsafe.Pointer(channel_param.p_sublib_for_monitor))
	} else {
		channel_param.p_sublib_for_monitor = (*C.int)(nil)
	}

	decoder_parameters.width = channel_param.frame_width
	decoder_parameters.height = channel_param.frame_height
	fmt.Println("*******************************************************************************")
	fmt.Printf("%+v\n", decoder_parameters)
	fmt.Println("*******************************************************************************")
	fmt.Printf("%+v\n", channel_param)
	fmt.Printf("sublib : ")
	C.testIntArry(channel_param.p_sublib_for_monitor, channel_param.sublib_count)
	fmt.Println("*******************************************************************************")

	var ret C.int
	if models.EngineConfig.Webrtc {
		C.MySetUseWebrtcFlag((C.int)(1))
	} else {
		C.MySetUseWebrtcFlag((C.int)(0))
	}
	if models.EngineConfig.Use_Rtmp {
		C.MySetUseRtmpFlag((C.int)(1))
	} else {
		C.MySetUseRtmpFlag((C.int)(0))
	}

	ret = C.MyAddChannel(&decoder_parameters, &channel_param)
	logs.Info(" ****** AddChannel out")
	// record 时候保存记录视频 ，1 false ,2 true
	if param.Record > 1 {
		StartRecord(&ChannelParam{
			ID:     channel.Id,
			Name:   channel.Name,
			Height: channel.Image_height,
			Width:  channel.Image_width,
		})
	}
	return detectWrapperErr(ret)
}

// 删除通道
func RemoveChannel(channel_id int) {
	ChannelsLock.Lock()
	defer ChannelsLock.Unlock()
	logs.Info(" ****** RemoveChannel in")
	fmt.Println("****** RemoveChannel in")
	logs.Info(fmt.Sprintf(" ****** RemoveChannel's id is %d", channel_id))
	StopRecord(channel_id)
	logs.Info(fmt.Sprintf(" ****** RecordChannel's  out id is %d", channel_id))
	C.MyDelChannel((C.int)(channel_id))
	logs.Info(" ****** RemoveChannel out")
	fmt.Println("****** RemoveChannel out")
}

// 添加分库
func AddSublib(sublib_id int, quality_score int) error {
	logs.Info(" ****** AddSublib in, id is %d", sublib_id)
	var ret C.int
	ret = C.MyAddSublib((C.int)(sublib_id))
	if ret != C.HSFE_EC_NO_ERROR {
		logs.Info(" ****** AddSublib err out")
		return detectErr(ret)
	}
	ret = C.MySetSublibQuality((C.int)(sublib_id), (C.int)(quality_score))
	logs.Info(" ****** AddSublib out")
	return detectErr(ret)
}

// 删除分库
func RemoveSublib(sublib_id int) error {
	logs.Info(" ****** RemoveSublib in, id is %d", sublib_id)
	var ret C.int
	ret = C.MyDelSublib((C.int)(sublib_id))
	if ret == C.HSFE_EC_SUBLIB_NOT_EXIST {
		logs.Warning(" ****** RemoveSublib:HSFE_EC_SUBLIB_NOT_EXIST")
		return nil
	}
	logs.Info(" ****** RemoveSublib out")
	return detectErr(ret)
}

// 设置分库质量
func SetSublibQuality(sublib_id int, quality_score int) error {
	logs.Info(" ****** SetSublibQuality in, id is %d, quality_score is %d", sublib_id, quality_score)
	var ret C.int
	ret = C.MySetSublibQuality((C.int)(sublib_id), (C.int)(quality_score))
	logs.Info(" ****** SetSublibQuality out")
	return detectErr(ret)
}

// 添加特征
func AddFeature(feature []models.Feature, sublib_id int) error {
	logs.Debug(" ****** AddFeature in")
	var size int = len(feature)
	var ret C.int = (C.int)(C.HSFE_EC_NO_ERROR)
	if size > 0 {
		var templete []C.HSFE_FaceTemplate = make([]C.HSFE_FaceTemplate, size)
		for i := 0; i < size; i++ {
			templete[i].template_id = (C.int)(feature[i].Id)
			feaTmpsl, _ := base64.StdEncoding.DecodeString(feature[i].Feature)
			templete[i].feature_size = (C.int)(len(string(feaTmpsl)))
			tmp := C.CString(string(feaTmpsl))
			defer C.MyFree(unsafe.Pointer(tmp))
			templete[i].p_feature = unsafe.Pointer(tmp)
			templete[i].result = (C.int)(feature[i].Success)
		}
		logs.Debug(" ****** AddFeature's id is %d-%d", feature[0].Id, feature[size-1].Id)

		var batch_size C.int = (C.int)(size)
		var sublib_cid C.int = (C.int)(sublib_id)

		ret = C.MyAddFeature(&templete[0], batch_size, sublib_cid)
	}
	logs.Debug(" ****** AddFeature out")
	return detectErr(ret)
}

func VerifyTwoFeatures(feature1 string, feature2 string) (score int, err error) {
	logs.Debug(" ****** VerifyTwoFeatures in")
	var ret C.int = (C.int)(C.HSFE_EC_NO_ERROR)
	var float_score C.float = (C.float)(0)
	featureRaw1, _ := base64.StdEncoding.DecodeString(feature1)
	featureRaw2, _ := base64.StdEncoding.DecodeString(feature2)
	ret = C.MyVerify((*C.uchar)(unsafe.Pointer(&featureRaw1[0])), (*C.uchar)(unsafe.Pointer(&featureRaw2[0])), (*C.float)(&float_score))
	score = (int)(float_score * 1000.0)
	logs.Info(fmt.Sprintf(" ****** score: %d", score))
	logs.Debug(" ****** VerifyTwoFeatures out")

	err = detectErr(ret)
	return
}

func SearchOneFeature(feature string, sublibs []int, face_width int, confidence float32, result_count int, result_threshold int) (results []models.Other_Match_, err error) {
	logs.Debug(" ****** SearchOneFeature in")

	featureRaw, _ := base64.StdEncoding.DecodeString(feature)

	var ret C.int = (C.int)(C.HSFE_EC_NO_ERROR)

	var sublib_array []C.int = make([]C.int, len(sublibs))

	var result_scores []C.int = make([]C.int, result_count)
	var result_sublibs []C.int = make([]C.int, result_count)
	var result_featureids []C.int = make([]C.int, result_count)

	for i := 0; i < len(sublibs); i++ {
		sublib_array[i] = (C.int)(sublibs[i])
	}

	ret = C.MySearchOneTemplate((*C.uchar)(&featureRaw[0]), (C.int)(len(featureRaw)), (C.int)(face_width), (C.float)(confidence), (C.int)(len(sublibs)), &sublib_array[0], (*C.int)(unsafe.Pointer(&result_count)), &result_scores[0], &result_featureids[0], &result_sublibs[0])

	for i := 0; i < result_count; i++ {
		if result_threshold <= int(result_scores[i]) {
			var cur_result models.Other_Match_
			cur_result.Feature_id = int(result_featureids[i])
			cur_result.Sublib_id = int(result_sublibs[i])
			cur_result.Score = int(result_scores[i])
			results = append(results, cur_result)
		}
	}

	logs.Debug(" ****** SearchOneFeature out")
	err = detectErr(ret)
	return
}

func SearchOnePhoto(path string, sublibs []int, result_count int, result_threshold int) (results []models.Other_Match_, err error) {
	logs.Debug(" ****** SearchOnePhoto in")

	tmp := C.CString(path)
	defer C.MyFree(unsafe.Pointer(tmp))

	var ret C.int = (C.int)(C.HSFE_EC_NO_ERROR)

	var sublib_array []C.int = make([]C.int, len(sublibs))

	var result_scores []C.int = make([]C.int, result_count)
	var result_sublibs []C.int = make([]C.int, result_count)
	var result_featureids []C.int = make([]C.int, result_count)

	for i := 0; i < len(sublibs); i++ {
		sublib_array[i] = (C.int)(sublibs[i])
	}

	ret = C.MySearchOnePhoto((*C.char)(unsafe.Pointer(tmp)), (C.int)(len(sublibs)), &sublib_array[0], (*C.int)(unsafe.Pointer(&result_count)), &result_scores[0], &result_featureids[0], &result_sublibs[0])

	for i := 0; i < result_count; i++ {
		if result_threshold <= int(result_scores[i]) {
			var cur_result models.Other_Match_
			cur_result.Feature_id = int(result_featureids[i])
			cur_result.Sublib_id = int(result_sublibs[i])
			cur_result.Score = int(result_scores[i])
			results = append(results, cur_result)
		}
	}

	logs.Debug(" ****** SearchOnePhoto out")
	err = detectErr(ret)
	return
}

// 删除特征
func RemoveFeature(Id, sublibId int) error {
	logs.Info(" ****** RemoveFeature in")
	logs.Info(fmt.Sprintf(" ****** RemoveFeature's id is %d", Id))
	var template_id C.int = (C.int)(Id)
	var sublib_id C.int = (C.int)(sublibId)
	var ret C.int
	ret = C.MyDelFeature(template_id, sublib_id)

	logs.Info(" ****** RemoveFeature out")
	return detectErr(ret)
}

// 删除所有特征
func RemoveAllFeature() error {
	logs.Info(" ****** RemoveAllFeature in")
	var ret C.int
	ret = C.MyDelAllFeature()
	logs.Info(" ****** RemoveAllFeature out")
	return detectErr(ret)
}

// 特征提取
func ExtractFaceFeature(filePath []string, batchSize int) ([]models.Feature, error) {
	extract_feature_mutex.Lock()
	defer extract_feature_mutex.Unlock()
	logs.Info(" ****** ExtractFaceFeature in")
	var feature []models.Feature
	var feature_packet []C.HSFE_ExtractFeaturePacket
	var err error
	var detail []byte

	for i := 0; i < batchSize; i++ {
		var feature_tmp models.Feature
		feature_tmp.Name = strings.TrimSuffix(filepath.Base(filePath[i]), ".jpg")
		feature = append(feature, feature_tmp)
	}

	for i := 0; i < batchSize; i++ {
		detail, err = ioutil.ReadFile(filePath[i])
		if err == nil {
			tmp := C.CString(filePath[i])
			defer C.MyFree(unsafe.Pointer(tmp))
			var feature_packet_tmp C.HSFE_ExtractFeaturePacket
			feature_packet_tmp.p_img_file_name = tmp
			feature_packet_tmp.p_gray = unsafe.Pointer(nil)
			feature_packet_tmp.gray_width = 0
			feature_packet_tmp.gray_height = 0
			feature_packet_tmp.error_code = -1
			feature_packet_tmp.quality_score = 800
			feature_packet = append(feature_packet, feature_packet_tmp)

			// 大于200k值时进行图片压缩
			imageSize := len(detail)
			if imageSize < 1024*200 {
				feature[i].Image_data = base64.StdEncoding.EncodeToString(detail)
			} else {
				image, err := ResizeImage(filePath[i], imageSize)
				if err == nil {
					feature[i].Image_data = image.ImageBuf

				} else {
					feature[i].Image_data = base64.StdEncoding.EncodeToString(detail)

				}
			}

		} else {
			return feature, err
		}
	}

	var ret C.int
	ret = C.MyExtractFeature(&feature_packet[0], C.int(batchSize))

	if ret != C.HSFE_EC_NO_ERROR { // 失败
		err = detectErr(ret)
	}

	for j := 0; j < batchSize; j++ {
		if feature_packet[j].error_code != C.HSFE_EC_NO_ERROR {
			feature[j].Feature = ""
			feature[j].Quality_score = 0
			feature[j].Success = -int(feature_packet[j].error_code)
		} else {
			feature[j].Feature = base64.StdEncoding.EncodeToString(
				[]byte(C.GoStringN((*C.char)(feature_packet[j].p_feature),
					feature_packet[j].feature_size)))
			feature[j].Quality_score = (int)(feature_packet[j].quality_score)
			feature[j].Img_width = (int)(feature_packet[j].gray_width)
			feature[j].Img_height = (int)(feature_packet[j].gray_height)
			feature[j].Success = 1
		}
	}

	logs.Info(" ****** ExtractFaceFeature out")

	return feature, err
}

func GetCurrentFeatureCount() int32 {
	logs.Info(" ****** GetCurrentFeatureCount in")
	count := C.MyGetCurrentFeatureCount()
	logs.Info(" ****** GetCurrentFeatureCount out")
	return int32(count)
}

func CropImage(cropImage *CropImageStruct) error {
	logs.Info(" ****** CropImage in")
	var ret C.int
	var src_image C.HSFE_Image
	var rect C.HSFE_Rect

	tmp := C.CString(cropImage.SrcImg.ImageBuf)
	defer C.MyFree(unsafe.Pointer(tmp))
	src_image.p_buf = unsafe.Pointer(tmp)
	src_image.size = (C.int)(cropImage.SrcImg.Size)
	src_image.width = (C.int)(cropImage.SrcImg.Width)
	src_image.height = (C.int)(cropImage.SrcImg.Height)
	src_image.format = (C.int)(cropImage.SrcImg.Format)

	rect.left = (C.int)(cropImage.CropRect.Left)
	rect.top = (C.int)(cropImage.CropRect.Top)
	rect.right = (C.int)(cropImage.CropRect.Right)
	rect.bottom = (C.int)(cropImage.CropRect.Bottom)

	crop_width := cropImage.CropRect.Right - cropImage.CropRect.Left
	crop_height := cropImage.CropRect.Bottom - cropImage.CropRect.Top
	var sldesImage []byte = make([]byte, crop_width*crop_height*3, crop_width*crop_height*3)
	ret = C.MyCropImage(&src_image, &rect, unsafe.Pointer(&sldesImage[0]))
	logs.Info("****** CropImage out")
	if ret != C.HSFE_EC_NO_ERROR {
		return errors.New(strconv.Itoa(int(ret)))
	} else {
		cropImage.CropImg = string(sldesImage[:])
		return nil
	}
}

// 特征提取2.0
func ExtractImageFeature(filePath string) (models.ExtractFeatureInfo, error) {
	extract_feature_mutex.Lock()
	defer extract_feature_mutex.Unlock()
	logs.Info(" ****** ExtractImageFeature in")
	var feature models.ExtractFeatureInfo
	var feature_packet C.HSFE_ImageFeaturePacket
	var err error
	var detail []byte

	detail, err = ioutil.ReadFile(filePath)
	if err == nil {
		tmp := C.CString(filePath)
		defer C.MyFree(unsafe.Pointer(tmp))
		var feature_packet_tmp C.HSFE_ImageFeaturePacket
		feature_packet_tmp.p_img_file_name = tmp
		feature_packet_tmp.quality_score = 800
		feature_packet = feature_packet_tmp

	} else {
		return feature, err
	}

	var ret C.int
	ret = C.MyExtractImageFeature(&feature_packet)

	if ret != C.HSFE_EC_NO_ERROR { // 失败
		err = detectErr(ret)
	}

	if feature_packet.error_code != C.HSFE_EC_NO_ERROR {
		feature.Quality_score = 0
		feature.Error_code = -int(feature_packet.error_code)
		err_msg := fmt.Sprintf("返回特征错误，错误码为： %d", -int(feature.Error_code))
		err = errors.New(err_msg)
	} else {
		feature.Image_feature = base64.StdEncoding.EncodeToString(
			[]byte(C.GoStringN((*C.char)(feature_packet.p_feature),
				feature_packet.feature_size)))
		feature.Feature_size = (int)(feature_packet.feature_size)
		feature.Img_size = len(detail)
		feature.Quality_score = (int)(feature_packet.quality_score)
		feature.Img_width = (int)(feature_packet.img_width)
		feature.Img_height = (int)(feature_packet.img_height)
		feature.Face_left = (int)(feature_packet.face_left)
		feature.Face_top = (int)(feature_packet.face_top)
		feature.Face_width = (int)(feature_packet.face_width)
		feature.Face_height = (int)(feature_packet.face_height)
		feature.Is_gray = (int)(feature_packet.is_gray)
		feature.Pitch = (int)(feature_packet.pitch)
		feature.Roll = int(feature_packet.roll)
		feature.Ustddev = int(feature_packet.ustddev)
	}

	logs.Info(" ****** ExtractImageFeature out")

	return feature, err
}

func ResizeImage(filePath string, image_size int) (ima Image, err1 error) {
	var image Image
	var img unsafe.Pointer
	file := C.CString(filePath)
	defer C.MyFree(unsafe.Pointer(file))

	ret := C.MyJPGLoadColor(file, &img, (*C.int)(unsafe.Pointer(&image.Width)), (*C.int)(unsafe.Pointer(&image.Height)))
	if ret != C.HSFE_EC_NO_ERROR {
		return image, errors.New("load color error")
	}

	// 原图 rgb/yuv->jpg
	var jpgBuf C.HSFE_Buffer
	var srcImage C.HSFE_Image
	srcImage.width = C.int(image.Width)
	srcImage.height = C.int(image.Height)
	srcImage.format = C.HSFE_IMG_FORMAT_COLOR
	srcImage.p_buf = img

	//计算压缩比例
	var resize float64
	var percent int
	if image_size > 1024*1024 {
		resize = float64(1024*1024) / float64(image_size)
		percent = int(resize * 100)
	} else {
		resize = 1
		percent = 100
	}

	JPGRet := C.MyCreateJPG(&srcImage, C.int(percent), C.int(90), &jpgBuf, (*C.struct__HSFE_Rect)(nil))
	if jpgBuf.p_buf != (unsafe.Pointer)(nil) && JPGRet == C.HSFE_EC_NO_ERROR {
		image.ImageBuf = base64.StdEncoding.EncodeToString(C.GoBytes(jpgBuf.p_buf, jpgBuf.size))
		C.MyFreeJPGBuf(jpgBuf.p_buf)
		image.Height = int(resize * float64(image.Height))
		image.Width = int(resize * float64(image.Width))
		image.Size = int(jpgBuf.size)
	} else {
		logs.Error(fmt.Sprintf("Create jpg failed, ret is %d", int(JPGRet)))
		return image, errors.New(fmt.Sprintf("Create jpg failed, ret is %d", int(JPGRet)))
	}
	logs.Info("*********resize image ok********", int(C.HSFE_IMG_FORMAT_COLOR), percent, jpgBuf.size, image.Width, image.Height)
	return image, nil
}
