#ifndef _MIDDLELAYERC_H_
#define _MIDDLELAYERC_H_

#include "../hswrapper/hsfe_wrapper.h"

/***************************************************************************/
// 返回结果结构体
#define OTHER_CANDIDATES_MAX 32
typedef struct _MonitorResult
{
    int HitFlag;
    int FaceId;
    int ChannelId;
    int Score;
    int FeatureId;
    int SublibId;
    int FrameId;
    int FaceLeft;
    int FaceTop;
    int FaceRight;
    int FaceBottom;
    int FrameSize;
    char *Frame;
    int FrameWidth;
    int FrameHeight;
    int FrameFormat;
    int CropFrameSize;
    void *CropFrame;
    int CropFrameWidth;
    int CropFrameHeight;
    void *CropFrameFeature;
    int FeatureSize;
    int CropRet;
    int OtherCandidatesCount;
    int OtherFeatureIds[OTHER_CANDIDATES_MAX];
    int OtherSublibIds[OTHER_CANDIDATES_MAX];
    int OtherScores[OTHER_CANDIDATES_MAX];
	int OtherHitFlags[OTHER_CANDIDATES_MAX];
} MonitorResultFromC;
/****************************************************************************/
int MyUpdateMjpegSetting(int resize, int quality);
/****************************************************************************/
/****************************************************************************/

//BEGIN...SERVER_VERSION - FOR RESTFUL API - Added by Jiefeng Lai
void IdentifyFromServer(int subid, char *pImage, int imglen, HSFE_PersonCandidate *pCands, int *pCandsCount);
void ExtractFeatureFromServer(HSFE_Buffer *pImageBuffer, int num, HSFE_Buffer *pFeatureBuffer);

void SetCallBackIdentifyFromServer();
void SetCallBackFeatureExtraction();
//END////////////////////////////////////////////////////////////////////////////////////////////////////


/*监控结果返回回调函数*/
void MonitorResultArrive(const HSFE_MonitorResult* MonitorResult);

/*设置监控结果回调函数*/
void SetCallBackMonitorResult();

/*跟踪回调函数*/
void TrackResultArrive(const HSFE_TrackResult* TrackResult);

/*设置跟踪回调函数*/
void SetCallBackTrackResult();

/*引擎开启*/
int MyEngineCreate(HSFE_EngineParameters *p_engine_parameters);

/*引擎开启*/
int MySetFeatureMaxCount(int max_count);

/*引擎销毁*/
void MyEngineDestory();

/*添加通道*/
int MyAddChannel(HSVD_DecoderParameters* DecoderParameters, HSFE_ChannelParameters* ChannelParameters);

/*删除通道*/
void MyDelChannel(int channel_id);

/*添加分库*/
int MyAddSublib(int sublib_id);

/*删除分库*/
int MyDelSublib(int sublib_id);

/*设置分库质量*/
int MySetSublibQuality(int sublib_id, int quality_score);

/*添加特征*/
int MyAddFeature(HSFE_FaceTemplate template_batch[], int batch_size, int sublib_id);

int MyAddPhoto(char *path, HSFE_FaceTemplate *temp, int sublib_id);

/*删除特征*/
int MyDelFeature(int template_id, int sublib_id);

/*删除所有特征*/
int MyDelAllFeature();

/*特征提取*/
int MyExtractFeature(HSFE_ExtractFeaturePacket pkt_batch[], int batch_size);
int MyExtractImageFeature(HSFE_ImageFeaturePacket *p_ifpkt);

/*获取当前引擎中的特征总数*/
int MyGetCurrentFeatureCount();

/*剪裁图像*/
int MyCropImage(HSFE_Image *p_img, HSFE_Rect *p_rect, void *p_croped_img);

int MyGetDecoderErrorCode();

int MyGetEngineErrorCode();

int MyGetSDKErrorCode();

int MyCreateJPG(HSFE_Image *p_color_img, int reduce_ratio_percent, int jpg_quality, HSFE_Buffer *p_jpg_buf, HSFE_Rect *p_roi);

int MyJPGLoadColor(char *file_path, void **pp_color_img, int *p_width, int *p_height);

int MyCreateJPGAndDrawLine(HSFE_Image *p_color_img,
                            int reduce_ratio, 
                            int jpg_quality,
                            HSFE_Buffer *p_jpg_buf,
                            HSFE_FaceRect *p_rects,
                            int rect_count);

int MyFreeJPGBuf(void *p_buf);

int MyVerify(unsigned char* feature1, unsigned char* feature2, float* score);

int MySearchOneTemplate(unsigned char *p_feature, int feature_size, int face_width, float confidence, int sublib_count, int *sublib_tosearch, int *result_count, int *scores, int *featureids, int *sublibs);

int MySearchOnePhoto(char *path, int sublib_count, int *sublib_tosearch, int *result_count, int *scores, int *featureids, int *sublibs);

/****************************************************************************/
int MyUnlockSendMutex(int channel_id);
int MySetMJPGPushChannelCount(int channel_id, int value);
/****************************************************************************/
/*申请空间*/
int * MyMallocIntArry(int size);
void MyMemcpyIntArry(int * des, int * src, int size);

char * MyMallocCharArry(int size);
void MyMemcpyCharArry(char * des, char * src, int size);

void MyFree(void * buf);
void MySetUseRtmpFlag(int val);
void MySetUseWebrtcFlag(int val);
int MyAddVideowriterChannel(GST_VIDEO_WRITER_Parameters *p_params);
void MyDelVideowriterChannel(int id);
int MyAddWebrtcChannel(GST_WEBRTC_Parameters *p_params, const char *p_url, int use_rtmp);  // p_url -> url of rtmp, use_rtmp -> 1: use rtmp, 0: return encoded buffer
void MyDelWebrtcChannel(int id);
/****************************************************************************/

/****************************************************************************/
/*测试输出函数*/
void testIntArry(int *list,int size);
void testCharArry(char *list);
/****************************************************************************/

#endif
