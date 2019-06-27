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
#define HS_DLL_EXPORT 1

#include <list>
#include <unordered_map>
#include <iostream>
#include "hsfe_wrapper.h"
#include "face_packets.h"
#include "my_common_defines.hpp"
#include "my_object_pool.hpp"
#ifdef WITH_CVDEC
#include "cvdec_interface.h"
#endif
#ifdef WITH_HOBOT_IPC
#include "hobot_ipc_client_interface.h"
#endif

int g_hsvd_error_code = HSVD_EC_NO_ERROR;
int g_hsfe_error_code = HSFE_EC_NO_ERROR;

std::list<void*> g_decoder_list;
std::unordered_map<int, void*> g_rtsp_server_map;
std::unordered_map<int, void*> g_rtmp_map;
std::unordered_map<int, void*> g_hobot_ipc_client_map;
std::unordered_map<int, void*> g_cvdec_map;
void* g_videowriter_map[32] = {nullptr};
void* g_webrtc_map[32] = {nullptr};
int g_webrtc_type_map[32] = {0};

DEFINE_OBJECT_POOL(CTrackPacket, TrackPacketPool);
DEFINE_OBJECT_POOL(CFramePacket, FramePacketPool);
DEFINE_OBJECT_POOL(CSnappedFace, SnappedFacePool);

void HSVD_STDCALL frame_arrive(const HSVD_Frame *p_frame, void *user_data)
{
    if (p_frame == nullptr)
    {
        return;
    }

    HSFE_Frame _frame;
    _frame.channel_id = p_frame->decoder_id;
    _frame.frame_id = 0;
    _frame.p_buf = p_frame->p_buf;

    hsfe_input_frame(&_frame);
}

#ifdef WITH_HOBOT_IPC
void pkt_arrive(std::shared_ptr<HobotIpcClientInterface::CIpcPacket> sp_pkt) {
    //std::cout << "*************** pkt arrive callback ***************" << std::endl;
    //std::cout << "frame width->" << sp_pkt->frame_width_ << " height->"<< sp_pkt->frame_height_ << " size->" << sp_pkt->frame_size_ << std::endl;
    //std::cout << "track list size->" << sp_pkt->track_list_.size() << std::endl;

    auto sp_frame = FramePacketPool::pop_sp();
    if (sp_frame == nullptr) {
        std::cout << "memory error!" << std::endl;
        return;
    }
    sp_frame->m_channel_id = sp_pkt->channel_id_;
    sp_frame->m_frame_id = sp_pkt->frame_id_;
    sp_frame->m_width = sp_pkt->frame_width_;
    sp_frame->m_height = sp_pkt->frame_height_;
    sp_frame->m_sp_yuv = sp_pkt->sp_frame_buf_;
    sp_frame->m_sp_gray = sp_pkt->sp_frame_buf_;
    sp_frame->m_sp_rgb.reset();

    auto sp_track_pkt = TrackPacketPool::pop_sp();
    if (sp_track_pkt == nullptr) {
        std::cout << "memory error!" << std::endl;
        return;
    }
    sp_track_pkt->m_sp_frame = sp_frame;

    for (auto & sp_track : sp_pkt->track_list_) {
        CTrackedFace _face;

        _face.m_id = sp_track->track_id_;
        _face.m_left = sp_track->bbox_left_;
        _face.m_top = sp_track->bbox_top_;
        _face.m_right = sp_track->bbox_right_;
        _face.m_bottom = sp_track->bbox_bottom_;
        _face.m_confidence = static_cast<float>(1000) / 1000;

        if (_face.m_left < 0) { _face.m_left = 0; }
        if (_face.m_top < 0) { _face.m_top = 0; }
        if (_face.m_right >= sp_frame->m_width ) { _face.m_right = sp_frame->m_width - 1; }
        if (_face.m_bottom >= sp_frame->m_height) { _face.m_bottom = sp_frame->m_height - 1; }

        _face.m_confidence = 0.2 * _face.m_confidence + 0.8;

        sp_track_pkt->m_tracked_face_list.emplace_back(_face);
    }

    for (auto & sp_snap : sp_pkt->snap_list_) {
        auto sp_snapped_face = SnappedFacePool::pop_sp();
        if (sp_snapped_face == nullptr) {
            std::cout << "memory error!" << std::endl;
            return;
        }
        sp_snapped_face->m_frame_x0 = sp_snap->pic_box_left_;
        sp_snapped_face->m_frame_y0 = sp_snap->pic_box_top_;
        sp_snapped_face->m_frame_x1 = sp_snap->pic_box_right_;
        sp_snapped_face->m_frame_y1 = sp_snap->pic_box_bottom_;
        sp_snapped_face->m_snap_x0 = sp_snap->face_box_left_;
        sp_snapped_face->m_snap_y0 = sp_snap->face_box_top_;
        sp_snapped_face->m_snap_x1 = sp_snap->face_box_right_;
        sp_snapped_face->m_snap_y1 = sp_snap->face_box_bottom_;
        //sp_snapped_face->m_width = snaps[i].width_;
        //sp_snapped_face->m_height = snaps[i].height_;
        sp_snapped_face->m_track_id = sp_snap->person_id_;
        //sp_snapped_face->m_format = snaps[i].format_;
        sp_snapped_face->m_quality = sp_snap->buf_size_;

        sp_snapped_face->m_sp_face_buf = sp_snap->sp_buf_;
        sp_track_pkt->m_sp_snapped_face_map[sp_snapped_face->m_track_id].emplace_back(sp_snapped_face);
    }


    hsfe_input_packet(sp_track_pkt);
}
#endif

