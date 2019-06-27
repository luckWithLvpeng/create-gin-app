/*********************************************************** 
* Date: 2016-10-20 
* 
* Author: 牟韵 
* 
* Department: 智能产品中心-研发部 
* 
* Company: 北京海鑫科金高科技股份有限公司 
* 
* Module: 业务数据包 
* 
* Brief: 
* 
* Note: 
* 
* CodePage: Pure UTF-8 
************************************************************/ 
#ifndef __FACE_PACKETS_H_BY_MOUYUN_2016_10_20__
#define __FACE_PACKETS_H_BY_MOUYUN_2016_10_20__

#include <memory>
#include <vector>
#include <map>
#include <list>
#include "my_common_defines.hpp"
#include "my_lock.hpp"
#include "hsfe_types.h"


/// < 比对结果判定器 
class CFaceResultJudger
{
public:
    CFaceResultJudger();
    ~CFaceResultJudger();

private:
    CFaceResultJudger(const CFaceResultJudger&);
    CFaceResultJudger& operator=(const CFaceResultJudger&);

private:
    std::vector<int> m_hit_pyramid; /// < 比中判定 
    std::vector<int> m_reject_pyramid; /// < 拒识判定 

private:
    bool build_hit_pyramid(const int _top, const int _bottom);
    bool build_reject_pyramid(const int _top, const int _bottom);

public:
    bool build(const int hit_top, const int hit_bottom, const int reject_top = 0, const int reject_bottom = 0);

    void destroy();

    bool judge_hit(const int _score, const int _count) const;

    bool judge_reject(const int _score, const int _count) const;
};

/// < 通道参数 
class CChannelParameters
{
public:
    CChannelParameters() {}
    ~CChannelParameters() {}

private:
    CChannelParameters(const CChannelParameters&);
    CChannelParameters& operator=(const CChannelParameters&);

public:
    int m_channel_id; /// < 通道的ID

    int m_channel_mode; /// < 通道模式（核验、监控） 

    int m_reduce_ratio; /// < 视频帧的缩小比例 

    int m_frame_width; /// < 视频帧的宽度 
    int m_frame_height; /// < 视频帧的高度 
    int m_frame_format; /// < 视频帧的格式 
    int m_color_channel_sequence; /// < 彩图的通道序列(BGR or RGB) 
    int m_frame_keep_flag; /// < 是否保持视频帧标志 

    HSFE_Rect m_detect_roi; /// < 检测ROI区域
    int m_face_track_step; /// < 跟踪步长
    int m_max_track_count; /// 每个视频同时跟踪最大人脸数
    int m_face_filter_structure; /// < 人脸脸型过滤阈值
    int m_face_filter_clearity; /// < 人脸清晰度过滤值，小于该值时，权重为清晰度 / face_filter_clearity, 大于时为1
    int m_face_filter_size; /// < 人脸尺寸过滤大小，小于该值时，权重为face_wdith / face_filter_size, 大于时为1
    int m_collect_cache_num;
    int m_pyramid_layer_min_height; /// < 金字塔阈值最小层间距
    int m_face_recorder_remove_dif; /// < face result reporter track_id移除间距
    int m_feature_recorder_remove_dif; /// < face result reporter feature移除间距

    int m_reduced_width; /// < 缩小图的宽度 
    int m_reduced_height; /// < 缩小图的高度 

    int m_face_min_size; /// < （跟踪时的）最小人脸 
    int m_face_max_size; /// < （跟踪时的）最大人脸 

    int m_face_filte_count; /// < 人脸过滤计数 

    int m_face_confidence; /// < 人脸置信度（百分制） 

    int m_pose_estimate_flag; /// < 角度估计的开关 
    int m_yaw_left; /// < 人脸的左偏角 
    int m_yaw_right; /// < 人脸的右偏角 

    int m_prefeature_flag; /// < 核验预提特征标志 

    int m_recog_top_threshold; /// < 识别的顶级阈值 
    int m_recog_bottom_threshold; /// < 识别的基准阈值 

    int m_reject_flag; /// < 拒识的开关 
    int m_reject_top_threshold; /// < 拒识的顶级阈值 
    int m_reject_bottom_threshold; /// < 拒识的基准阈值 

    int m_max_result_count; /// < 监控结果的最大保留数量 

    int m_merge_time_out_ms; /// < 监控结果的融合超时时间 

    std::vector<int> m_sublib_for_monitor; /// < 监控的分库 
    std::map<int,int> m_sublib_quality; /// < 监控的分库质量

