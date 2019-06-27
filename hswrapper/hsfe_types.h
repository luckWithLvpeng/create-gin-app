/*********************************************************** 
* Date: 2016-00-00 
* 
* Author: 牟韵 
* 
* Department: 智能产品中心-研发部 
* 
* Company: 北京海鑫科金高科技股份有限公司 
* 
* Module: 海鑫人脸引擎-数据类型 
* 
* Brief: 
* 
* Note: 
* 
* CodePage: Pure UTF-8 
************************************************************/ 
#ifndef __HSFE_TYPES_H_BY_MOUYUN_2016_12_13__
#define __HSFE_TYPES_H_BY_MOUYUN_2016_12_13__

#ifdef __cplusplus
extern "C" {
#endif

    ////////////////////////////////////枚举//////////////////////////////////////
    /// < 错误码 
    typedef enum _HisignFaceEngineErrorCode
    {
        HSFE_EC_UNKNOWN_ERROR = (-1), /// < 未知错误 
        HSFE_EC_NO_ERROR = 0, /// < 正确 
        HSFE_EC_INVALID_PARAMETER = 1, /// < 无效参数 
        HSFE_EC_INVALID_MODE = 2, /// < 无效模式 
        HSFE_EC_INVALID_FEATURE_SIZE = 3, /// < 无效的特征大小 
        HSFE_EC_ENGINE_NOT_CREATED = 4, /// < 引擎未创建 
        HSFE_EC_ENGINE_ALREADY_CREATED = 5, /// < 引擎未创建 
        HSFE_EC_MODULE_NOT_INIT = 6, /// < 模块未初始化 
        HSFE_EC_MODULE_ALREADY_INITED = 7, /// < 模块已初始化 
        HSFE_EC_CREATE_LOGGER_FAILED = 8, /// < 创建日志模块失败 
        HSFE_EC_FILE_ERROR = 9, /// < 文件错误 
        HSFE_EC_MEMORY_ERROR = 10, /// < 内存错误 
        HSFE_EC_SDK_ERROR = 11, /// < 算法错误，具体的算法错误码请调用hsfe_get_sdk_last_error()获取 
        HSFE_EC_THREAD_ERROR = 12, /// < 线程错误 
        HSFE_EC_THREAD_QUEUE_FULL = 13, /// < 线程队列已满 
        HSFE_EC_THREAD_QUEUE_DISABLED = 14, /// < 线程队列已禁用 
        HSFE_EC_CHANNEL_NOT_EXIST = 15, /// < 无效的通道 
        HSFE_EC_CHANNEL_ALREADY_EXIST = 16, /// < 通道已存在 
        HSFE_EC_SUBLIB_NOT_EXIST = 17, /// < 分库不存在 
        HSFE_EC_TEMPLATE_NOT_EXIST = 18, /// < 模板不存在 
        HSFE_EC_TEMPLATE_POOL_IS_FULL = 19, /// < 模板池已满 
        HSFE_EC_NO_FEATURE_TO_VERIFY = 20, /// < 核验输入不存在有效特征 
        HSFE_EC_NO_FACE_IN_IMG = 21, /// < 图片中无人脸 
        HSFE_EC_EXCEED_MAX_CHANNEL_COUNT = 22,
        HSFE_EC_MULTI_FACES_IN_IMG = 23,
        HSFE_EC_FILE_TOO_SMALL = 24,
        HSFE_EC_LOW_QUALITY_IMG = 25,
        HSFE_EC_IMG_EXCEED_MAX_RESOLUTION = 26,
        HSFE_EC_FACE_POSE_ERROR = 27
    }HSFE_ErrorCode;

    /// < 通道模式 
    typedef enum _HSFE_ChannelMode
    {
        HSFE_CM_INVALID = 0, /// < 无效模式 
        HSFE_CM_VERIFY = 0x01, /// < 核验（1:1） 
        HSFE_CM_ONLY_MONITOR_MAX_FACE = 0x02, /// < 监控最大人脸（1:N） 
        HSFE_CM_MONITOR_ALL = 0x04, /// < 监控所有人脸（n:N），正常去重
        HSFE_CM_MONITOR_ALL_UNIQUE = 0x05, /// < 监控所有人脸（n:N），加后处理去重
        HSFE_CM_MONITOR_ALL_RAW = 0x06, /// < 监控所有人脸（n:N），发送所有跟踪人脸
        HSFE_CM_MONITOR_ALL_SNAP = 0x07 /// < 监控所有人脸（n:N），抓拍机模式
    }HSFE_ChannelMode;

    /// < 引擎业务流程代码 
    typedef enum _HSFE_ProcessID
    {
        HSFE_PID_INVALID = 0, /// < 无效 
        HSFE_PID_TRACK = 1, /// < 人脸跟踪 
        HSFE_PID_LOCATE = 2, /// < 人脸特征点定位 
        HSFE_PID_POSE = 2, /// < 人脸角度估计 
        HSFE_PID_EXTRACT = 3, /// < 人脸特征提取 
        HSFE_PID_COMPARE = 4 /// < 人脸特征比对 
    }HSFE_ProcessID;

    /// < 图像格式枚举 
    typedef enum _HSFE_ImageFormat
    {
        HSFE_IMG_FORMAT_INVALID = 0,
        HSFE_IMG_FORMAT_COLOR = 1,
        HSFE_IMG_FORMAT_GRAY = 2,
        HSFE_IMG_FORMAT_I420 = 3,
        HSFE_IMG_FORMAT_NV21 = 4
    }HSFE_ImageFormat;

    /// < 彩图的通道序列 
    typedef enum _HSFE_ColorChannelSequence
    {
        HSFE_CCS_INVALID = 0,
        HSFE_CCS_BGR = 1,
        HSFE_CCS_RGB = 2
    }HSFE_ColorChannelSequence;


    ////////////////////////////////////////////////////////结构体//////////////////////////////////////////////////////// 

    /// < 缓存  
    typedef struct _HSFE_Buffer
    {
        void *p_buf;
        int size;
    }HSFE_Buffer;

    /// < 矩形 
    typedef struct _HSFE_Rect
    {
        int left;
        int top;
        int right;
        int bottom;
    }HSFE_Rect;

    /// < 人脸矩形 
    typedef struct _HSFE_FaceRect
    {
        int channel_id; /// < 人脸所属的通道 
        int frame_id; /// < 人脸所属的视频帧 
        int face_id; /// < 人脸ID 

        int left;
        int top;
        int right;
        int bottom;

        int confidence; /// < 人脸置信度 百分制 [0, 100]
    }HSFE_FaceRect;

    /// < 灰度图像 
    typedef struct _HSFE_Image
    {
        const void *p_buf; /// < 图像数据 
        int size; /// < 图像大小
        int width; /// < 图像宽度 
        int height; /// < 图像高度 
        int format; /// < 图像格式 
    }HSFE_Image;

    /// < 视频帧 
    typedef struct _HSFE_Frame
    {
        int channel_id; /// < 通道ID 
        int frame_id; /// < 帧ID 
        const void *p_buf; /// < 图像数据  
    }HSFE_Frame;

    /// < 通道参数 
    typedef struct _HSFE_ChannelParameters
    {
        int channel_id; /// < 通道ID，非负整数 

        int reduce_ratio; /// < 缩图比例，建议值：1、2、4、8 

        int frame_width; /// < 图像宽度 
        int frame_height; /// < 图像高度 
        int frame_format; /// < 图像格式 参见枚举HSFE_ImageFormat 
        int color_channel_sequence; /// < 彩图的通道序列，当图像格式设为HSFE_IMG_FORMAT_COLOR时，需要设置此标志，参见枚举HSFE_ColorChannelSequence 
        int frame_keep_flag; /// < 是否保持视频帧 

        HSFE_Rect detect_roi; /// < 检测ROI区域
        int face_track_step; /// < 跟踪步长
        int max_track_count; /// 每个视频同时跟踪最大人脸数
        int face_filter_structure; /// < 人脸脸型过滤阈值
        int face_filter_clearity; /// < 人脸清晰度过滤值，小于该值时，权重为清晰度 / face_filter_clearity, 大于时为1
        int face_filter_size; /// < 人脸尺寸过滤大小，小于该值时，权重为face_wdith / face_filter_size, 大于时为1
        int pyramid_layer_min_height; /// < 金字塔阈值最小层间距
        int face_recorder_remove_dif; /// < face result reporter track_id移除间距
        int feature_recorder_remove_dif; /// < face result reporter feature移除间距

        int face_filte_count; /// < 人脸过滤计数，取值范围[1, N]，设为N时表示N个人脸中取1个（目前为第N个）送到后端处理，所以设为1时表示不过滤 

        int face_min_size; /// < 最小人脸 当前版本(V831)最小值为20 
        int face_max_size; /// < 最大人脸 不应超过图像宽度 

        int face_confidence; /// < 人脸置信度（百分制，建议设为75分左右） 
        int face_pose_estimate_flag; /// < 人脸角度估计的开关：0代表关，1代表开 
        int face_yaw_left; /// < 人脸水平左偏角 -90°至0° 
        int face_yaw_right; /// < 人脸水平右偏角 0°至90° 

        int prefeature_flag; /// < 预提特征开关：0代表关 1代表开 

        /// < 关于识别阈值和拒识阈值的设置，具体请知悉金字塔阈值模型 

        int recog_top_threshold; /// < 识别的顶级阈值 阈值范围[0, 1000] 一般应大于500 视算法版本与实际业务情况而定 
        int recog_bottom_threshold; /// < 识别的基准阈值 阈值范围[0, 1000] 必须小于recog_top_threshold 

        int reject_flag; /// < 拒识的开关 0代表关，1代表开 
        int reject_top_threshold; /// < 拒识的顶级阈值 阈值范围[0, 1000] 应小于recog_bottom_threshold 
        int reject_bottom_threshold; /// < 拒识的基准阈值 阈值范围[0, 1000] 必须小于reject_top_threshold 

        int channel_mode; /// < 通道模式，参见枚举HSFE_ChannelMode 

        int max_result_count; /// < （监控模式）比对结果保留数量 比对结果多于此数量时引擎会将差的结果丢弃 

        int merge_time_out_ms; /// < （监控模式）结果融合超时时间（毫秒） 人脸ID第一次到来时开始计时 超过此时间时报出结果 

        /// < 分库和分库的数量为置为0代表监控时比对所有分库 
        const int *p_sublib_for_monitor; /// < 监控的分库 
        int sublib_count; /// <  监控的分库的数量 
        int use_tracker; /// < whether init face_tracker(1) or not(0)
        int liveness_thresh;
    }HSFE_ChannelParameters;

    /// < 提取特征数据包 
    typedef struct _HSFE_ExtractFeaturePacket
    {
        const char *p_img_file_name;

        const void *p_gray;
        int gray_width;
        int gray_height;

        int error_code; /// < 错误码 参见枚举HSFE_ErrorCode 

        const void *p_feature;
        int feature_size;
        int quality_score;
    } HSFE_ExtractFeaturePacket;

    typedef struct _HSFE_ImageFeaturePacket
    {
        const char *p_img_file_name;
        int error_code;  // refer to HSFE_ErrorCode
        const void *p_feature;
        int feature_size;  // unit: byte
        int quality_score;  // reserved, return 800 for now.
        int file_size;  // image file size, unit: byte
        int img_width;
        int img_height;
        int is_gray;  // return 1 if the image is gray, return 0 if is colorful.
        int face_left;  // face bounding box, left
        int face_top;  // face bounding box, top
        int face_width;  // face bounding box, width
        int face_height;  // face bounding box, height
        int pitch;  // face pose, pitch angle, range from -90 to 90
        int yaw;  // face pose, yaw angle, range from -90 to 90
        int roll;  // face pose, roll angle, range from -90 to 90
        int ustddev;  // standard deviation value of U from YUV, it can be used to evaluate the quality of image.
    } HSFE_ImageFeaturePacket;

    /// < 核验的输入数据 
    typedef struct _HSFE_VerifyInputData
    {
        /// < 图像文件名 
        const char **pp_img_file_name;
        int img_file_count;

        /// < 图像数据 
        const HSFE_Image *p_img;
        int img_count;

        /// < 特征数据 
        const HSFE_Buffer *p_feature;
        int feature_count;

        int time_out_ms; /// < 超时时间，单位ms 

        int enable_reject; /// < 拒识开关 0:禁用 非0:开启 

        int top_threshold; /// < 顶级阈值 
        int bottom_threshold; /// < 基准阈值 

        /// < 若进行全通道比对，比对通道指针和其数量均设为0即可 
        int *p_channel_to_verify; /// < 比对通道 
        int channel_count; /// < 比对通道的数量 
    }HSFE_VerifyInputData;

    /// < 人脸模板（用于入库） 
    typedef struct _HSFE_FaceTemplate
    {
        int template_id;
        const void *p_feature;
        int feature_size;
        int result; /// < 入库结果 
    }HSFE_FaceTemplate;

    //////////////////////////////////回调函数及其结构体//////////////////////////////////////// 

    /// < 流程耗时 
    typedef struct _HSFE_ConsumedTimeResult
    {
        int process_id; /// < 流程ID 参见枚举HSFE_ProcessID 
        int consumed_time_us; /// < 流程耗时 
        void *p_user_data; /// < 用户数据 
    }HSFE_ConsumedTimeResult;

    /// < 跟踪结果 
    typedef struct _HSFE_TrackResult
    {
        int error_code; /// < 错误码 参见枚举HSFE_ErrorCode 

        int face_count; /// < 人脸数量 
        const HSFE_FaceRect *p_face_rect; /// < 人脸坐标数组（由face_count表示大小，已按人脸大小降序排列） 

        int channel_id;
        int scene_id;
        int scene_width;
        int scene_height;
        int scene_format;
        const void *p_scene; /// < 场景图  
        void *p_user_data; /// < 用户数据 
    }HSFE_TrackResult;

    /// < 核验结果 
    typedef struct _HSFE_VerifyResult
    {
        int consumed_time_ms; /// < 核验耗时(从核验数据输入至得到核验结果消耗的时间) 

        int hit_flag; /// < 比中标志 0:未比中 1:比中 
        int score; /// < 比分 

        int face_id; /// < 比中的人脸的ID 
        HSFE_Rect face_rect; /// < 比中的人脸的坐标 

        int channel_id;
        int scene_id;
        int scene_width;
        int scene_height;
        const void *p_scene; /// < 比中的人脸所属的场景图 

        void *p_user_data; /// < 用户数据 
    }HSFE_VerifyResult;

    /// < 监控结果单元  
    typedef struct _HSFE_MonitorResultUnit
    {
        int score; /// < 比分 
        int template_id; /// < 比中模板的ID 
        int sublib_id; /// < 比中的模板所属分库的ID 
        int hit_flag; /// < 比中标志 0:未比中 1:比中 

        HSFE_Rect face_rect; /// < 人脸坐标 

        int face_id;
        int face_width;
        int face_height;
        int face_format;
        const void *p_face;

        int scene_id;
        int scene_width;
        int scene_height;
        int scene_format;
        const void *p_scene; /// < 场景图 
    }HSFE_MonitorResultUnit;

    /// < 监控结果 
    typedef struct _HSFE_MonitorResult
    {
        int hit_flag; /// < 比中标志 0:未比中 1:比中 

        int channel_id; /// < 通道ID 

        int face_id; /// < 人脸ID 

        const HSFE_MonitorResultUnit *p_result_unit; /// < 结果集合 
        int unit_count; /// < 结果的数量 

        void *p_user_data; /// < 用户数据 
        int feature_size; /// < 特征长度
        const void *p_feature; /// < 特征
    }HSFE_MonitorResult;

    typedef struct _HSFE_PersonCandidate
    {
        float score;
        int template_id;
    }HSFE_PersonCandidate;

    typedef struct _HSFE_EngineParameters
    {
        int batch_size;
        int chrominance_thresh;
        int collect_time_out_us;
        int small_face_score_decay; /// < score decay value when face is small
        int min_pic_size;
    }HSFE_EngineParameters;


#ifdef __cplusplus
}
#endif

#endif
