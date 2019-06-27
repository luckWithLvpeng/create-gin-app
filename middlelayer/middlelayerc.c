#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <pthread.h>
#include "middlelayerc.h"

#define CHANNEL_COUNT 32
int count = 0;
pthread_mutex_t videowriter_mutex[CHANNEL_COUNT];
pthread_mutex_t mutex[CHANNEL_COUNT];
pthread_mutex_t send_mutex[CHANNEL_COUNT];
pthread_mutex_t webrtc_mutex[CHANNEL_COUNT];
int exist_channels[CHANNEL_COUNT];
//int exist_channels_size = 0;
int use_rtmp = 0;
int use_webrtc = 0;
int mjpeg_reize = 3;
int mjpeg_quality = 70;
/***************************************************************************/

//BEGIN...SERVER_VERSION - FOR RESTFUL API - Added by Jiefeng Lai
void IdentifyFromServer(int subid, char *pImage, int imglen, HSFE_PersonCandidate *pCands, int *pCandsCount)
{
	IdentifyFromServerWithGo(subid, pImage, imglen, pCands, pCandsCount);
}
void ExtractFeatureFromServer(HSFE_Buffer *pImageBuffer, int num, HSFE_Buffer *pFeatureBuffer)
{
	ExtractFeatureFromServerWithGo(pImageBuffer, num, pFeatureBuffer);
}
void SetCallBackIdentifyFromServer()
{
	hsfe_set_identify_from_server_callback(IdentifyFromServer);
}
void SetCallBackFeatureExtraction()
{
	hsfe_set_feature_extraction_from_server_callback(ExtractFeatureFromServer);
}
//END//////////////////////////////////////////////////////////////////////////////////////////////////////

