/*********************************************************** 
* Date: 2019-04-19
* 
* Author: Wang Xinyu
* 
* Email: wangxinyu@alleyes.com.cn
*
* Department: 
* 
* Company: 
* 
* Module: gst-videowriter
* 
* Brief: push stream from appsrc to filesink
*
* Note: 
* 
* CodePage: Pure UTF-8 
************************************************************/ 
#ifndef __GST_VIDEO_WRITER_INTERFACE_H__
#define __GST_VIDEO_WRITER_INTERFACE_H__


#ifdef __cplusplus
extern "C" {
#endif
    typedef enum _GST_VIDEO_WRITER_ErrorCode {
        GST_VIDEO_WRITER_EC_UNKNOWN_ERROR = (-1),
        GST_VIDEO_WRITER_EC_NO_ERROR = 0,
        GST_VIDEO_WRITER_EC_INVALID_PARAMETER = 1,
        GST_VIDEO_WRITER_EC_EXCEPTION = 2
    } GST_VIDEO_WRITER_ErrorCode;

    typedef enum _GST_VIDEO_WRITER_FrameFormat {
        GST_VIDEO_WRITER_FRAME_FORMAT_INVALID = 0,
        GST_VIDEO_WRITER_FRAME_FORMAT_BGR = 1,
        GST_VIDEO_WRITER_FRAME_FORMAT_NV12 = 2,
        GST_VIDEO_WRITER_FRAME_FORMAT_I420 = 3
    } GST_VIDEO_WRITER_FrameFormat;

    typedef struct _GST_VIDEO_WRITER_FaceRect {
        int channel_id;
        int frame_id;
        int face_id;

        int left;
        int top;
        int right;
        int bottom;

        int confidence;
    } GST_VIDEO_WRITER_FaceRect;

    typedef struct _GST_VIDEO_WRITER_Image {
        const void *p_buf;
        int size;
        int width;
        int height;
        int format;
    } GST_VIDEO_WRITER_Image;

    /// < parameters
    typedef struct _GST_VIDEO_WRITER_Parameters {
        int id;
        int format; /// < only support BGR and I420 for now.
        int width;
        int height;
        int size;
        int fps;
        int compress_ratio;
        const char *p_video_path;
    } GST_VIDEO_WRITER_Parameters;

    void* gst_videowriter_create(const GST_VIDEO_WRITER_Parameters*, int *p_error_code);
    void gst_videowriter_destroy(void *p_videowriter);

    int gst_videowriter_start(void *p_source);
    int gst_videowriter_pause(void *p_source);
    int gst_videowriter_stop(void *p_source);
    void gst_videowriter_push(void *p_source, void *data);
    void gst_videowriter_push_and_draw_bbox(void *p_source, GST_VIDEO_WRITER_Image *img, const GST_VIDEO_WRITER_FaceRect *p_rects, const int rect_count);

#ifdef __cplusplus
}
#endif

#endif  // __GST_VIDEO_WRITER_INTERFACE_H__
