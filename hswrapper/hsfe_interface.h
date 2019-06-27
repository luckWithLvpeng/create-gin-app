/*********************************************************** 
* Date: 2016-07-04 
* 
* Author: 牟韵 
* 
* Department: 智能产品中心-研发部 
* 
* Company: 北京海鑫科金高科技股份有限公司 
* 
* Module: 海鑫人脸引擎接口 
* 
* Brief: 
* 
* Note:  TX1版本 
* 
* CodePage: Pure UTF-8 
************************************************************/ 
#ifndef __HSFE_INTERFACE_H_BY_MOUYUN_2016_07_04__
#define __HSFE_INTERFACE_H_BY_MOUYUN_2016_07_04__

#include "hsfe_types.h"

#define HSFE_DLL_INTERFACE
#define HSFE_STDCALL

#ifdef __cplusplus
extern "C" {
#endif

    ////////////////////////////////////////////////回调//////////////////////////////////////////////////////////////
    typedef void(HSFE_STDCALL *PFUNC_HSFE_PROCESS_CONSUMED_TIME_ARRIVE)(const HSFE_ConsumedTimeResult*);
    typedef void(HSFE_STDCALL *PFUNC_HSFE_TRACK_RESULT_ARRIVE)(const HSFE_TrackResult*);
    typedef void(HSFE_STDCALL *PFUNC_HSFE_VERIFY_RESULT_ARRIVE)(const HSFE_VerifyResult*);
    typedef void(HSFE_STDCALL *PFUNC_HSFE_MONITOR_RESULT_ARRIVE)(const HSFE_MonitorResult*);
    typedef void(HSFE_STDCALL *PFUNC_HSFE_IDENTIFY_FROM_SERVER)(int subid, char *pImage, int imglen, HSFE_PersonCandidate *pCands, int *pCandsCount);
	typedef void(HSFE_STDCALL *PFUNC_HSFE_FEATURE_EXTRACTION_FROM_SERVER)(HSFE_Buffer *pImageBuffer, int num, HSFE_Buffer *pFeatureBuffer);
	typedef void(HSFE_STDCALL *PFUNC_HSFE_CHANNEL_STOPPED)(int channel_id);
	
    HSFE_DLL_INTERFACE void hsfe_set_process_consumed_time_callback(PFUNC_HSFE_PROCESS_CONSUMED_TIME_ARRIVE cb);
    HSFE_DLL_INTERFACE void hsfe_set_track_callback(PFUNC_HSFE_TRACK_RESULT_ARRIVE cb);
    HSFE_DLL_INTERFACE void hsfe_set_verify_callback(PFUNC_HSFE_VERIFY_RESULT_ARRIVE cb);
    HSFE_DLL_INTERFACE void hsfe_set_monitor_callback(PFUNC_HSFE_MONITOR_RESULT_ARRIVE cb);
    HSFE_DLL_INTERFACE void hsfe_set_identify_from_server_callback(PFUNC_HSFE_IDENTIFY_FROM_SERVER cb);
	HSFE_DLL_INTERFACE void hsfe_set_feature_extraction_from_server_callback(PFUNC_HSFE_FEATURE_EXTRACTION_FROM_SERVER cb);
	HSFE_DLL_INTERFACE void hsfe_set_channel_stopped_callback(PFUNC_HSFE_CHANNEL_STOPPED cb);
    ////////////////////////////////////////////////回调//////////////////////////////////////////////////////////////


    ////////////////////////////////////////////////用户数据//////////////////////////////////////////////////////////
    HSFE_DLL_INTERFACE void hsfe_set_user_data(void *p_user_data);
    HSFE_DLL_INTERFACE void* hsfe_get_user_data();
    ////////////////////////////////////////////////用户数据//////////////////////////////////////////////////////////

    HSFE_DLL_INTERFACE void hsfe_clear_face_record(const int channel_id); /// < 清除一个通道的人脸记录缓存 

    HSFE_DLL_INTERFACE int hsfe_get_sdk_last_error(); /// < 获取算法SDK的最后一次错误码 当引擎接口返回HSFE_EC_SDK_ERROR时，可以通过此接口可获取算法SDK的错误码 


    HSFE_DLL_INTERFACE int hsfe_get_sdk_version(); /// 获取算法SDK版本
    HSFE_DLL_INTERFACE int hsfe_get_feature_size(); /// 获取算法特征大小
    HSFE_DLL_INTERFACE int hsfe_create(const HSFE_EngineParameters *p_engine_parameters); /// < 创建引擎 不支持多线程 
    HSFE_DLL_INTERFACE void hsfe_destroy();/// < 销毁引擎 不支持多线程 

    HSFE_DLL_INTERFACE int hsfe_add_channel(const HSFE_ChannelParameters *p_channel_parametes); /// < 添加通道 支持多线程 
    HSFE_DLL_INTERFACE int hsfe_remove_channel(const int channel_id); /// < 移除通道 支持多线程 

    HSFE_DLL_INTERFACE int hsfe_input_frame(const HSFE_Frame *p_frame); /// < 输入视频帧 支持多线程 

    HSFE_DLL_INTERFACE int hsfe_input_verify_data(const HSFE_VerifyInputData *p_verify_data); /// < 输入核验数据 支持多线程 


    HSFE_DLL_INTERFACE int hsfe_add_sublib(const int sublib_id); /// < 添加分库 调用者须保证分库id的全局唯一性 支持多线程 
    HSFE_DLL_INTERFACE int hsfe_set_sublib_quality(const int sublib_id, const int quality_score); /// < 设置分库质量分数score = score *( quality_score(1-1000) ) /800 支持多线程 
    HSFE_DLL_INTERFACE int hsfe_remove_sublib(const int sublib_id); /// < 移除分库 属于分库的模板也将被一起移除 支持多线程 
    HSFE_DLL_INTERFACE int hsfe_remove_all_sublibs(); /// < 移除所有分库 所有模板也将被一起移除 支持多线程 

    HSFE_DLL_INTERFACE int hsfe_add_template(HSFE_FaceTemplate template_batch[], const int batch_size, const int sublib_id); /// < 添加模板 支持多线程 调用者须保证指定的分库存在以及特征数据与具体的算法版本相匹配 
    HSFE_DLL_INTERFACE int hsfe_add_photo(const char *url, HSFE_FaceTemplate *temp, const int sublib_id);
    HSFE_DLL_INTERFACE int hsfe_remove_template(const int template_id, const int sublib_id); /// < 移除模板 支持多线程 
    HSFE_DLL_INTERFACE int hsfe_remove_all_templates(); /// < 移除所有模板 支持多线程 

    HSFE_DLL_INTERFACE int hsfe_set_template_max_count(const int max_count);
    HSFE_DLL_INTERFACE int hsfe_get_template_current_count(int *p_currrent_count);

    HSFE_DLL_INTERFACE int hsfe_load_color(const char *p_file_name, void **pp_color_img, int *p_width, int *p_height);
    HSFE_DLL_INTERFACE int hsfe_color_to_gray(const unsigned char *p_color, const int _width, const int _height, void *p_gray_img, int color_channel_sequence);

    HSFE_DLL_INTERFACE int hsfe_extract_face_feature(HSFE_ExtractFeaturePacket pkts[], const int pkt_count); /// < 批量提取人脸特征 不支持多线程 
    HSFE_DLL_INTERFACE int hsfe_extract_feature(HSFE_ImageFeaturePacket *p_ifpkt); /// < 批量提取人脸特征 不支持多线程 

    HSFE_DLL_INTERFACE int hsfe_crop_img(const HSFE_Image *p_img, HSFE_Rect *p_rect, void *p_croped_img); /// < 裁剪图片 不支持多线程 

    HSFE_DLL_INTERFACE int hsfe_create_jpg_and_draw_line(const HSFE_Image *p_color_img, const int reduce_ratio, const int jpg_quality, HSFE_Buffer *p_jpg_buf, const HSFE_FaceRect *p_rects, const int rect_count);
    HSFE_DLL_INTERFACE int hsfe_create_jpg(const HSFE_Image *p_color_img, const int reduce_ratio_percent, const int jpg_quality, HSFE_Buffer *p_jpg_buf, const HSFE_Rect *p_roi);

    HSFE_DLL_INTERFACE int hsfe_free_jpg_buf(void *p_buf);

    HSFE_DLL_INTERFACE int hsfe_verify(const unsigned char* feature1, const unsigned char* feature2, float* score);
    HSFE_DLL_INTERFACE int hsfe_photo_verify(const char* url1, const char* url2, float* score);
    HSFE_DLL_INTERFACE int hsfe_search_one_template(unsigned char *p_feature, const int feature_size, const int face_width, const float confidence, const int sublib_count, int *sublib_tosearch, int *result_count, int *scores, int *featureids, int *sublibs);
    HSFE_DLL_INTERFACE int hsfe_search_one_photo(const char *url, const int sublib_count, int *sublib_tosearch, int *result_count, int *scores, int *featureids, int *sublibs);
#ifdef __cplusplus
}
#endif

#endif