    CFaceResultJudger m_monitor_result_judger; /// < 监控结果判定器 
    int m_liveness_thresh;

public:
    inline bool is_invalid_mode() const
    {
        return m_channel_mode == HSFE_CM_INVALID;
    }

    inline bool is_keep_frame() const
    {
        return m_frame_keep_flag != 0;
    }

    inline bool is_verify() const
    {
        return (m_channel_mode & HSFE_CM_VERIFY) == HSFE_CM_VERIFY;
    }

    inline bool is_monitor() const
    {
        return is_monitor_only_max_face() || is_monitor_all();
    }

    inline bool is_monitor_only_max_face() const
    {
        return (m_channel_mode & HSFE_CM_ONLY_MONITOR_MAX_FACE) == HSFE_CM_ONLY_MONITOR_MAX_FACE;
    }

    inline bool is_monitor_all() const
    {
        return (m_channel_mode & HSFE_CM_MONITOR_ALL) == HSFE_CM_MONITOR_ALL;
    }

    inline bool is_estimate_pose() const
    {
        return m_pose_estimate_flag != 0;
    }

    inline bool is_valid_yaw(const float _yaw) const
    {
        return _yaw > m_yaw_left && _yaw < m_yaw_right;
    }

    inline bool is_rejection_open() const
    {
        return m_reject_flag != 0;
    }

    inline bool is_prepare_feature() const
    {
        return m_prefeature_flag != 0;
    }
};

/// < 跟踪到的人脸 
class CTrackedFace
{
public:
    CTrackedFace() { clear(); }
    ~CTrackedFace() {}

    CTrackedFace(const CTrackedFace &_copy) :
        m_id(_copy.m_id),
        m_left(_copy.m_left),
        m_top(_copy.m_top),
        m_right(_copy.m_right),
        m_bottom(_copy.m_bottom),
        m_confidence(_copy.m_confidence)
    {
    }

    CTrackedFace& operator=(const CTrackedFace &_rht)
    {
        if (&_rht != this)
        {
            m_id = _rht.m_id;
            m_left = _rht.m_left;
            m_top = _rht.m_top;
            m_right = _rht.m_right;
            m_bottom = _rht.m_bottom;
            m_confidence = _rht.m_confidence;
        }

        return *this;
    }

    bool operator<(const CTrackedFace &_rht)
    {
        return (m_right - m_left) * (m_bottom - m_top) > (_rht.m_right - _rht.m_left) * (_rht.m_bottom - _rht.m_top);
    }

public:
    int m_id;
    int m_left;
    int m_top;
    int m_right;
    int m_bottom;
    float m_confidence;

public:
    void clear()
    {
        m_id = 0;
        m_left = 0;
        m_top = 0;
        m_right = 0;
        m_bottom = 0;
        m_confidence = 0.00f;
    }

    inline int get_face_width() const
    {
        return m_right - m_left;
    }

    inline int get_face_height() const
    {
        return m_bottom - m_top;
    }
};

/// < face snap 
class CSnappedFace
{
public:
    CSnappedFace() { clear(); }
    ~CSnappedFace() {}

public:
    int m_frame_x0;  // face bbox relative to frame
    int m_frame_y0;
    int m_frame_x1;
    int m_frame_y1;
    int m_snap_x0;  // face bbox relative to snap image
    int m_snap_y0;
    int m_snap_x1;
    int m_snap_y1;
    int m_frame_width;
    int m_frame_height;
    int m_width;
    int m_height;
    int m_track_id;
    int m_format;
    int m_quality;
    std::shared_ptr<CBuffer> m_sp_face_buf;

public:
    void clear()
    {
        m_frame_x0 = 0;
        m_frame_y0 = 0;
        m_frame_x1 = 0;
        m_frame_y1 = 0;
        m_snap_x0 = 0;
        m_snap_y0 = 0;
        m_snap_x1 = 0;
        m_snap_y1 = 0;
        m_frame_width = 0;
        m_frame_height = 0;
        m_width = 0;
        m_height = 0;
        m_track_id = 0;
        m_format = HSFE_IMG_FORMAT_I420;  // I420 by default
        m_quality = 0;
        m_sp_face_buf.reset();
    }
};


/// < 视频帧数据包 
class CFramePacket
{
public:
    CFramePacket() { clear(); }
    ~CFramePacket() {}

private:
    CFramePacket(const CFramePacket&);
    CFramePacket& operator=(const CFramePacket&);

public:
    int m_channel_id; /// < 通道ID 
    int m_frame_id; /// < 帧ID 

    int m_width;
    int m_height;

