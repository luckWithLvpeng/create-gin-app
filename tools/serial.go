package tools

import (
	//"container/list"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/astaxie/beego"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"strconv"
	//"net"
	"container/list"
	"sync"
	"time"

	"github.com/tarm/serial"
)

const (
	MessageLen = 4
	Parity     = serial.ParityNone // 奇校验
	StopBit    = serial.Stop1      // 停止位 1
	DataNum    = 8                 // 数据位 8

)

var RelayId int = 0

var (
	SerialControlTime, err  = beego.AppConfig.Int("SerialControlTime")
	SerialCloseSpaceTime, _ = beego.AppConfig.Int("SerialCloseSpaceTime")
	AllSerialOpenFlag, _    = beego.AppConfig.Bool("SerialOpenFlag")
	SerialAddr              = beego.AppConfig.String("SerialAddr")
	SerialPort              = beego.AppConfig.String("SerialPort")
	SerialBaud, _           = beego.AppConfig.Int("SerialBaud")
	SerialType, _           = beego.AppConfig.Int("SerialType")
	SerialNetAddr           = beego.AppConfig.String("SerialNetAddr")
	//RelayId, _             = beego.AppConfig.Int("SerialRelayId")
	//SerialNotHitOpen, _    = beego.AppConfig.Bool("SerialNotHitOpen")
	SerialOpen  = []byte{0xa0, 0x01, 0x01, 0xa2}
	SerialClose = []byte{0xa0, 0x01, 0x00, 0xa1}
)

type SerialClient struct {
	Client  *serial.Port
	NetConn net.Conn

	WriteCh    chan<- []byte
	SerialType int    //网络继电器时用 type 1 usb;2 net
	SerialIp   string // 网络继电器时用
	SerialPort string
	SerialBaud int
	SerialLock sync.Mutex
}

type SerialInfo struct {
	NumId      int
	SerialType int    //继电器类型 type 1 usb; 2 网络
	SerialAddr string // 网络继电器时用
	SerialPort int
}

type SerialNum struct {
	Id                int
	Serial_id         int
	SerialOpen        []byte
	SerialClose       []byte
	Chanel_id         []int
	SerialLock        sync.Mutex
	Time              time.Time
	SerialOpenFlag    bool
	SerialNotHitOpen  bool
	SerialHitOpen     bool
	SerialVicitorOpen bool
	LogicdbIds        []int
	SerialInfo        *SerialInfo
	Client            *serial.Port
	NetConn           net.Conn

	WriteCh    chan<- []byte
	PersonList *list.List
}

var serialClientUSB *SerialClient
var PersonList *list.List
var SerialMap = make(map[int]*SerialNum)
var SerialMapLock sync.Mutex
var UsbClient *serial.Port

func init() {
	if !AllSerialOpenFlag {
		return
	}
	PersonList = list.New()
	serialClient, _  := NewSerialClient(1, SerialNetAddr, SerialPort, SerialBaud)
	serialClientUSB = serialClient

	go Open()
	Subscribe(SerialCallback)
}