int create(const HSFE_EngineParameters *p_engine_parameters)
{
    int _result = hsfe_create(p_engine_parameters);
    if (_result != HSFE_EC_NO_ERROR)
    {
        g_hsfe_error_code = _result;
        return HSFE_WRAPPER_EC_ENGINE_ERROR;
    }

    _result = hsvd_environment_init();
    if (_result != HSVD_EC_NO_ERROR)
    {
        hsfe_destroy();
        g_hsvd_error_code = _result;
        return HSFE_WRAPPER_EC_DECODER_ERROR;
    }

    g_decoder_list.clear();
#ifdef WITH_HOBOT_IPC
    g_hobot_ipc_client_map.clear();
#endif
    g_rtmp_map.clear();
    g_rtsp_server_map.clear();

#ifdef WITH_CVDEC
    _result = cvdec_environment_init();
    if (_result != CVDEC_EC_NO_ERROR) {
        hsfe_destroy();
        return HSFE_WRAPPER_EC_DECODER_ERROR;
    }
    g_cvdec_map.clear();
#endif

    return HSFE_WRAPPER_EC_NO_ERROR;
}

void destroy()
{
    hsfe_destroy();

    for (auto & p_decoder : g_decoder_list)
    {
        hsvd_destroy(p_decoder);
    }
    g_decoder_list.clear();

    for (std::unordered_map<int, void*>::iterator it = g_rtmp_map.begin(); it != g_rtmp_map.end(); it++) {
        hsrtmp_destroy(it->second);
        g_rtmp_map.erase(it->first);
    }
#ifdef WITH_CVDEC
     for (std::unordered_map<int, void*>::iterator it = g_cvdec_map.begin(); it != g_cvdec_map.end(); it++) {
        cvdec_destroy(it->second);
        g_cvdec_map.erase(it->first);
    }
#endif
}

#ifdef WITH_HOBOT_IPC
int add_hobotipc_client(const HSVD_DecoderParameters *p_decoder_params) {
    int error_code = 0;
    if (g_hobot_ipc_client_map.empty()) {
        std::cout << "appwrapper: init ipc engine......" << std::endl;
        error_code = HobotIpcClientInterface::init_ipc_engine();
        if (error_code != 0) {
            std::cout << "init_ipc_engine failed. error: " << error_code << std::endl;
            return HSFE_WRAPPER_EC_HOBOT_IPC_ERROR;
        }
    }
    std::cout << "appwrapper: creating ipc client......" << std::endl;

    HobotIpcClientInterface::ConfigParam cp;
    cp.url = p_decoder_params->p_video_url;
    cp.channel_id = p_decoder_params->decoder_id;
    cp.frame_width = p_decoder_params->width;
    cp.frame_height = p_decoder_params->height;

    void *p_client = HobotIpcClientInterface::create_ipc_client(&cp, &error_code);

    if (error_code != 0) {
        std::cout << "CIpcClient initialize error: " << error_code << std::endl;
        return HSFE_WRAPPER_EC_HOBOT_IPC_ERROR;
    }

    HobotIpcClientInterface::set_pkt_arrive_callback(p_client, pkt_arrive);

    error_code = HobotIpcClientInterface::start(p_client);
    if (error_code != 0) {
        std::cout << "CIpcClient start error: " << error_code << std::endl;
        return HSFE_WRAPPER_EC_HOBOT_IPC_ERROR;
    }

    g_hobot_ipc_client_map[p_decoder_params->decoder_id] = p_client;

    return HSFE_WRAPPER_EC_NO_ERROR;
}
#endif