    std::shared_ptr<CBuffer> m_sp_rgb; /// < RGB图 
    std::shared_ptr<CBuffer> m_sp_yuv; /// < YUV图 
    std::shared_ptr<CBuffer> m_sp_gray; /// < 灰度图 

public:
    void clear()
    {
        m_channel_id = 0;
        m_frame_id = 0;

        m_width = 0;
        m_height = 0;

        m_sp_rgb.reset();
        m_sp_yuv.reset();
        m_sp_gray.reset();
    }
};

/// < 比对结果 
class CCompareResult
{
public:
    CCompareResult() { clear(); }
    ~CCompareResult() {};

private:
    CCompareResult(const CCompareResult&);
    CCompareResult& operator=(const CCompareResult&);

public:
    /// < 重载小于操作符，在其中实现大于逻辑，以此来实现降序排列 
    bool operator<(const CCompareResult &_rht)
    {
        /// < 比较策略： 比分次数优先考虑，比中次数其次考虑 
        if (m_best_score > _rht.m_best_score)
        {
            return true;
        }
        else if (m_best_score == _rht.m_best_score)
        {
            return m_hit_count > _rht.m_hit_count;
        }
        else
        {
            return false;
        }
    }

public:
    int m_template_id; /// < 比中模板的ID 
    int m_sublib_id; /// < 比中模板所属分库的ID 

    int m_best_score; /// < 最高比分 
    int m_hit_count; /// < 比中次数 

    //std::shared_ptr<CSnappedFace> m_sp_snap;
    //std::shared_ptr<CFramePacket> m_sp_frame;

public:
    void clear()
    {
        m_template_id = 0;
        m_sublib_id = 0;

        m_best_score = 0;
        m_hit_count = 0;

        //m_sp_snap.reset();
        //m_sp_frame.reset();
    }
};

/// < 人脸数据包 
class CFacePacket
{
public:
    CFacePacket() { clear(); }
    ~CFacePacket() {}

private:
    CFacePacket(const CFacePacket&);
    CFacePacket& operator=(const CFacePacket&);

public:
    int m_process_result; /// < 流程执行结果 

    long long m_create_time_ms; /// < 创建时间 
    long long m_consumed_time_us; /// < 流程耗时 

    int m_max_face_flag; /// < 是否为最大人脸的标志 

    int m_face_id; /// < 人脸ID 
    int m_face_vanish_count;

    //CTrackedFace m_face_rect; /// < 人脸坐标 
    std::list<std::shared_ptr<CSnappedFace>> m_sp_snapped_face_list;

    float m_pitch; /// < 俯仰角度 
    float m_yaw; /// < 水平角度 

    std::shared_ptr<CBuffer> m_sp_feature_points; /// < 特征点数据 
    std::shared_ptr<CBuffer> m_sp_feature; /// < 特征数据 

    std::shared_ptr<CFramePacket> m_sp_frame; /// < 人脸所在的视频帧 

    std::list<std::shared_ptr<CCompareResult> > m_result_list; /// < 监控结果 

    std::shared_ptr<CChannelParameters> m_sp_channel_parameters; /// < 人脸所属通道的参数 

public:
    void clear()
    {
        m_process_result = 0;

        m_create_time_ms = 0;
        m_consumed_time_us = 0;

        m_max_face_flag = 0;

        m_face_id = 0;
        m_face_vanish_count = 0;

        //m_face_rect.clear();
        m_sp_snapped_face_list.clear();

        m_pitch = 0.00f;
        m_yaw = 0.00f;

        m_sp_feature_points.reset();
        m_sp_feature.reset();

        m_sp_frame.reset();

        m_result_list.clear();

        m_sp_channel_parameters.reset();
    }

    inline bool is_max_face() const
    {
        return m_max_face_flag != 0;
    }
};

////////////////////////////////////////////////////////////////////////// 

/// < 比对数据包 
class CComparePacket
{
public:
    CComparePacket() { clear(); }
    ~CComparePacket() {}

private:
    CComparePacket(const CComparePacket&);
    CComparePacket& operator=(const CComparePacket&);

public:
    int m_process_result; /// < 流程执行结果 

    long long m_consumed_time_us; /// < 流程耗时 

    std::list<std::shared_ptr<CFacePacket> > m_face_set; /// < 人脸集合 

public:
    void clear()
    {
        m_process_result = 0;
        m_consumed_time_us = 0;
        m_face_set.clear();
    }
};

/// < 跟踪数据包 
class CTrackPacket
{
public:
    CTrackPacket() { clear(); }
    ~CTrackPacket() {}

private:
    CTrackPacket(const CTrackPacket&);
    CTrackPacket& operator=(const CTrackPacket&);

public:
    int m_process_result; /// < 流程执行结果 

