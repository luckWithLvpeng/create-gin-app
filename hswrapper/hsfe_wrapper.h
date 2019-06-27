/*********************************************************** 
* Date: 2019-05-07
* 
* Author: wangxinyu
* 
* Email: wangxinyu@alleyes.com.cn
* 
* Department: R&D
* 
* Company: Alleyes
* 
* Module: wrapper
* 
* Brief: wrapper layer between go and c++
* 
* Note: 
************************************************************/ 
#ifndef __HSFE_WRAPPER_H__
#define __HSFE_WRAPPER_H__

#define HS_DLL_INTERFACE
#define HS_STDCALL

#include "hsvd_interface.h" /// < 解码器 
#include "hsfe_interface.h" /// < 引擎 
#ifdef USE_RTSP
    #include "hsrs_interface.h" /// < RTSP Server
#endif
#include "hsrtmp_interface.h" /// < RTMP
#include "gst_webrtc_interface.h"
#include "gst_videowriter_interface.h"

#ifdef __cplusplus
extern "C" {
#endif

#define HSFE_WRAPPER_EC_NO_ERROR 0
#define HSFE_WRAPPER_EC_ENGINE_ERROR 1
#define HSFE_WRAPPER_EC_DECODER_ERROR 2
#define HSFE_WRAPPER_EC_RTSP_SERVER_ERROR 3
#define HSFE_WRAPPER_EC_RTMP_ERROR 4
#define HSFE_WRAPPER_EC_HOBOT_IPC_ERROR 5



    HS_DLL_INTERFACE int create(const HSFE_EngineParameters*);
    HS_DLL_INTERFACE void destroy();

    HS_DLL_INTERFACE int add_channel(const HSVD_DecoderParameters*, const HSFE_ChannelParameters*);
    HS_DLL_INTERFACE void remove_channel(const int channel_id);

#ifdef USE_RTSP
    HS_DLL_INTERFACE int add_rtsp_server(HSRS_Parameters *param);
    HS_DLL_INTERFACE void remove_rtsp_server(const int channel_id);
    HS_DLL_INTERFACE void push_rtsp_and_draw_bbox(const int channel_id, HSFE_Image *img, HSFE_FaceRect *p_rects, const int rect_count);
#endif
    HS_DLL_INTERFACE int add_rtmp_channel(const HSVD_DecoderParameters*);
    HS_DLL_INTERFACE void remove_rtmp_channel(const int channel_id);
    HS_DLL_INTERFACE void push_rtmp_and_draw_bbox(const int channel_id, HSFE_Image *img, HSFE_FaceRect *p_rects, const int rect_count);

    HS_DLL_INTERFACE int add_webrtc_channel(const GST_WEBRTC_Parameters*, PFUNC_GST_WEBRTC_SAMPLE_ARRIVE cb, const char *p_url, int use_rtmp);
    HS_DLL_INTERFACE void remove_webrtc_channel(const int channel_id);
    HS_DLL_INTERFACE void push_webrtc_and_draw_bbox(const int channel_id, HSFE_Image *img, HSFE_FaceRect *p_rects, const int rect_count);

    HS_DLL_INTERFACE int add_videowriter_channel(const GST_VIDEO_WRITER_Parameters *p_params);
    HS_DLL_INTERFACE void remove_videowriter_channel(const int channel_id);
    HS_DLL_INTERFACE void push_videowriter_and_draw_bbox(const int channel_id, HSFE_Image *img, HSFE_FaceRect *p_rects, const int rect_count);

    HS_DLL_INTERFACE int get_decoder_error_code();
    HS_DLL_INTERFACE int get_engine_error_code();

#ifdef __cplusplus
}
#endif

#endif
