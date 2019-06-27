#ifndef HISIGN_HOBOT_IPC_CLIENT_INTERFACE_H_
#define HISIGN_HOBOT_IPC_CLIENT_INTERFACE_H_

#include <list>
#include <string>
#include "my_common_defines.hpp"
#include "my_buffer.hpp"

namespace HobotIpcClientInterface {

    typedef enum _ErrorCode {
        EC_UNKNOWN_ERROR = (-1),
        EC_NO_ERROR = 0,
        EC_INVALID_PARAMETER = 1,
        EC_CLI_NOT_INITED = 2,
        EC_CLI_ALREADY_INITED = 3,
        EC_THREAD_ERROR = 4,
        EC_ENGINE_INIT_ERROR = 5,
        EC_ENGINE_NOT_INITED = 6,
        EC_ADD_CLI_ERROR = 7,
        EC_NVDEC_INIT_ERROR = 8,
        EC_MEMORY_ERROR = 9,
        EC_GET_LOGGER_ERROR = 10
    } ErrorCode;

    typedef struct _ConfigParam {
        std::string url;
        int channel_id;
        int skip_frames;
        int frame_width;
        int frame_height;
    } ConfigParam;

    typedef struct _Point {
        int x;
        int y;
        float score;
    } Point;

    typedef struct _Attribute {
        int type;
        int int_val;
        std::string str_val;
    } Attribute;

    class CSnapPacket {
        public:
            CSnapPacket() {
                clear();
            }
            ~CSnapPacket() {}

        private:
            CSnapPacket(const CSnapPacket&);
            CSnapPacket& operator=(const CSnapPacket&);

        public:
            int person_id_;
            std::string img_type_;
            int buf_size_;
            std::shared_ptr<CBuffer> sp_buf_;
            int face_box_top_;
            int face_box_left_;
            int face_box_right_;
            int face_box_bottom_;
            int pic_box_top_;
            int pic_box_left_;
            int pic_box_right_;
            int pic_box_bottom_;

        private:
            void clear() {
                buf_size_ = 0;
                sp_buf_.reset();
                face_box_top_ = 0;
                face_box_left_ = 0;
                face_box_right_ = 0;
                face_box_bottom_ = 0;
                pic_box_top_ = 0;
                pic_box_left_ = 0;
                pic_box_right_ = 0;
                pic_box_bottom_ = 0;
            }
    };

    class CTrackPacket {
        public:
            CTrackPacket() {
                clear();
            }
            ~CTrackPacket() {}

        private:
            CTrackPacket(const CTrackPacket&);
            CTrackPacket& operator=(const CTrackPacket&);

        public:
            int track_id_;
            int bbox_top_;
            int bbox_left_;
            int bbox_right_;
            int bbox_bottom_;
            std::list<Point> landmarks_;
            std::list<Attribute> attributes_;
        private:
            void clear() {
                track_id_ = 0;
                bbox_top_ = 0;
                bbox_left_ = 0;
                bbox_right_ = 0;
                bbox_bottom_ = 0;
                landmarks_.clear();
                attributes_.clear();
            }
    };

    class CIpcPacket {
        public:
            CIpcPacket() {
                clear();
            }
            ~CIpcPacket() {}

        private:
            CIpcPacket(const CIpcPacket&);
            CIpcPacket& operator=(const CIpcPacket&);

        public:
            int channel_id_;
            int frame_id_;

            int frame_width_;
            int frame_height_;
            int frame_size_;

            std::shared_ptr<CBuffer> sp_frame_buf_;
            std::list<std::shared_ptr<CSnapPacket>> snap_list_;
            std::list<std::shared_ptr<CTrackPacket>> track_list_;
        public:
            void clear() {
                channel_id_ = 0;
                frame_id_ = 0;

                frame_width_ = 0;
                frame_height_ = 0;
                frame_size_ = 0;

                sp_frame_buf_.reset();
                snap_list_.clear();
                track_list_.clear();
            }
    };

    typedef void(*PFUNC_PACKET_ARRIVE)(std::shared_ptr<CIpcPacket> sp_pkt);
    int init_ipc_engine();
    int close_ipc_engine();
    void* create_ipc_client(const ConfigParam*, int *p_ipc_client);
    void destroy_ipc_client(void *p_ipc_client);
    void set_pkt_arrive_callback(void *p_source, PFUNC_PACKET_ARRIVE cb);
    int start(void *p_source);

}

#endif  // HISIGN_HOBOT_IPC_CLIENT_INTERFACE_H_

