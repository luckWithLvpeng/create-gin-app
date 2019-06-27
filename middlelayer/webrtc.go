package middlelayer

/*
#include "middlelayerc.h"
#include <stdlib.h>
*/
import "C"

import (
	"bytes"
	"eme/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	gosocketio "github.com/graarh/golang-socketio"

	"github.com/graarh/golang-socketio/transport"
	webrtc "github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"

	"github.com/astaxie/beego/logs"
)

//Message 客户端发送和接收的信息格式
//
//  From 客户端的socket.ID
//  To  盒子的Mac 地址
//  ChannelID 要请求的视频通道ID
//  Sdp  base64 编码的sdp
//  Type  消息的类型
//  Msg  当消息类型为错误的时候，附带的信息
//  Candidate  candidate 验证参数
//  SDPMid  candidate 验证参数
//  SDPMLineIndex  candidate 验证参数
//  UsernameFragment  candidate 验证参数
type Message struct {
	From             string `json:"from"`
	To               string `json:"to"`
	ChannelID        int    `json:"channelID"`
	Sdp              string `json:"sdp"`
	Type             string `json:"type"`
	Msg              string `json:"msg"`
	Candidate        string `json:"candidate"`
	SDPMid           string `json:"sdpMid"`
	SDPMLineIndex    uint16 `json:"sdpMLineIndex"`
	UsernameFragment string `json:"usernameFragment"`
}

var (
	client             *gosocketio.Client
	webURL             = gosocketio.GetUrl("39.105.67.236", 8080, false)
	websocketTransport = &transport.WebsocketTransport{
		PingInterval:   10 * time.Second,
		PingTimeout:    30 * time.Second,
		ReceiveTimeout: 30 * time.Second,
		SendTimeout:    30 * time.Second,
		BufferSize:     1024 * 32,
	}
	err    error
	pcs    = make(map[string]*webrtc.PeerConnection)
	config = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			// {
			// 	URLs: []string{"stun:39.105.67.236:3478"},
			// },
			// {
			// 	URLs:           []string{"turn:39.105.67.236:3478"},
			// 	Username:       "hsbox",
			// 	Credential:     "123",
			// 	CredentialType: webrtc.ICECredentialTypePassword,
			// },
		},
	}
	tracks         = make(map[int]*webrtc.Track)
	pipelinesLock  sync.Mutex
	mac            = getMacAddr()
	pcsLock        sync.Mutex
	ClientsMap     = make(map[int][]string) //key 为channel id  value 为客户端id字符串数组
	clientsMapLock sync.Mutex
	clientsStatus  = make(map[string]bool) //key 为clientid

)

const (
	videoClockRate = 90000
	audioClockRate = 48000
)

var BoxMac = mac

