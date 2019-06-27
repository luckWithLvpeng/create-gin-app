/*********************************************************** 
* Date: 2017-11-15
* 
* Author: Wang Xinyu
* 
* Email: wangxinyu@hisign.com.cn
* 
* Department: 智能产品中心-研发部 
* 
* Company: 北京海鑫科金高科技股份有限公司 
* 
* Module: rtmp
* 
* Brief: push stream from appsrc to rtmp
* 
* Note:
*       
* 
* CodePage: Pure UTF-8 
************************************************************/ 
#ifndef __HS_RTMP_INTERFACE_H__
#define __HS_RTMP_INTERFACE_H__


#ifdef __cplusplus
extern "C" {
#endif
    typedef enum _HSRTMP_ErrorCode {
        HSRTMP_EC_UNKNOWN_ERROR = (-1),
        HSRTMP_EC_NO_ERROR = 0,
        HSRTMP_EC_INVALID_PARAMETER = 1,
        HSRTMP_EC_EXCEPTION = 2
    } HSRTMP_ErrorCode;

    typedef enum _HSRTMP_FrameFormat {
        HSRTMP_FRAME_FORMAT_INVALID = 0,
        HSRTMP_FRAME_FORMAT_BGR = 1,
        HSRTMP_FRAME_FORMAT_NV12 = 2,
        HSRTMP_FRAME_FORMAT_I420 = 3
    } HSRTMP_FrameFormat;

    typedef struct _HSRTMP_FaceRect {
        int channel_id; /// < 人脸所属的通道 
        int frame_id; /// < 人脸所属的视频帧 
        int face_id; /// < 人脸ID 

        int left;
        int top;
        int right;
        int bottom;

        int confidence; /// < 人脸置信度 百分制 [0, 100]
    } HSRTMP_FaceRect;

    typedef struct _HSRTMP_Image {
        const void *p_buf;
        int size;
        int width;
        int height;
        int format;
    } HSRTMP_Image;

    /// < parameters
    typedef struct _HSRTMP_Parameters {
        int id;
        int format; /// < only support BGR and I420 for now.
        int width;
        int height;
        int size;
        int fps;
        int compress_ratio;
        const char *p_url;
    } HSRTMP_Parameters;

    void* hsrtmp_create(const HSRTMP_Parameters*, int *p_error_code);
    void hsrtmp_destroy(void *p_rtmp);

    int hsrtmp_start(void *p_source);
    int hsrtmp_pause(void *p_source);
    int hsrtmp_stop(void *p_source);
    void hsrtmp_push(void *p_source, void *data);
    void hsrtmp_push_and_draw_bbox(void *p_source, HSRTMP_Image *img, const HSRTMP_FaceRect *p_rects, const int rect_count);



#ifdef __cplusplus
}
#endif

#endif