func SerialUpdateStatus(serial SerialNum) {
	SerialMapLock.Lock()
	if v, ok := SerialMap[serial.Id]; ok {
		v.SerialHitOpen = serial.SerialHitOpen
		v.SerialNotHitOpen = serial.SerialNotHitOpen
		v.SerialOpenFlag = serial.SerialOpenFlag
		v.SerialVicitorOpen = serial.SerialVicitorOpen
		v.Chanel_id = serial.Chanel_id
		v.LogicdbIds = serial.LogicdbIds
		v.SerialInfo = serial.SerialInfo

		if serial.SerialInfo.SerialType == 1 {
			a := 161 + serial.SerialInfo.NumId
			var open []byte
			open = append(open,0xa0)
			open = append(open,byte(serial.SerialInfo.NumId))
			open = append(open,0x01)
			open = append(open,byte(a))
			v.SerialOpen = open
			
			b := 160 + serial.SerialInfo.NumId
			var closed []byte
			closed = append(closed,0xa0)
			closed = append(closed,byte(serial.SerialInfo.NumId))
			closed = append(closed,0x00)
			closed = append(closed,byte(b))
			v.SerialClose = closed		
			//SerialMap[serial.Id] = &serial
		}
		if serial.SerialInfo.SerialType == 2 {
			SerialConnection(&serial)
			//< 11 开1 12 开2， 可后边加 :妙数 自动释放
			second := SerialCloseSpaceTime / 1000
			if second < 1 {
				second = 1
			}
			var SerialOpen string = fmt.Sprintf("%d:%d", 10+serial.SerialInfo.NumId, second)
			v.SerialOpen = []byte(SerialOpen)
			fmt.Println("net serial serial open string", SerialOpen)
			///< 21 关1; 22 关2
			var SerialClose string = fmt.Sprintf("%d", 20+serial.SerialInfo.NumId)
			v.SerialClose = []byte(SerialClose)
			//SerialMap[serial.Id] = &serial
		}
		fmt.Println("serial ******", v)
	} else {
		serial.PersonList = list.New()
		if serial.SerialInfo.SerialType == 1 {

			a := 161 + serial.SerialInfo.NumId
			serial.SerialOpen = append(serial.SerialOpen, 0xa0)
			serial.SerialOpen = append(serial.SerialOpen, byte(serial.SerialInfo.NumId))
			serial.SerialOpen = append(serial.SerialOpen, 0x01)
			serial.SerialOpen = append(serial.SerialOpen, byte(a))
			b := 160 + serial.SerialInfo.NumId
			serial.SerialClose = append(serial.SerialClose, 0xa0)
			serial.SerialClose = append(serial.SerialClose, byte(serial.SerialInfo.NumId))
			serial.SerialClose = append(serial.SerialClose, 0x00)
			serial.SerialClose = append(serial.SerialClose, byte(b))
			SerialMap[serial.Id] = &serial
		}
		if serial.SerialInfo.SerialType == 2 {
			SerialConnection(&serial)
			//< 11 开1 12 开2， 可后边加 :妙数 自动释放
			second := SerialCloseSpaceTime / 1000
			if second < 1 {
				second = 1
			}
			var SerialOpen string = fmt.Sprintf("%d:%d", 10+serial.SerialInfo.NumId, second)
			serial.SerialOpen = []byte(SerialOpen)
			fmt.Println("net serial serial open string", SerialOpen)
			///< 21 关1; 22 关2
			var SerialClose string = fmt.Sprintf("%d", 20+serial.SerialInfo.NumId)
			serial.SerialClose = []byte(SerialClose)
			SerialMap[serial.Id] = &serial
		}

	}
	SerialMapLock.Unlock()
	return
}

func SerialCallback(client MQTT.Client, message MQTT.Message) {

	if !AllSerialOpenFlag {
		return
	}
	beego.Info("**enter serial mqttcallback**")
	var matchResult MatchResultMQTT
	err := json.Unmarshal(message.Payload(), &matchResult)
	if err != nil {
		beego.Error(err)
		return
	}
	///< open
	matchResult.Time = time.Now()
	serials := GetMatchSerial(matchResult.Channel_id)
	if len(serials) <= 0 {
		return
	}

	for _, serial := range serials {
		if !serial.SerialOpenFlag {
			continue
		}

		serial.Time = time.Now()
		if matchResult.Hit_flag {
			if serial.SerialHitOpen {
				/*if !serial.SerialVicitorOpen && matchResult.Sublib_id == 2 {
					continue
				}*/
				///改为按选择分库开门
				for i := 0; i < len(serial.LogicdbIds); i++ {
					if serial.LogicdbIds[i] == matchResult.Sublib_id {
						PersonList.PushBack(serial)
						continue
					}
				}

			}

		} else {
			if serial.SerialNotHitOpen {
				PersonList.PushBack(serial)
			}
		}

	}

}

func NewSerialClient(serialType int, SerialIp string, SerialPort string, SerialBaunt int) (*SerialClient, error) {
	var client = new(SerialClient)
	client.SerialType = serialType
	client.SerialIp = SerialIp
	client.SerialPort = SerialPort
	client.SerialBaud = SerialBaunt
	var err error

	err = client.SerialConnection()
	if err != nil {
		//beego.Error(fmt.Sprintf("serial connection failed, error is %s\n"), err.Error())
		return client, err
	}
	beego.Info("**serial connection success**")
	return client, nil
}