func connect() {
	client, err = gosocketio.Dial(webURL, websocketTransport)

	if err != nil {
		log.Println("链接", webURL, "失败, 2 秒后重连")
		// 连接失败后， 间隔一段时间后再次重连
		time.AfterFunc(time.Second*2, connect)
		return
	}
	client.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Println("Disconnected")
		connect()
	})
	client.On(gosocketio.OnError, func(err error) {
		log.Println("Webrtc gosocketio has err:", err)

	})
	client.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
		client.Emit("createOrJoin", mac)
	})

	client.On("log", func(h *gosocketio.Channel, args []string) {
		//log.Println("log from server:", strings.Join(args, " "))
	})
	client.On("created", func(h *gosocketio.Channel, room string) {
		log.Println("created room ", room)
		EnableChannelSourceCtl(mac)
		checkClientStatus()
		callAllClientsHasErrByChannel()

	})
	client.On("askToConnect", func(h *gosocketio.Channel, msg Message) {
		// 客户端请求建立连接
		if msg.ChannelID <= 0 {
			client.Emit("messageToBrowser", Message{
				Type:      "error",
				To:        msg.From,
				From:      msg.To,
				ChannelID: msg.ChannelID,
				Msg:       "need channelID for view",
			})
			return
		}
		clientsMapLock.Lock()
		if _, ok := ClientsMap[msg.ChannelID]; ok {
			ClientsMap[msg.ChannelID] = append(ClientsMap[msg.ChannelID], msg.From)
		} else {
			ClientsMap[msg.ChannelID] = []string{msg.From}
		}
		clientsMapLock.Unlock()

		logs.Debug(fmt.Sprintf("...Client ask from %v to view channel: %v...", msg.From, msg.ChannelID))
		if ClientsType[msg.ChannelID] == 0 {
			// logs.Debug(fmt.Sprintf("...System start new record at %v...",now))
			err := createPeerConnection(msg.From, msg.ChannelID)
			if err != nil {
				client.Emit("messageToBrowser", Message{
					Type:      "error",
					To:        msg.From,
					ChannelID: msg.ChannelID,
					From:      msg.To,
					Msg:       err.Error(),
				})
				return
			}
			client.Emit("messageToBrowser", Message{
				Type:      "ready",
				To:        msg.From,
				ChannelID: msg.ChannelID,
				From:      msg.To,
				Msg:       "{\"ForceRtmp\":0}",
			})
		} else {
			AddChannleChan <- msg.ChannelID
			select {
			case url := <-ChannelSourceStart:
				time.Sleep(5 * time.Second)
				client.Emit("messageToBrowser", Message{
					Type:      "ready",
					To:        msg.From,
					ChannelID: msg.ChannelID,
					From:      msg.To,
					Msg:       url,
				})
			case <-time.After(3 * time.Second):
				client.Emit("messageToBrowser", Message{
					Type:      "error",
					To:        msg.From,
					ChannelID: msg.ChannelID,
					From:      msg.To,
					Msg:       "Push stream to cloud timeout!",
				})

			}

		}

	})
	client.On("messageToBox", func(h *gosocketio.Channel, msg Message) {
		CIDStr := strconv.Itoa(msg.ChannelID)
		pcsLock.Lock()
		pc := pcs[msg.From+CIDStr]
		pcsLock.Unlock()
		if pc != nil {
			if msg.Type == "offer" {
				offer := webrtc.SessionDescription{}
				tmpbyte, err := base64.StdEncoding.DecodeString(msg.Sdp)
				defer func() {
					if e := recover(); e != nil {
						client.Emit("messageToBrowser", Message{
							Type:      "error",
							To:        msg.From,
							From:      msg.To,
							ChannelID: msg.ChannelID,
							Msg:       fmt.Sprintf("run time panic: %v", e),
						})
						return
					}
				}()
				if err != nil {
					client.Emit("messageToBrowser", Message{
						Type:      "error",
						To:        msg.From,
						ChannelID: msg.ChannelID,
						From:      msg.To,
						Msg:       err.Error(),
					})
					return
				}
				json.Unmarshal(tmpbyte, &offer)
				err = pc.SetRemoteDescription(offer)
				if err != nil {
					client.Emit("messageToBrowser", Message{
						Type:      "error",
						To:        msg.From,
						ChannelID: msg.ChannelID,
						From:      msg.To,
						Msg:       err.Error(),
					})
					return
				}
				answer, err := pc.CreateAnswer(nil)
				if err != nil {
					client.Emit("messageToBrowser", Message{
						Type:      "error",
						To:        msg.From,
						ChannelID: msg.ChannelID,
						From:      msg.To,
						Msg:       err.Error(),
					})
					return
				}
				err = pc.SetLocalDescription(answer)
				if err != nil {
					client.Emit("messageToBrowser", Message{
						Type:      "error",
						To:        msg.From,
						ChannelID: msg.ChannelID,
						From:      msg.To,
						Msg:       err.Error(),
					})
					return
				}
				tmpbyte, err = json.Marshal(answer)
				if err != nil {
					client.Emit("messageToBrowser", Message{
						Type:      "error",
						To:        msg.From,
						From:      msg.To,
						ChannelID: msg.ChannelID,
						Msg:       err.Error(),
					})
					return
				}
				client.Emit("messageToBrowser", Message{
					Type:      "answer",
					To:        msg.From,
					ChannelID: msg.ChannelID,
					From:      msg.To,
					Sdp:       base64.StdEncoding.EncodeToString(tmpbyte),
				})
			} else if msg.Type == "candidate" {
				pc.AddICECandidate(webrtc.ICECandidateInit{
					Candidate:        msg.Candidate,
					SDPMid:           &msg.SDPMid,
					SDPMLineIndex:    &msg.SDPMLineIndex,
					UsernameFragment: msg.UsernameFragment,
				})
			}
		}

		if msg.Type == "checkStatus" {
			k := msg.From + strconv.Itoa(msg.ChannelID)
			clientsStatus[k] = true
		} else if msg.Type == "clientClose" {
			fmt.Println("-------clientClose-----")
			fmt.Println(msg.From, msg.ChannelID)
			delClient(msg.From, msg.ChannelID)
		} else if msg.Type == "stunFailed" {
			logs.Error("-------stunFailed-----")
			delAllPeerConnection(msg.ChannelID)
			callAllClientsUseRtmp(msg.ChannelID)
		}

	})
}

