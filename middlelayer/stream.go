// Package mjpeg implements a simple MJPEG streamer.
//
// Stream objects implement the http.Handler interface, allowing to use them with the net/http package like so:
//	stream = mjpeg.NewStream()
//	http.Handle("/camera", stream)
// Then push new JPEG frames to the connected clients using stream.UpdateJPEG().
package middlelayer

/*
#include "middlelayerc.h"
*/
import "C"

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Pipline struct {
	sync.RWMutex
	failcount int32
	ch        chan *[]byte
	id        int
}

// Stream represents a single video feed.
type Stream struct {
	sync.RWMutex
	pipline_maxid int
	channel       map[int]Pipline
	sending				int32
	// FrameInterval time.Duration
}

const boundaryWord = "MJPEGBOUNDARY"
const headerf = "\r\n" +
	"--" + boundaryWord + "\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n" +
	"\r\n"

//var push_channel_list = make(map[int]int)

// https://blog.golang.org/go-maps-in-action
var counter = struct {
	sync.RWMutex
	push_channel_list map[int]int
}{push_channel_list: make(map[int]int)}

func editPushChannelList(channel_id int, add bool) {
	var exist bool = false
	counter.Lock()
	for index, _ := range counter.push_channel_list {
		if index == channel_id {
			exist = true
			if add {
				counter.push_channel_list[index]++
			} else {
				counter.push_channel_list[index]--
			}
			break
		}
	}
	if !exist {
		counter.push_channel_list[channel_id] = 1
	}
	counter.Unlock()
	C.MySetMJPGPushChannelCount((C.int)(channel_id),
		(C.int)(counter.push_channel_list[channel_id]))
}

// ServeHTTP responds to HTTP requests with the MJPEG stream, implementing the http.Handler interface.
func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var id string
	if id = r.Form.Get("id"); len(id) == 0 {
		return
	}
	channel_id, _ := strconv.Atoi(id)

	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundaryWord)

	c := make(chan *[]byte, 1)

	s.Lock()
	s.pipline_maxid += 1
	pipline_maxid := s.pipline_maxid

	var channel = Pipline{id: channel_id, failcount: 0, ch: c}

	s.channel[pipline_maxid] = channel

	editPushChannelList(channel_id, true)
	s.Unlock()

	log.Println("Stream:", r.RemoteAddr, "connected. pipline id:", pipline_maxid)

	for {
		// time.Sleep(s.FrameInterval)
		b := <-c
		_, err := w.Write(*b)
		if err != nil {
			break
		}
	}

	s.Lock()
	if val, ok := s.channel[pipline_maxid]; ok {
		editPushChannelList(val.id, false)
		delete(s.channel, pipline_maxid)
	}
	s.Unlock()

	log.Println("Stream:", r.RemoteAddr, "disconnected. connection id:", pipline_maxid)
}

// UpdateJPEG pushes a new JPEG frame onto the clients.
func (s *Stream) UpdateJPEG(jpeg []byte, channelId int) {
	header := fmt.Sprintf(headerf, len(jpeg))
	frame := make([]byte, (len(jpeg)+len(header))*2)
	copy(frame, header)
	copy(frame[len(header):], jpeg)
	go s.SendJPEG(frame, channelId)
}

func (s *Stream) SendJPEG(frame []byte, channelId int) {
	if is_sending := atomic.LoadInt32(&s.sending); is_sending > 0 {
		C.MyUnlockSendMutex((C.int)(channelId))
		fmt.Println("///MJPG Drop///")
		return
	}
	s.Lock()
	atomic.StoreInt32(&s.sending, 1)
	var pipline_id_todelete int = -1
	//遍历当前channel下所有pipline，推送jpg
	for pipline_id, c := range s.channel {
		// Select to skip streams which are sleeping to drop frames.
		if c.id == channelId {
			//当前connection失败次数
			val := atomic.LoadInt32(&c.failcount)
			if val > 96 {
				pipline_id_todelete = pipline_id
				break
			}
			//如果已经阻塞
			if val > 0 {
				atomic.AddInt32(&c.failcount, 1)
				continue
			}
			//未阻塞，变为阻塞状态，开始发图
			atomic.AddInt32(&c.failcount, 1)
			select {
			case <-time.After(time.Second * 1):
				println("write channel timeout channel_id:", channelId, ", pipline_id:", pipline_id)
				pipline_id_todelete = pipline_id
			case c.ch <- &frame:
			}
			if pipline_id_todelete >= 0 {
				break
			}
			//解除阻塞
			atomic.StoreInt32(&c.failcount, 0)
		}
	}
	if pipline_id_todelete >= 0 {
		if val, ok := s.channel[pipline_id_todelete]; ok {
			editPushChannelList(val.id, false)
			delete(s.channel, pipline_id_todelete)
		}
	}
	atomic.StoreInt32(&s.sending, 0)
	s.Unlock()

	C.MyUnlockSendMutex((C.int)(channelId))
	///< updata
	SyncIpcStatus(channelId)
}

// NewStream initializes and returns a new Stream.
func NewStream() *Stream {
	return &Stream{
		//m:			 make(map[chan []byte]bool),
		pipline_maxid: -1,
		channel:       make(map[int]Pipline),
		sending:       0,
		// FrameInterval: 5 * time.Millisecond,
	}
}