func SerialConnection(serialClient *SerialNum) error {
	var err error
	if serialClient.SerialInfo.SerialType == 1 {

		serialClient.Client, err = serial.OpenPort(&serial.Config{
			Name:        serialClient.SerialInfo.SerialAddr,
			Baud:        serialClient.SerialInfo.SerialPort,
			Parity:      Parity,
			StopBits:    StopBit,
			Size:        DataNum,
			ReadTimeout: time.Second * 5,
		})
	}
	if serialClient.SerialInfo.SerialType == 2 {
		serialClient.NetConn, err = net.DialTimeout("tcp",
			serialClient.SerialInfo.SerialAddr+":"+strconv.Itoa(serialClient.SerialInfo.SerialPort),
			time.Second*5)
		if err != nil {
			fmt.Println("Error connecting:", err)
			serialClient.NetConn = nil
			return err
		}
		beego.Info("Connecting to err"+serialClient.SerialInfo.SerialAddr+":",
			serialClient.SerialInfo.SerialPort)
		return nil
	}

	if err != nil {
		beego.Info("Connecting to err"+serialClient.SerialInfo.SerialAddr+":",
			serialClient.SerialInfo.SerialPort, err)
		return err
	}

	fmt.Println("Connecting to "+serialClient.SerialInfo.SerialAddr+":",
		serialClient.SerialInfo.SerialPort)
	return nil
}

func (sc *SerialClient) SerialConnection() error {
	var err error
	if sc.SerialType == 1 {
		sc.Client, err = serial.OpenPort(&serial.Config{
			Name:        sc.SerialPort,
			Baud:        sc.SerialBaud,
			Parity:      Parity,
			StopBits:    StopBit,
			Size:        DataNum,
			ReadTimeout: time.Second * 5,
		})
	}
	if sc.SerialType == 2 {
		sc.NetConn, err = net.DialTimeout("tcp", sc.SerialIp, time.Second*5)
		if err != nil {
			fmt.Println("Error connecting:", err)
			sc.NetConn = nil
			return err
		}
		beego.Info("Connecting to err"+sc.SerialIp+":", sc.SerialBaud)
		return nil
	}

	if err != nil {
		beego.Info("Connecting to err"+sc.SerialPort+":", sc.SerialBaud, err)
		return err
	}

	fmt.Println("Connecting to "+sc.SerialPort+":", sc.SerialBaud)
	return nil
}
func GetMatchSerial(channel_id int) (serials []*SerialNum) {
	for _, v := range SerialMap {
		for i := 0; i < len(v.Chanel_id); i++ {
			if v.Chanel_id[i] == channel_id {
				serials = append(serials, v)
				//beego.Info("*****match open serials status*****", *v)
			}
		}
	}
	return

}

func (sc *SerialClient) OpenSerial(serial *SerialNum) (int, error) {
	var len int
	var err error
	if sc.SerialType == 1 {
		if sc.Client == nil {
			//sc.Client.Close()
			err = sc.SerialConnection()
			if err != nil {
				return 0, err
			}
		}
		//sc.SerialLock.Lock()
		//defer sc.SerialLock.Unlock()
		len, err = sc.Client.Write(serial.SerialOpen)
		if err != nil || len != MessageLen {
			beego.Info("Error to open serial because of ", err.Error())
			err = sc.SerialConnection()
			if err != nil {
				return 0, err
			}
			sc.Client.Write(serial.SerialOpen)
		}

		if SerialCloseSpaceTime < 100 {
			SerialCloseSpaceTime = 100
		}

		time.Sleep(time.Millisecond * time.Duration(SerialCloseSpaceTime))

		len, err = sc.Client.Write(serial.SerialClose)
		if err != nil || len != MessageLen {
			beego.Info("Error to close serial because of ", err.Error())
			err = sc.SerialConnection()
			if err != nil {
				return 0, err
			}
			sc.Client.Write(serial.SerialClose)
		}

		return len, nil
	}
	if sc.SerialType == 2 {

		sc.SerialConnection()
		if sc.NetConn == nil {
			fmt.Println("netConn is nil")
			err = sc.SerialConnection()
			if err != nil {
				fmt.Println("serial connection error", err.Error())
				return 0, err
			}
		}
		defer sc.NetConn.Close()
		fmt.Println("open", serial.SerialOpen)
		len, err = sc.NetConn.Write(serial.SerialOpen)
		if err != nil || len < 4 {
			if err != nil {
				fmt.Println("open error", err.Error())
			}
			sc.NetConn.Close()
			err = sc.SerialConnection()
			if err == nil {
				len, err = sc.NetConn.Write(serial.SerialOpen)
			}
			return 0, err

		}

		time.Sleep(time.Millisecond * 200)
		return len, err

	}

	return len, err

}