func createPeerConnection(clientID string, channelID int) error {
	CIDStr := strconv.Itoa(channelID)

	// 假如一个客户端重复请求同一个视频，关闭历史建立的链接
	tmpPC := pcs[clientID+CIDStr]
	if tmpPC != nil {
		pcsLock.Lock()
		tmpPC.Close()
		delete(pcs, clientID+CIDStr)
		pcsLock.Unlock()
	}
	// 创建 pc
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		//打印 ice 状态的变化
		logs.Error(fmt.Sprintf("Connection State has changed %s \n", connectionState.String()))
		switch connectionState.String() {
		case "connected":
			AddChannleChan <- channelID
			<-ChannelSourceStart
		// case "failed":
		// delAllPeerConnection()
		// callAllClientsUseRtmp()
		case "disconnected":
			// if ForceRtmp == 0 {
			// 	delClient(clientID,channelID)
			// }else{
			peerConnection.Close()
			// }

		}

		// // 当远程pc 失去连接 ,关闭本地客户端，清除数据
		// if connectionState.String() == "disconnected" {
		// 	peerConnection.Close()
		// }
	})
	peerConnection.OnICECandidate(func(ICECandidate *webrtc.ICECandidate) {

		//client.Emit("messageToBrowser",Message{
		//	Type: "candidate",
		//	Candidate:        ICECandidate.Candidate,
		//	SDPMid:           ICECandidate.SDPMid,
		//	SDPMLineIndex:    ICECandidate.SDPMLineIndex,
		//	UsernameFragment: ICECandidate.UsernameFragment,
		//})
	})
	addStream(peerConnection, channelID, clientID)
	pcsLock.Lock()
	pcs[clientID+CIDStr] = peerConnection
	pcsLock.Unlock()
	return nil
}

func delClient(clientID string, channelID int) {
	logs.Error(fmt.Sprintf("...delClient  id is : %v,channel is %v...", clientID, channelID))
	DelChannleChan <- channelID
	CIDStr := strconv.Itoa(channelID)
	pcsLock.Lock()
	delClientFromClients(clientID, channelID)
	if _, ok := pcs[clientID+CIDStr]; ok {
		pcs[clientID+CIDStr].Close()
		delete(pcs, clientID+CIDStr)
	}
	pcsLock.Unlock()

}

func delClientFromClients(clientID string, channelID int) {
	i := -1
	for k, v := range ClientsMap[channelID] {
		if v == clientID {
			i = k
			break
		}
	}
	if i != -1 {
		clientsMapLock.Lock()
		ClientsMap[channelID] = append(ClientsMap[channelID][:i], ClientsMap[channelID][i+1:]...)
		clientsMapLock.Unlock()
	}
}
func delAllPeerConnection(channelID int) {
	CIDStr := strconv.Itoa(channelID)

	delete(tracks, channelID)
	for k, peerConnection := range pcs {
		if result := strings.HasSuffix(k, CIDStr); result {
			pcsLock.Lock()
			peerConnection.Close()
			delete(pcs, k)
			pcsLock.Unlock()
		}

	}

}