#ifdef WITH_CVDEC
int add_cvdec(const HSVD_DecoderParameters *p_decoder_params) {
    std::string source_url = p_decoder_params->p_video_url;
    source_url = source_url.substr(6);
    CVDEC_DecoderParameters cvdec_params;
    memcpy(&cvdec_params, p_decoder_params, sizeof(CVDEC_DecoderParameters));
    cvdec_params.p_video_url = source_url.c_str();

    int error_code = CVDEC_EC_NO_ERROR;
    void *p_cvdec = cvdec_create(&cvdec_params, &error_code);
    if (error_code != CVDEC_EC_NO_ERROR) {
        g_hsvd_error_code = error_code;
        return HSFE_WRAPPER_EC_DECODER_ERROR;
    }
    cvdec_set_frame_arrive_callback(p_cvdec, (PFUNC_CVDEC_FRAME_ARRIVE)frame_arrive, nullptr);
    int _result = cvdec_start(p_cvdec);
    if (_result != CVDEC_EC_NO_ERROR) {
        cvdec_destroy(p_cvdec);
        g_hsvd_error_code = _result;
        return HSFE_WRAPPER_EC_DECODER_ERROR;
    } else {
        g_cvdec_map[p_decoder_params->decoder_id] = p_cvdec;
    }
    return HSFE_WRAPPER_EC_NO_ERROR;
}
#endif

int add_channel(const HSVD_DecoderParameters *p_decoder_params, const HSFE_ChannelParameters *p_engine_params)
{
    std::string source_url = p_decoder_params->p_video_url;
    if (source_url.find("hobotipc") == 0) {
        const_cast<HSFE_ChannelParameters*>(p_engine_params)->use_tracker = 0;
        const_cast<HSFE_ChannelParameters*>(p_engine_params)->face_filte_count = 1;
    } else {
        const_cast<HSFE_ChannelParameters*>(p_engine_params)->use_tracker = 1;
    }
    /// < 初始化通道 
    printf("appwrapper: initializing hsfe channel...... \n");
    int _result = hsfe_add_channel(p_engine_params);
    if (_result != HSFE_EC_NO_ERROR)
    {
        g_hsfe_error_code = _result;
        return HSFE_WRAPPER_EC_ENGINE_ERROR;
    }

#ifdef WITH_HOBOT_IPC
    if (p_engine_params->use_tracker == 0) {
        _result = add_hobotipc_client(p_decoder_params);
        if (_result != HSFE_WRAPPER_EC_NO_ERROR) hsfe_remove_channel(p_engine_params->channel_id);
        return _result;
    }
#endif

#ifdef WITH_CVDEC
    if (source_url.find("cvdec-") == 0) {
        _result = add_cvdec(p_decoder_params);
        if (_result != HSFE_WRAPPER_EC_NO_ERROR) hsfe_remove_channel(p_engine_params->channel_id);
        return _result;
    }
#endif

    /// < 创建解码器 
    printf("appwrapper: creating decoder...... \n");
    int error_code = HSVD_EC_NO_ERROR;
    void *p_decoder = hsvd_create(p_decoder_params, &error_code);
    if (error_code != HSVD_EC_NO_ERROR)
    {
        hsfe_remove_channel(p_engine_params->channel_id);
        g_hsvd_error_code = error_code;
        return HSFE_WRAPPER_EC_DECODER_ERROR;
    }

    hsvd_set_frame_arrive_callback(p_decoder, frame_arrive, nullptr);

    /// < 启动解码 
    printf("appwrapper: starting decoder...... \n");
    _result = hsvd_start(p_decoder);
    if (_result != HSVD_EC_NO_ERROR)
    {
        hsvd_destroy(p_decoder);
        hsfe_remove_channel(p_engine_params->channel_id);
        g_hsvd_error_code = _result;
        return HSFE_WRAPPER_EC_DECODER_ERROR;
    }
    else
    {
        g_decoder_list.push_back(p_decoder);
    }

    return HSFE_WRAPPER_EC_NO_ERROR;
}