func OpenSerial(serialClient *SerialNum) (int, error) {
	var len int
	var err error
	if serialClient.SerialInfo.SerialType == 1 {
		if serialClientUSB.Client == nil {
			//sc.Client.Close()
			err = serialClientUSB.SerialConnection()
			if err != nil {
				return 0, err
			}
		}
		//sc.SerialLock.Lock()
		//defer sc.SerialLock.Unlock()
		len, err = serialClientUSB.Client.Write(serialClient.SerialOpen)
		if err != nil || len != MessageLen {
			if err != nil {
				beego.Info("Error to open serial because of ", err.Error())
			}

			err = serialClientUSB.SerialConnection()
			if err != nil {
				return 0, err
			}
			serialClientUSB.Client.Write(serialClient.SerialOpen)
		}

		if SerialCloseSpaceTime < 100 {
			SerialCloseSpaceTime = 100
		}

		time.Sleep(time.Millisecond * time.Duration(SerialCloseSpaceTime))

		len, err = serialClientUSB.Client.Write(serialClient.SerialClose)
		if err != nil || len != MessageLen {
			if err != nil {
				beego.Info("Error to open serial because of ", err.Error())
			}

			err = serialClientUSB.SerialConnection()
			if err != nil {
				return 0, err
			}
			serialClientUSB.Client.Write(serialClient.SerialClose)
		}

		return len, nil
	}
	if serialClient.SerialInfo.SerialType == 2 {

		SerialConnection(serialClient)
		if serialClient.NetConn == nil {
			fmt.Println("netConn is nil")
			err = SerialConnection(serialClient)
			if err != nil {
				fmt.Println("serial connection error", err.Error())
				return 0, err
			}
		}
		defer serialClient.NetConn.Close()
		fmt.Println("open", serialClient.SerialOpen)
		len, err = serialClient.NetConn.Write(serialClient.SerialOpen)
		if err != nil || len < 4 {
			if err != nil {
				fmt.Println("open error", err.Error())
			}
			serialClient.NetConn.Close()
			err = SerialConnection(serialClient)
			if err == nil {
				len, err = serialClient.NetConn.Write(serialClient.SerialOpen)
			}
			return 0, err

		}

		time.Sleep(time.Millisecond * 200)
		return len, err

	}

	return len, err

}

func dequeue() interface{} {
	if PersonList.Len() == 0 {
		return nil
	}
	e := PersonList.Front()
	result := PersonList.Remove(e)
	return result
}

func IsTimeout(ctreatTime time.Time) bool {
	if SerialControlTime < 1000 {
		SerialControlTime = 1000
	}
	if time.Now().Second()-ctreatTime.Second() > SerialControlTime/1000 {
		//fmt.Println(time.Now().Second() - ctreatTime.Second())
		return true
	}

	return false
}

func Open() {
	for {
		if PersonList.Len() == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		} else {
			ele := dequeue().(*SerialNum)
			if ele == nil {
				continue
			}

			//if IsTimeout(ele.Time) {
			//fmt.Println("timeout", PersonList.Len())
			//continue
			//} else {
			fmt.Println("open channel:", ele.Chanel_id, "serial_id", ele.Id)
			OpenSerial(ele)
			continue
			//}
		}
		return
	}
}