func addStream(peerConnection *webrtc.PeerConnection, channelID int, clientID string) {
	// 在这里添加视频流
	if tracks[channelID] != nil {
		peerConnection.AddTrack(tracks[channelID])
		return
	}
	// 没有则先创建这个通道视频
	VideoTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeVP8, rand.Uint32(), mac+strconv.Itoa(channelID), mac+strconv.Itoa(channelID))
	tracks[channelID] = VideoTrack
	if err != nil {
		sendErrorToClient(err, channelID, clientID)
	}
	_, err = peerConnection.AddTrack(VideoTrack)
	if err != nil {
		sendErrorToClient(err, channelID, clientID)
	}
}

func callAllClientsHasErrByChannel() {
	for {
		errChannelMap := <-ChannelErrorChan
		for channelID, err := range errChannelMap {
			for _, clientID := range ClientsMap[channelID] {
				sendErrorToClient(err, channelID, clientID)
			}
		}
	}
}

func sendErrorToClient(err error, channelID int, clientID string) {
	client.Emit("messageToBrowser", Message{
		Type:      "error",
		To:        clientID,
		ChannelID: channelID,
		From:      mac,
		Msg:       err.Error(),
	})
}
func callAllClientsUseRtmp(channelID int) {
	logs.Debug(fmt.Sprintf("...callAllClientsUseRtmp channel id is %v: ...", channelID))
	ForceRtmpChan <- channelID
	playUrl := <-ChangeRtmpChan
	time.Sleep(5 * time.Second)
	for _, clientID := range ClientsMap[channelID] {
		client.Emit("messageToBrowser", Message{
			Type:      "ready",
			To:        clientID,
			ChannelID: channelID,
			From:      mac,
			Msg:       playUrl,
		})
	}
}
func checkClientStatus() {
	go func() {
		for {
			for channelID, cs := range ClientsMap {
				for _, clientID := range cs {
					clientsStatus[clientID+strconv.Itoa(channelID)] = false
					client.Emit("messageToBrowser", Message{
						Type:      "checkStatus",
						To:        clientID,
						ChannelID: channelID,
						From:      mac,
					})

				}
			}
			time.Sleep(10 * time.Second)
			for channelID, cs := range ClientsMap {
				for _, clientID := range cs {
					if s, ok := clientsStatus[clientID+strconv.Itoa(channelID)]; ok && !s {
						delClient(clientID, channelID)
					}
				}
			}

		}
	}()

}

//export HandlePipelineBufferFromGo
func HandlePipelineBufferFromGo(buffer unsafe.Pointer, bufferLen C.int, duration C.int, channelID C.int) {
	pipelinesLock.Lock()
	track, ok := tracks[int(channelID)]
	pipelinesLock.Unlock()
	if ok {
		var samples uint32
		samples = uint32(videoClockRate * (float32(duration) / 1000000000))
		if err := track.WriteSample(media.Sample{Data: C.GoBytes(buffer, bufferLen), Samples: samples}); err != nil {
			log.Println("Webrtc track writeSample has err:", err)
			delete(tracks, int(channelID))
		}
	} else {
		//fmt.Printf("discarding buffer, no pipeline with id %d \n", int(channelID))
	}
	//	C.free(buffer)
}

// getMacAddr gets the MAC hardware
// address of the host machine
func getMacAddr() string {
	var address string
	inter, err := net.InterfaceByName("eth0")
	if err != nil {
		if interfaces, err := net.Interfaces(); err == nil {
			for _, interf := range interfaces {
				if (interf.Flags&net.FlagUp) != 0 && bytes.Compare(interf.HardwareAddr, nil) != 0 {
					if address = interf.HardwareAddr.String(); len(address) < 1 {
						continue
					}
					address = strings.Replace(address, ":", "", -1)
					break
				}
			}
		}
		fmt.Println("not eth0 mac address:", address)
		return address
	}
	//mac地址
	address = inter.HardwareAddr.String()
	address = strings.Replace(address, ":", "", -1)
	return address
}

// InitWebtrc 初始化 webrtc
func InitWebtrc() {
	if models.EngineConfig.Webrtc {
		connect()
	}
}
