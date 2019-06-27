/*********************************************************** 
* Date: 2018-09-30
* 
* Author: Wang Xinyu
* 
* Email: wangxinyu@alleyes.com.cn
* 
* Department: 
* 
* Company: 
* 
* Module: gst-webrtc
* 
* Brief: push stream from appsrc to udpsrc
* 
* Note:
*       
* 
* CodePage: Pure UTF-8 
************************************************************/ 
#ifndef __GST_WEBRTC_INTERFACE_H__
#define __GST_WEBRTC_INTERFACE_H__


#ifdef __cplusplus
extern "C" {
#endif
    typedef enum _GST_WEBRTC_ErrorCode {
        GST_WEBRTC_EC_UNKNOWN_ERROR = (-1),
        GST_WEBRTC_EC_NO_ERROR = 0,
        GST_WEBRTC_EC_INVALID_PARAMETER = 1,
        GST_WEBRTC_EC_EXCEPTION = 2
    } GST_WEBRTC_ErrorCode;

    typedef enum _GST_WEBRTC_FrameFormat {
        GST_WEBRTC_FRAME_FORMAT_INVALID = 0,
        GST_WEBRTC_FRAME_FORMAT_BGR = 1,
        GST_WEBRTC_FRAME_FORMAT_NV12 = 2,
        GST_WEBRTC_FRAME_FORMAT_I420 = 3
    } GST_WEBRTC_FrameFormat;

    typedef struct _GST_WEBRTC_FaceRect {
        int channel_id;
        int frame_id;
        int face_id;

        int left;
        int top;
        int right;
        int bottom;

        int confidence;
    } GST_WEBRTC_FaceRect;

    typedef struct _GST_WEBRTC_Image {
        const void *p_buf;
        int size;
        int width;
        int height;
        int format;
    } GST_WEBRTC_Image;

    /// < parameters
    typedef struct _GST_WEBRTC_Parameters {
        int id;
        int format; /// < only support BGR and I420 for now.
        int width;
        int height;
        int size;
        int fps;
        int compress_ratio;
        int port;
    } GST_WEBRTC_Parameters;

    typedef struct _GST_WEBRTC_Sample {
        int id;
        void *p_buf;
        int size;
        int duration;
    } GST_WEBRTC_Sample;

    typedef void(*PFUNC_GST_WEBRTC_SAMPLE_ARRIVE)(const GST_WEBRTC_Sample*);

    void* gst_webrtc_create(const GST_WEBRTC_Parameters*, int *p_error_code);
    void gst_webrtc_destroy(void *p_webrtc);
    void gst_webrtc_set_sample_arrive_callback(void *p_webrtc, PFUNC_GST_WEBRTC_SAMPLE_ARRIVE cb);

    int gst_webrtc_start(void *p_source);
    int gst_webrtc_pause(void *p_source);
    int gst_webrtc_stop(void *p_source);
    void gst_webrtc_push(void *p_source, void *data);
    void gst_webrtc_push_and_draw_bbox(void *p_source, GST_WEBRTC_Image *img, const GST_WEBRTC_FaceRect *p_rects, const int rect_count);



#ifdef __cplusplus
}
#endif

#endif