void MonitorResultArrive(const HSFE_MonitorResult* MonitorResult)
{
    MonitorResultFromC result;
    result.HitFlag = MonitorResult->hit_flag;
    result.FaceId = MonitorResult->face_id;
    result.ChannelId = MonitorResult->channel_id;
    result.FeatureSize = MonitorResult->feature_size;
    result.CropFrameFeature = MonitorResult->p_feature;

    result.Score = MonitorResult->p_result_unit[0].score;
    result.FeatureId = MonitorResult->p_result_unit[0].template_id;
    result.SublibId = MonitorResult->p_result_unit[0].sublib_id;
    result.FrameId = MonitorResult->p_result_unit[0].scene_id;
    result.FaceLeft = MonitorResult->p_result_unit[0].face_rect.left;
    result.FaceTop = MonitorResult->p_result_unit[0].face_rect.top;
    result.FaceRight = MonitorResult->p_result_unit[0].face_rect.right;
    result.FaceBottom = MonitorResult->p_result_unit[0].face_rect.bottom;
    result.CropFrameWidth = MonitorResult->p_result_unit[0].face_width;
    result.CropFrameHeight = MonitorResult->p_result_unit[0].face_height;
    result.FrameWidth = MonitorResult->p_result_unit[0].scene_width;
    result.FrameHeight = MonitorResult->p_result_unit[0].scene_height;
    result.FrameFormat = MonitorResult->p_result_unit[0].scene_format;
    result.Frame = (char*)MonitorResult->p_result_unit[0].p_scene;
    // Compress snap face image to JPG
    HSFE_Image srcImage;
    srcImage.p_buf = MonitorResult->p_result_unit[0].p_face;
    srcImage.width = result.CropFrameWidth;
    srcImage.height = result.CropFrameHeight;
    srcImage.format = result.FrameFormat;

    HSFE_Rect rect;
    rect.left = 0;
    rect.top = 0;
    rect.right = srcImage.width;
    rect.bottom = srcImage.height;
    int rect_width = rect.right - rect.left;
    int rect_height = rect.bottom - rect.top;

    rect.left -= rect_width / 3;
    rect.right += rect_width / 3;
    rect.top -= rect_height / 3;
    rect.bottom += rect_height / 3;
    if (rect.left < 0) rect.left = 0;
    if (rect.top < 0) rect.top = 0;
    if (rect.right >= srcImage.width) rect.right = srcImage.width - 1;
    if (rect.bottom >= srcImage.height) rect.bottom = srcImage.height - 1;

    HSFE_Buffer jpgBuf;
    result.CropRet = MyCreateJPG(&srcImage, 100, 95, &jpgBuf, &rect);
    result.CropFrame = jpgBuf.p_buf;
    result.CropFrameSize = jpgBuf.size;

    int other_candidates_count = MonitorResult->unit_count;
    result.OtherCandidatesCount = other_candidates_count;
    for (int i = 0; i < other_candidates_count; i++) {
        result.OtherFeatureIds[i] = MonitorResult->p_result_unit[i].template_id;
        result.OtherSublibIds[i] = MonitorResult->p_result_unit[i].sublib_id;
        result.OtherScores[i] = MonitorResult->p_result_unit[i].score;
        result.OtherHitFlags[i] = MonitorResult->p_result_unit[i].hit_flag;
    }

    // 保存剪裁图像
    //FILE *pFile = fopen("crop.rgb", "w"); 
    //fwrite (result.Frame, 1, result.FrameSize, pFile);
    //fflush(pFile);
    //fclose(pFile);

    MonitorCallbackFromGo(&result);
    if (jpgBuf.p_buf != NULL && result.CropRet == HSFE_EC_NO_ERROR) {
        MyFreeJPGBuf(jpgBuf.p_buf);
    }
}
/****************************************************************************/
void SetCallBackMonitorResult()
{
    hsfe_set_monitor_callback(MonitorResultArrive);
}
/****************************************************************************/
void TrackResultArrive(const HSFE_TrackResult* TrackResult)
{
    HSFE_Image colorImage;
    colorImage.format = TrackResult->scene_format;
    colorImage.width = TrackResult->scene_width;
    colorImage.height = TrackResult->scene_height;
    colorImage.p_buf = TrackResult->p_scene;

    if (pthread_mutex_trylock(&videowriter_mutex[TrackResult->channel_id]) == 0) {
        push_videowriter_and_draw_bbox(TrackResult->channel_id, &colorImage, TrackResult->p_face_rect, TrackResult->face_count);
        pthread_mutex_unlock(&videowriter_mutex[TrackResult->channel_id]);
    }
    if (use_webrtc == 1) {
        push_webrtc_and_draw_bbox(TrackResult->channel_id, &colorImage, TrackResult->p_face_rect, TrackResult->face_count);
    }
    if (use_rtmp == 1) {
        push_rtmp_and_draw_bbox(TrackResult->channel_id, &colorImage, TrackResult->p_face_rect, TrackResult->face_count);
    } else {
        int exist = 0;
        if (pthread_mutex_trylock(&mutex[TrackResult->channel_id]) == 0) {
            if (exist_channels[TrackResult->channel_id] > 0) {
                exist = 1;
            }
            pthread_mutex_unlock(&mutex[TrackResult->channel_id]);
        }
        if(exist == 1) {
            if (pthread_mutex_trylock(&send_mutex[TrackResult->channel_id]) == 0) {
                HSFE_Buffer buf;
                int ret = MyCreateJPGAndDrawLine(&colorImage, mjpeg_reize, mjpeg_quality, &buf,
                        TrackResult->p_face_rect, TrackResult->face_count);
                TrackResultArriveFromGo(buf, TrackResult->channel_id);
                MyFreeJPGBuf(buf.p_buf);
            }
        }
    }
}
/****************************************************************************/
void SetCallBackTrackResult()
{
    hsfe_set_track_callback(TrackResultArrive);
}
/****************************************************************************/
void WebrtcSampleArrive(const GST_WEBRTC_Sample *sample)
{
    // call WebrtcSampleArriveFromGo
    // printf("webrtc sample --------\n");
    HandlePipelineBufferFromGo(sample->p_buf, sample->size, sample->duration, sample->id);
    // char *p = (char *)sample->p_buf;
    // for (int i = 0; i < 100 && i < sample->size; i++) {
    //     printf("%x, ", p[i]);
    // }
    // printf("\n --------\n");
    // printf("%d -> %d --- duration -> %d ****************************\n", sample->id, sample->size, sample->duration);
}
/****************************************************************************/
int MyEngineCreate(HSFE_EngineParameters *p_engine_parameters)
{
    int ret = 0;
    for (int i=0; i<CHANNEL_COUNT; i++) {
        exist_channels[i] = 0;
        pthread_mutex_init(&mutex[i], NULL);
        pthread_mutex_init(&send_mutex[i], NULL);
        pthread_mutex_init(&videowriter_mutex[i], NULL);
        pthread_mutex_init(&webrtc_mutex[i], NULL);
    }
    ret = create(p_engine_parameters);
    return ret;
}
/****************************************************************************/
int MyUpdateMjpegSetting(int resize, int quality)
{
    mjpeg_reize = resize;
    mjpeg_quality = quality;
    return HSFE_EC_NO_ERROR;
}
/****************************************************************************/
int MySetFeatureMaxCount(int max_count)
{
    int ret = 0;
    ret = hsfe_set_template_max_count(max_count);
    return ret;
}
/****************************************************************************/
void MyEngineDestory()
{
    for (int i=0; i<CHANNEL_COUNT; i++) {
        pthread_mutex_unlock(&mutex[i]);
        pthread_mutex_unlock(&send_mutex[i]);
        pthread_mutex_unlock(&videowriter_mutex[i]);
        pthread_mutex_unlock(&webrtc_mutex[i]);
        pthread_mutex_destroy(&mutex[i]);
        pthread_mutex_destroy(&send_mutex[i]);
        pthread_mutex_destroy(&videowriter_mutex[i]);
        pthread_mutex_destroy(&webrtc_mutex[i]);
    }
    destroy();
}
/****************************************************************************/
int MyAddChannel(HSVD_DecoderParameters* DecoderParameters, HSFE_ChannelParameters* ChannelParameters)
{
    int ret = 0;
    if (use_rtmp) {
        add_rtmp_channel(DecoderParameters);
    }
    ret = add_channel(DecoderParameters, ChannelParameters);
    return ret;
}
/****************************************************************************/
void MyDelChannel(int channel_id)
{
    hsfe_clear_face_record(channel_id);
    remove_channel(channel_id);
    if (use_rtmp) {
        remove_rtmp_channel(channel_id);
    }
}
/****************************************************************************/
int MyAddSublib(int sublib_id)
{
    int ret = 0;
    ret = hsfe_add_sublib(sublib_id);
    return ret;
}
/****************************************************************************/
int MyDelSublib(int sublib_id)
{
    int ret = 0;
    ret = hsfe_remove_sublib(sublib_id);
    return ret;
}
/****************************************************************************/
int MySetSublibQuality(int sublib_id, int quality_score)
{
       int ret = 0;
       ret = hsfe_set_sublib_quality(sublib_id, quality_score);
       return ret;
}
/****************************************************************************/
int MyAddFeature(HSFE_FaceTemplate template_batch[], int batch_size, int sublib_id)
{
    int ret = 0;
    ret = hsfe_add_template(template_batch, batch_size, sublib_id);
    return ret;
}
/****************************************************************************/
int MyAddPhoto(char *path, HSFE_FaceTemplate *temp, int sublib_id)
{
    int ret = 0;
#ifdef USE_HOBOT
    ret = hsfe_add_photo(path, temp, sublib_id);
#endif
    return ret;
}
/****************************************************************************/
int MyDelFeature(int template_id, int sublib_id)
{
    int ret = 0;
    ret = hsfe_remove_template(template_id, sublib_id);
    return ret;
}
/****************************************************************************/
int MyDelAllFeature()
{
    int ret = 0;
    ret = hsfe_remove_all_templates();
    return ret;
}
/****************************************************************************/
int MyExtractFeature(HSFE_ExtractFeaturePacket pkt_batch[], int batch_size)
{
    int ret = 0;
    ret = hsfe_extract_face_feature(pkt_batch, batch_size);
    return ret;
}
int MyExtractImageFeature(HSFE_ImageFeaturePacket *p_ifpkt)
{
    int ret = 0;
    ret = hsfe_extract_feature(p_ifpkt);
    return ret;
}
/****************************************************************************/
int MyGetCurrentFeatureCount()
{
    int ret = 0;
    int currentCount = 0;
    ret = hsfe_get_template_current_count(&currentCount);
    //printf("ret = %d, currentCount = %d\n", ret, currentCount);
    return currentCount;
}
/****************************************************************************/
int MyCreateJPG(HSFE_Image *p_color_img, int reduce_ratio_percent, int jpg_quality, HSFE_Buffer *p_jpg_buf, HSFE_Rect *p_roi)
{
    int ret = 0;
    ret = hsfe_create_jpg(p_color_img, reduce_ratio_percent, jpg_quality, p_jpg_buf, p_roi);
    return ret;
}
/****************************************************************************/
int MyJPGLoadColor(char *file_path, void **pp_color_img, int *p_width, int *p_height)
{
    int ret = 0;
    ret = hsfe_load_color(file_path, pp_color_img, p_width, p_height);

    return ret;

}

