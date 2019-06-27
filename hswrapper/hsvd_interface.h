/*********************************************************** 
* Date: 2017-10-24
* 
* Author: Wang Xinyu
* 
* Email: wangxinyu@hisign.com.cn
* 
* Department: 智能产品中心-研发部 
* 
* Company: 北京海鑫科金高科技股份有限公司 
* 
* Module: 海鑫视频解码器接口 
* 
* Brief: 可解码文件视频、RTSP视频、USB Camera和NV Camera 
* 
* Note: 此解码器基于gstreamer实现， 
*       使用前需提前安装好相关依赖 
* 
* CodePage: Pure UTF-8 
************************************************************/ 
#ifndef __HSVD_INTERFACE_H__
#define __HSVD_INTERFACE_H__

#define HSVD_STDCALL
#define HSVD_DLL_INTERFACE

#ifdef __cplusplus
extern "C" {
#endif

    /// < 错误码 
    typedef enum _HSVD_ErrorCode
    {
        HSVD_EC_UNKNOWN_ERROR = (-1),
        HSVD_EC_NO_ERROR = 0,
        HSVD_EC_INVALID_PARAMETER = 1,
        HSVD_EC_NOT_INITED = 2,
        HSVD_EC_ALREADY_INITED = 3,
        HSVD_EC_NV_SDK_ERROR = 4,
        HSVD_EC_THREAD_ERROR = 5, 
        HSVD_EC_EXCEPTION = 6, 
        HSVD_EC_CONVERT_IMG_FAILED = 7, 
        HSVD_EC_MAP_NVXIO_IMG_FAILED = 8,
        HSVD_EC_OPEN_FAILED = 9,
        HSVD_EC_CAPTURE_FAILED = 10,
        HSVD_EC_FORMAT_NOT_SUPPORTED = 11,
        HSVD_EC_GET_LOGGER_ERROR = 12
    }HSVD_ErrorCode;

    /// < 视频类型 
    typedef enum _HSVD_VideoType
    {
        HSVD_VIDEO_TYPE_INVALID = 0,
        HSVD_VIDEO_TYPE_FILE = 1,
        HSVD_VIDEO_TYPE_RTSP = 2,
        HSVD_VIDEO_TYPE_USB_CAMERA = 3,
        HSVD_VIDEO_TYPE_NV_CAMERA = 4
    }HSVD_VideoType;

    /// < 解码格式 
    typedef enum _HSVD_DecodeType
    {
        HSVD_DECODE_TYPE_INVALID = 0,
        HSVD_DECODE_TYPE_RGB = 1,
        HSVD_DECODE_TYPE_RGBA = 2,
        HSVD_DECODE_TYPE_YV12 = 3,
        HSVD_DECODE_TYPE_NV12 = 4,
        HSVD_DECODE_TYPE_I420 = 5,
        HSVD_DECODE_TYPE_BGR = 10
    }HSVD_DecodeType;

    /// < 解码器参数 用于添加视频接入 
    typedef struct _HSVD_DecoderParameters
    {
        int video_type; /// < 视频类型，参见枚举HSVD_VideoType 
        /** 
        * 文件视频：绝对路径 
        * RTSP: "rtsp://username:password@ip" 
        * USB Camera: 相机序号，比如0，1，2 
        * NV Camera: 相机序号，比如0，1，2 
        */ 
        const char *p_video_url;

        int loop;
        int decoder_id; /// < 解码器ID，由调用者提供，最好唯一 
        int decode_format; /// < 解码格式，参见枚举HSVD_DecodeType 

        /** 
        * 仅当视频类型为USB Camera或NV Camera时 
        * 视频的宽高和帧率的设置才有效 
        */ 

        int width;
        int height;
        int fps;
        int frame_skip_num;
    }HSVD_DecoderParameters;

    /// < 解码得到的帧 用于解码输出 
    typedef struct _HSVD_Frame
    {
        int error_code;
        int decoder_id;

        void *p_buf;
        int size;
        int format;

        int width;
        int height;
        long long frame_id;
    }HSVD_Frame;

    /// < 解码器事件 
    typedef struct _HSVD_Event
    {
        int decoder_id;
        int event_id;
    }HSVD_Event;

    typedef void(HSVD_STDCALL *PFUNC_HSVD_FRAME_ARRIVE)(const HSVD_Frame*, void *user_data);
    typedef void(HSVD_STDCALL *PFUNC_HSVD_EVENT_ARRIVE)(const HSVD_Event*);

    HSVD_DLL_INTERFACE int hsvd_environment_init();

    HSVD_DLL_INTERFACE void* hsvd_create(const HSVD_DecoderParameters*, int *p_error_code);
    HSVD_DLL_INTERFACE void hsvd_destroy(void *p_decoder);

    HSVD_DLL_INTERFACE int hsvd_get_decoder_id(void *p_decoder, int *p_decoder_id);

    HSVD_DLL_INTERFACE void hsvd_set_frame_arrive_callback(void *p_decoder, PFUNC_HSVD_FRAME_ARRIVE cb, void *user_data);
    HSVD_DLL_INTERFACE void hsvd_set_event_arrive_callback(void *p_decoder, PFUNC_HSVD_EVENT_ARRIVE cb);

    HSVD_DLL_INTERFACE int hsvd_start(void *p_decoder);
    HSVD_DLL_INTERFACE int hsvd_pause(void *p_decoder);
    HSVD_DLL_INTERFACE int hsvd_stop(void *p_decoder);



#ifdef __cplusplus
}
#endif

#endif