void remove_channel(const int channel_id)
{
    for (auto d_itr = g_decoder_list.begin(); d_itr != g_decoder_list.end(); ++d_itr)
    {
        int decoder_id = 0;
        if (hsvd_get_decoder_id(*d_itr, &decoder_id) == HSVD_EC_NO_ERROR && decoder_id == channel_id)
        {
            hsvd_destroy(*d_itr);
            g_decoder_list.erase(d_itr);
            break;
        }
    }
#ifdef WITH_HOBOT_IPC
    if (g_hobot_ipc_client_map.count(channel_id) > 0) {
        std::cout << "before delete" << std::endl;
        HobotIpcClientInterface::destroy_ipc_client(g_hobot_ipc_client_map[channel_id]);
        std::cout << "after delete" << std::endl;
        g_hobot_ipc_client_map.erase(channel_id);
        if (g_hobot_ipc_client_map.empty()) {
            HobotIpcClientInterface::close_ipc_engine();
        }
    }
#endif
#ifdef WITH_CVDEC
    if (g_cvdec_map.count(channel_id) > 0) {
        cvdec_destroy(g_cvdec_map[channel_id]);
        g_cvdec_map.erase(channel_id);
    }
#endif
    std::cout << "before remove hsfe channel" << std::endl;
    hsfe_remove_channel(channel_id);
    std::cout << "after remove hsfe channel" << std::endl;
}

int get_decoder_error_code()
{
    return g_hsvd_error_code;
}

int get_engine_error_code()
{
    return g_hsfe_error_code;
}

// Use decoder parameters here to be compitable with eme/middlelayer
int add_rtmp_channel(const HSVD_DecoderParameters *p_decoder_params)
{
    /// < Create RTMP 
    printf("appwrapper: creating rtmp...... \n");
    HSRTMP_Parameters rtmp_param;
    rtmp_param.id = p_decoder_params->decoder_id;
    rtmp_param.format = p_decoder_params->decode_format == HSVD_DECODE_TYPE_I420 ? HSRTMP_FRAME_FORMAT_I420 : HSRTMP_FRAME_FORMAT_BGR;
    rtmp_param.width = p_decoder_params->width/2;
    rtmp_param.height = p_decoder_params->height/2;
    rtmp_param.size = rtmp_param.width*rtmp_param.height*3;
    //rtmp_param.fps = 24/(p_decoder_params->frame_skip_num + 1);
    rtmp_param.fps = 25;
    //rtmp_param.compress_ratio = 4;
    std::string url = "rtmp://localhost/live/" + std::to_string(rtmp_param.id);
    rtmp_param.p_url = url.c_str();
    int error_code = HSRTMP_EC_NO_ERROR;
    void *p_rtmp = hsrtmp_create(&rtmp_param, &error_code);
    if (error_code != HSRTMP_EC_NO_ERROR) {
        hsrtmp_destroy(p_rtmp);
        printf("appwrapper: creating rtmp failed...... \n");
        return HSFE_WRAPPER_EC_RTMP_ERROR;
    }
    g_rtmp_map[rtmp_param.id] = p_rtmp;

    return HSFE_WRAPPER_EC_NO_ERROR;
}

void remove_rtmp_channel(const int channel_id)
{
    hsrtmp_destroy(g_rtmp_map[channel_id]);
    g_rtmp_map.erase(channel_id);
}

void push_rtmp_and_draw_bbox(const int channel_id, HSFE_Image *img, HSFE_FaceRect *p_rects, const int rect_count)
{
    hsrtmp_push_and_draw_bbox(g_rtmp_map[channel_id], (HSRTMP_Image*)img, (HSRTMP_FaceRect*)p_rects, rect_count);
}

int add_webrtc_channel(const GST_WEBRTC_Parameters *p_params, PFUNC_GST_WEBRTC_SAMPLE_ARRIVE cb, const char *p_url, int use_rtmp)
{
    /// < Create webrtc
    printf("appwrapper: creating gst webrtc...... \n");
    if (use_rtmp) {
        HSRTMP_Parameters rtmp_param;
        rtmp_param.id = p_params->id;
        rtmp_param.format = p_params->format == GST_WEBRTC_FRAME_FORMAT_I420 ? HSRTMP_FRAME_FORMAT_I420 : HSRTMP_FRAME_FORMAT_BGR;
        rtmp_param.width = p_params->width;
        rtmp_param.height = p_params->height/2;
        rtmp_param.size = rtmp_param.width*rtmp_param.height*3;
        rtmp_param.fps = p_params->fps;
        rtmp_param.p_url = p_url;
        int error_code = HSRTMP_EC_NO_ERROR;
        void *p_rtmp = hsrtmp_create(&rtmp_param, &error_code);
        if (error_code != HSRTMP_EC_NO_ERROR) {
            hsrtmp_destroy(p_rtmp);
            printf("appwrapper: creating rtmp failed...... \n");
            return HSFE_WRAPPER_EC_RTMP_ERROR;
        }
        g_webrtc_type_map[rtmp_param.id] = 1;
        g_webrtc_map[rtmp_param.id] = p_rtmp;
    } else {
        int error_code = GST_WEBRTC_EC_NO_ERROR;
        void *p_webrtc = gst_webrtc_create(p_params, &error_code);
        if (error_code != GST_WEBRTC_EC_NO_ERROR) {
            gst_webrtc_destroy(p_webrtc);
            printf("appwrapper: creating gst-webrtc failed...... \n");
            return HSFE_WRAPPER_EC_RTMP_ERROR;
        }
        gst_webrtc_set_sample_arrive_callback(p_webrtc, cb);
        g_webrtc_type_map[p_params->id] = 0;
        g_webrtc_map[p_params->id] = p_webrtc;
    }
    return HSFE_WRAPPER_EC_NO_ERROR;
}