/****************************************************************************/
int MyCropImage(HSFE_Image *p_img, HSFE_Rect *p_rect, void *p_croped_img)
{
    int ret = 0;
    ret = hsfe_crop_img(p_img, p_rect, p_croped_img);
    return ret;
}
/****************************************************************************/
int MyGetDecoderErrorCode()
{
    int ret = 0;
    ret = get_decoder_error_code();
    return ret;
}
/****************************************************************************/
int MyGetEngineErrorCode()
{
    int ret = 0;
    ret = get_engine_error_code();
    return ret;
}
/****************************************************************************/
int MyGetSDKErrorCode()
{
    int ret = 0;
    ret = hsfe_get_sdk_last_error();
    return ret;
}
/****************************************************************************/
int MyCreateJPGAndDrawLine(HSFE_Image *p_color_img,
                            int reduce_ratio,
                            int jpg_quality,
                            HSFE_Buffer *p_jpg_buf,
                            HSFE_FaceRect *p_rects,
                            int rect_count)
{
    int ret = 0;
    ret = hsfe_create_jpg_and_draw_line(p_color_img, 
            reduce_ratio, 
            jpg_quality, 
            p_jpg_buf, 
            p_rects, 
            rect_count);
    return ret;
}
/****************************************************************************/
int MyFreeJPGBuf(void *p_buf)
{
    int ret = 0;
    ret = hsfe_free_jpg_buf(p_buf);
    return ret;
}
/****************************************************************************/
int MyVerify(unsigned char* feature1, unsigned char* feature2, float* score) {
    return hsfe_verify(feature1, feature2, score);
}
/****************************************************************************/
int MySearchOneTemplate(unsigned char *p_feature, int feature_size, int face_width, float confidence, int sublib_count, int *sublib_tosearch, int *result_count, int *scores, int *featureids, int *sublibs) {
    return hsfe_search_one_template(p_feature, feature_size, face_width, confidence, sublib_count, sublib_tosearch, result_count, scores, featureids, sublibs);
}
/****************************************************************************/
int MySearchOnePhoto(char *path, int sublib_count, int *sublib_tosearch, int *result_count, int *scores, int *featureids, int *sublibs) {
#ifdef USE_HOBOT
    return hsfe_search_one_photo(path, sublib_count, sublib_tosearch, result_count, scores, featureids, sublibs);
#else
    return HSFE_EC_NO_ERROR;
#endif
}
/****************************************************************************/
int MyUnlockSendMutex(int channel_id)
{
    pthread_mutex_unlock(&send_mutex[channel_id]);
}
/****************************************************************************/
int MySetMJPGPushChannelCount(int channel_id, int value)
{
    int ret = 0;
    pthread_mutex_lock(&mutex[channel_id]);
    if (channel_id < CHANNEL_COUNT)
        exist_channels[channel_id] = value;
    pthread_mutex_unlock(&mutex[channel_id]);
    return ret;
}
/****************************************************************************/
int * MyMallocIntArry(int size)
{
    int * des = (int*)malloc(size * sizeof(int));
    memset(des, 0, size * sizeof(int));
    return des;
}
/****************************************************************************/
void MyMemcpyIntArry(int * des, int * src, int size)
{
    for(int i = 0; i < size; ++i)
    {
        des[i] = src[i];
    }
}
/****************************************************************************/
char * MyMallocCharArry(int size)
{
    char * des = (char*)malloc(size + 1);
    memset(des, '\0', size + 1);
    return des;
}
/****************************************************************************/
void MyMemcpyCharArry(char * des, char * src, int size)
{
    memcpy(des, src, size);
}
/****************************************************************************/
void MyFree(void * buf)
{
    free(buf);
}
/****************************************************************************/
void MySetUseWebrtcFlag(int val)
{
    use_webrtc = val;
}
/****************************************************************************/
void MySetUseRtmpFlag(int val)
{
    use_rtmp = val;
}
/****************************************************************************/
int MyAddVideowriterChannel(GST_VIDEO_WRITER_Parameters *p_params)
{
    int ret = 0;
    pthread_mutex_lock(&videowriter_mutex[p_params->id]);
    ret = add_videowriter_channel(p_params);
    pthread_mutex_unlock(&videowriter_mutex[p_params->id]);
    return ret;
}
void MyDelVideowriterChannel(int id)
{
    pthread_mutex_lock(&videowriter_mutex[id]);
    remove_videowriter_channel(id);
    pthread_mutex_unlock(&videowriter_mutex[id]);
}
int MyAddWebrtcChannel(GST_WEBRTC_Parameters *p_params, const char *p_url, int use_rtmp)
{
    int ret = 0;
    pthread_mutex_lock(&webrtc_mutex[p_params->id]);
    ret = add_webrtc_channel(p_params, WebrtcSampleArrive, p_url, use_rtmp);
    pthread_mutex_unlock(&webrtc_mutex[p_params->id]);
    return ret;
}
void MyDelWebrtcChannel(int id)
{
    pthread_mutex_lock(&webrtc_mutex[id]);
    remove_webrtc_channel(id);
    pthread_mutex_unlock(&webrtc_mutex[id]);
}
/****************************************************************************/
void testIntArry(int *list, int size)
{
    for(int i = 0; i < size; ++i)
    {
        //printf("%d\n",list[i]);
    }
}
/****************************************************************************/
void testCharArry(char *list)
{
    //printf("%s\n",list);
}
/****************************************************************************/