    long long m_consumed_time_us; /// < 流程耗时 

    std::shared_ptr<CFramePacket> m_sp_frame;

    std::list<CTrackedFace> m_tracked_face_list;
    std::map<int, std::list<std::shared_ptr<CSnappedFace>>> m_sp_snapped_face_map;

    std::shared_ptr<CChannelParameters> m_sp_channel_parameters; /// < 通道参数 

public:
    void clear()
    {
        m_process_result = 0;

        m_consumed_time_us = 0;

        m_sp_frame.reset();

        m_tracked_face_list.clear();
        m_sp_snapped_face_map.clear();

        m_sp_channel_parameters.reset();
    }
};

/// < 核验模板 
class CVerifyTemplate
{
public:
    CVerifyTemplate() { clear(); }
    ~CVerifyTemplate() {}

private:
    DISABLE_CLASS_COPY(CVerifyTemplate);

public:
    long long m_create_time_ms; /// < 创建时间 
    int m_time_out_ms; /// < 超时时间 单位ms 

    std::vector<int> m_channel_to_verify; /// < 核验的通道 

    int m_reject_enable_flag;

    int m_top_threshold; /// < 顶级阈值 
    int m_bottom_threshold; /// < 基准阈值 
    CFaceResultJudger m_judger; /// < 判定器 

    std::list<std::shared_ptr<CBuffer> > m_feature_list; /// < 特征队列 

    int m_hit_flag; /// < 是否比中的标志 
    int m_hit_count; /// < 比中次数 
    int m_miss_count; /// < 未比中次数 
    int m_best_score; /// < 最佳分数 

    int m_face_id; /// < 人脸ID 
    CTrackedFace m_face_rect; /// < 人脸坐标 

    std::shared_ptr<CFramePacket> m_sp_frame;

public:
    void clear()
    {
        m_create_time_ms = 0;
        m_time_out_ms = 0;

        m_channel_to_verify.clear();

        m_reject_enable_flag = 0;

        m_top_threshold = 0;
        m_bottom_threshold = 0;
        m_judger.destroy();

        m_feature_list.clear();

        m_hit_flag = 0;
        m_hit_count = 0;
        m_miss_count = 0;
        m_best_score = 0;

        m_face_id = 0;
        m_face_rect.clear();

        m_sp_frame.reset();
    }

public:
    inline bool is_reject_enable() const
    {
        return m_reject_enable_flag != 0;
    }
};

////////////////////////////////////////////////////////////////////////// 

/// < 核验结果 
class CVerifyResult
{
public:
    CVerifyResult() { clear(); }
    ~CVerifyResult() {}

private:
    CVerifyResult(const CVerifyResult&);
    CVerifyResult& operator=(const CVerifyResult&);

public:
    int m_consumed_time_ms; /// < 从输入核验数据至得到核验结果的耗时 

    int m_hit_flag; /// < 比中标志 
    int m_score; /// < 比分 

    int m_face_id; /// < 比中的人脸的ID 
    CTrackedFace m_face_rect; /// < 人脸坐标 
    std::shared_ptr<CFramePacket> m_sp_frame; /// < 人脸所属的视频帧 

public:
    void clear()
    {
        m_consumed_time_ms = 0;

        m_hit_flag = 0;
        m_score = 0;

        m_face_id = 0;
        m_face_rect.clear();
        m_sp_frame.reset();
    }
};

/// < 监控结果融合包 
class CMonitorResultMergePacket
{
public:
    CMonitorResultMergePacket() { clear(); };
    ~CMonitorResultMergePacket() {};

private:
    CMonitorResultMergePacket(const CMonitorResultMergePacket&);
    CMonitorResultMergePacket& operator=(const CMonitorResultMergePacket&);

public:
    long long m_create_time_ms; /// < 创建时间 
    int m_merge_time_out_ms; /// < 融合超时时间 
    int m_merge_count; /// < 融合次数 

    int m_hit_flag; /// < 比中标志 

    int m_channel_id; /// < 通道ID 
    int m_face_id; /// < 人脸ID 

    std::shared_ptr<CFacePacket> m_sp_face_pkt;

public:
    void clear()
    {
        m_create_time_ms = 0;
        m_merge_time_out_ms = 0;
        m_merge_count = 0;

        m_hit_flag = 0;

        m_channel_id = 0;
        m_face_id = 0;

        m_sp_face_pkt.reset();
    }
};

int hsfe_input_packet(std::shared_ptr<CTrackPacket> &sp_pkt);

#endif

