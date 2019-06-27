/////////////////////////////////////////////////////////////////////////////////////////////////////////
//FILENAME: HrzService.go                                                                              //
//DESCRIPTION: A INTERFACES FOR CALLING THE HTTP SERVICE INTERFACES PROVIDED BY HORIZON ROBOTICS.      //
//             1.ADD PERSON AND IMAGE TO HORIZON SERVICE                                               //
//             2.DELETE PERSON AND IMAGE FROM HORIZON SERVICE                                          //
//             3.ADD SUBLIB TO HORIZON SERVICE                                                         //
//             4.DELETE SUBLIB FROM HORIZON SERVICE                                                    //
//             5.IMAGE RECOGNITION VIA HORIZON SERVICE                                                 //
//COMPANY: ARGUS                                                                                       //
//AUTHOR: JIEFENG LAI                                                                                  //
//DATE: Jun.20, 2018                                                                                   //
/////////////////////////////////////////////////////////////////////////////////////////////////////////
package restapi

import (
	"bytes"
	"github.com/astaxie/beego"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type HrzService struct {
	//Template operation URLs
	AddUrl string
	DelUrl string
	QueryUrl string
	ListImgUrl string

	//Sub Lib operation URLs
	AddLibUrl string
	DelLibUrl string
	QueryLibUrl string

	//Recognition URL
	RecogUrl string
}

var HRZService HrzService
func (v *HrzService) HrzServiceInit() {
	v.AddUrl = "http://127.0.0.1:9001/face_register"
	v.DelUrl = "http://127.0.0.1:9001/face_unregister"
	v.QueryUrl = "http://127.0.0.1:9001/person_list"
	v.ListImgUrl = "http://127.0.0.1:9001/image_list"
	v.AddLibUrl = "http://127.0.0.1:9001/category_create"
	v.DelLibUrl = "http://127.0.0.1:9001/category_delete"
	v.QueryLibUrl = "http://127.0.0.1:9001/category_list"
	v.RecogUrl = "http://127.0.0.1:9001/image_recognition"

	ip := beego.AppConfig.String("HrzServiceIP")
	port := beego.AppConfig.String("HrzServicePort")
	ipport := ip + ":" + port
	v.AddUrl = strings.Replace(v.AddUrl, "127.0.0.1:9001", ipport, -1)
	v.DelUrl = strings.Replace(v.DelUrl, "127.0.0.1:9001", ipport, -1)
	v.QueryUrl = strings.Replace(v.QueryUrl, "127.0.0.1:9001", ipport, -1)
	v.ListImgUrl = strings.Replace(v.ListImgUrl, "127.0.0.1:9001", ipport, -1)
	v.AddLibUrl = strings.Replace(v.AddLibUrl, "127.0.0.1:9001", ipport, -1)
	v.DelLibUrl = strings.Replace(v.DelLibUrl, "127.0.0.1:9001", ipport, -1) 
	v.QueryLibUrl = strings.Replace(v.QueryLibUrl, "127.0.0.1:9001", ipport, -1)
	v.RecogUrl = strings.Replace(v.RecogUrl, "127.0.0.1:9001", ipport, -1)
}

//Get Status Code
//0 - success;  <0 - failure
func (v *HrzService) GetStatusCode(strStatusMsg string) int {
	var pos = -1
	var ret = -1
	var err error
	token := strStatusMsg
	if pos = strings.Index(token, "\"is_success\""); pos < 0 {
		return ret
	}
	token = token[(pos + len("\"is_success\"") + 1):]
	if pos = strings.Index(token, "\n");     pos < 0 {
		return ret
	}

	token = token[:pos]
	token = strings.Replace(token, ",", "", -1)
	token = strings.Replace(token, ":", "", -1)
	token = strings.TrimSpace(token)
	if ret, err = strconv.Atoi(token); err != nil {
		ret = -2
		fmt.Println(err)
	}

	return ret
}

//Add Person
func (v *HrzService) Add(template_id string, base64img string, category string) error {
	if len(v.AddUrl) < 1 {
		v.HrzServiceInit()
	}

	var personPic PersonPic
	personPic.Index = 0
	personPic.Type = "jpg"
	personPic.Picdata = base64img

	var personInfo PersonInfo
	personInfo.Id = template_id
	personInfo.Category = category
	personInfo.Pictures = append(personInfo.Pictures, personPic)
	data, err := json.Marshal(personInfo)
	if err != nil {
		return fmt.Errorf(("json.Marshal: " + err.Error()))
	}

	var client http.Client
	var resp *http.Response
	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", v.AddUrl, body)
	if err != nil {
		return fmt.Errorf(("http.NewRequest: " + err.Error()))
	}

	//send request and get response
	if resp, err = client.Do(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	data, _ = ioutil.ReadAll(resp.Body)
	strStatusMsg := string(data)
	//fmt.Println(strStatusMsg)

	if nRet := v.GetStatusCode(strStatusMsg); nRet != 0 {
		if nRet == -6 {
			strStatusMsg = "This template has existed in the same sub lib."
		}
		return fmt.Errorf(("hrzService.Add failed: " + strStatusMsg))
	}
	return nil
}

func (v *HrzService) PostData(strUrl string, data []byte) (string, error) {
	var client http.Client
	var resp *http.Response
	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", strUrl, body)
	if err != nil {
		return "", fmt.Errorf(("json.Marshal: " + err.Error()))
	}

	//send request and get response
	if resp, err = client.Do(req); err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, _ = ioutil.ReadAll(resp.Body)

	strStatusMsg := string(data)
	//fmt.Println(strStatusMsg)
	return strStatusMsg, nil
}

func (v *HrzService) TemplateExisted (template_id string, category string) (bool, error) {
	if len(v.ListImgUrl) < 1 {
		v.HrzServiceInit()
	}

	var err error
	var data []byte
	var personDel PersonDelete
	personDel.Id = template_id
	personDel.Category = category
	data, err = json.Marshal(personDel)
	if err != nil {
		return false, fmt.Errorf(("json.Marshal: " + err.Error()))
	}

	var strStatusMsg = ""
	if strStatusMsg, err = v.PostData(v.ListImgUrl, data); err != nil {
		return false, err
	}

	var bRet = true
	nRet := v.GetStatusCode(strStatusMsg)
	if nRet != 0 { 
		bRet = false
	}
	return bRet, nil
}

//Delete template
func (v *HrzService) Delete(template_id string, category string) (error) {
	if len(v.DelUrl) < 1 {
		v.HrzServiceInit()
	}

	var err error
	var data []byte
	var personDel PersonDelete
	personDel.Id = template_id
	personDel.Category = category
	data, err = json.Marshal(personDel)
	if err != nil {
		return fmt.Errorf(("json.Marshal: " + err.Error()))
	}

	var client http.Client
	var resp *http.Response
	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", v.DelUrl, body)
	if err != nil {
		return fmt.Errorf(("json.Marshal: " + err.Error()))
	}
	//send request and get response
	if resp, err = client.Do(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	data, _ = ioutil.ReadAll(resp.Body)
	strStatusMsg := string(data)
	//fmt.Println("return val: " + strStatusMsg)

	if nRet := v.GetStatusCode(strStatusMsg); nRet != 0 {
		// "-5" means "User does not exist"
		if nRet != -5 {
			return fmt.Errorf(("hrzService.Delete failed: " + strStatusMsg))
		}
	}
	return nil
}
//AddLib
func (v *HrzService) AddLib(category string) error {
	return v.AddDeleteLib(v.AddLibUrl, category)
}
//Delete Lib
func (v *HrzService) DelLib(category string) error {
	return v.AddDeleteLib(v.DelLibUrl, category)
}

func (v *HrzService) QueryLib(strUrl string) error {
	return v.Query(strUrl)
}
func (v *HrzService) QueryPerson(strUrl string, category string) error {
	if len(category) > 0 {
		strUrl = strUrl + "?" + category
	}
	return v.Query(strUrl)
}
func (v *HrzService) Query(strUrl string) error {
	if len(v.QueryUrl) < 1 {
		v.HrzServiceInit()
	}

	var err error
	client := &http.Client{}
	var req *http.Request
	if req, err = http.NewRequest("GET", strUrl, nil); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return err
	}

	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	strStatusMsg := string(data)
	fmt.Println(strStatusMsg)
	return nil
}
func (v *HrzService) AddDeleteLib(strUrl string, category string) error {
	if len(v.AddLibUrl) < 1 {
		v.HrzServiceInit()
	}

	var err error
	var data []byte
	var lib LibCategory
	lib.Category = category
	data, err = json.Marshal(lib)
	if err != nil {
		return fmt.Errorf(("json.Marshal: " + err.Error()))
	}

	var client http.Client
	var resp *http.Response
	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", strUrl, body)
	if err != nil {
		return fmt.Errorf(("json.Marshal: " + err.Error()))
	}
	//send request and get response
	if resp, err = client.Do(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	data, _ = ioutil.ReadAll(resp.Body)
	strStatusMsg := string(data)
	//fmt.Println("return val: " + strStatusMsg)

	if nRet := v.GetStatusCode(strStatusMsg); nRet != 0 {
		// "-4" means "sub lib does not exist"
		if nRet != -4 {
			return fmt.Errorf(("hrzService.DeleteLib failed: " + strStatusMsg))
		}
	}
	return nil
}

//Image recognition
func (v *HrzService) RecogInHrz(base64img string, category string, topK int) (*RecogResp, error) {
	if len(v.RecogUrl) < 1 {
		v.HrzServiceInit()
	}

	var pic Picture
	pic.Type = "jpg"
	pic.Picdata = base64img

	var img Image4Recog
	img.DeviceId = "camera1"
	img.TrackId = "1000"
	//img.ReqType = 0
	img.Category = category
	//img.TopK = 1
	img.TopK = topK
	img.Pictures = append(img.Pictures, pic)
	data, err := json.Marshal(img)
	if err != nil {
		return nil, fmt.Errorf(("json.Marshal: " + err.Error()))
	}

	var client http.Client
	var resp *http.Response
	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", v.RecogUrl, body)
	if err != nil {
		return nil, fmt.Errorf(("http.NewRequest: " + err.Error()))
	}

	if resp, err = client.Do(req); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ = ioutil.ReadAll(resp.Body)
	//fmt.Println(string(data))

	var recResp RecogResp
	err = json.Unmarshal(data, &recResp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	//for i:=0; i<len(recResp.Details); i++ {
		//fmt.Println(recResp.Details[i])
	//}
	return &recResp, nil
}