void remove_webrtc_channel(const int channel_id)
{
    if (g_webrtc_map[channel_id] == nullptr) return;
    if (g_webrtc_type_map[channel_id]) {
        hsrtmp_destroy(g_webrtc_map[channel_id]);
    } else {
        gst_webrtc_destroy(g_webrtc_map[channel_id]);
    }
    g_webrtc_map[channel_id] = nullptr;
    g_webrtc_type_map[channel_id] = 0;
}

void push_webrtc_and_draw_bbox(const int channel_id, HSFE_Image *img, HSFE_FaceRect *p_rects, const int rect_count)
{
    if (g_webrtc_map[channel_id] == nullptr) return;
    if (g_webrtc_type_map[channel_id]) {
        hsrtmp_push_and_draw_bbox(g_webrtc_map[channel_id], (HSRTMP_Image*)img, (HSRTMP_FaceRect*)p_rects, rect_count);
    } else {
        gst_webrtc_push_and_draw_bbox(g_webrtc_map[channel_id], (GST_WEBRTC_Image*)img, (GST_WEBRTC_FaceRect*)p_rects, rect_count);
    }
}

int add_videowriter_channel(const GST_VIDEO_WRITER_Parameters *p_params)
{
    if (g_videowriter_map[p_params->id] != nullptr) return GST_VIDEO_WRITER_EC_NO_ERROR;
    printf("appwrapper: creating gst videowriter...... \n");
    int error_code = GST_VIDEO_WRITER_EC_NO_ERROR;
    void *p_videowriter = gst_videowriter_create(p_params, &error_code);
    if (error_code != GST_VIDEO_WRITER_EC_NO_ERROR) {
        gst_videowriter_destroy(p_videowriter);
        printf("appwrapper: creating videowriter failed...... \n");
        return error_code;
    }
    g_videowriter_map[p_params->id] = p_videowriter;
    return GST_VIDEO_WRITER_EC_NO_ERROR;
}

void remove_videowriter_channel(const int channel_id)
{
    if (g_videowriter_map[channel_id] == nullptr) return;
    gst_videowriter_destroy(g_videowriter_map[channel_id]);
    printf("6!\n");
    g_videowriter_map[channel_id] = nullptr;
    printf("7!\n");
}

void push_videowriter_and_draw_bbox(const int channel_id, HSFE_Image *img, HSFE_FaceRect *p_rects, const int rect_count)
{
    if (g_videowriter_map[channel_id] == nullptr) return;
    gst_videowriter_push_and_draw_bbox(g_videowriter_map[channel_id], (GST_VIDEO_WRITER_Image*)img, (GST_VIDEO_WRITER_FaceRect*)p_rects, rect_count);
}
#ifdef USE_RTSP
int add_rtsp_server(HSRS_Parameters *param)
{
    int error_code = 0;
    printf("appwrapper: creating rtsp server...... \n");
    void *p_rtsp_server = hsrs_create(param, &error_code);
    if (error_code != 0)
    {
        remove_channel(param->id);
        return HSFE_WRAPPER_EC_RTSP_SERVER_ERROR;
    }
    g_rtsp_server_map[param->id] = p_rtsp_server;

    printf("Create RTSP server: id->%d, addr: %p \n", param->id, p_rtsp_server);

    return HSFE_WRAPPER_EC_NO_ERROR;
}

void remove_rtsp_server(const int channel_id)
{
    printf("Remove RTSP server: id->%d, addr: %p \n", channel_id, g_rtsp_server_map[channel_id]);
    hsrs_destroy(g_rtsp_server_map[channel_id]);
    printf("Remove RTSP server: ---------------- \n");
    g_rtsp_server_map.erase(channel_id);
}

void push_rtsp_and_draw_bbox(const int channel_id, HSFE_Image *img, HSFE_FaceRect *p_rects, const int rect_count)
{
    hsrs_push_and_draw_bbox(g_rtsp_server_map[channel_id], (HSRS_Image*)img, (HSRS_FaceRect*)p_rects, rect_count);
}
#endif

