package restapi

//Person (un)register/////////////////////////////
type PersonPic struct {
	Index int `json:"index"` //index 0
	Picdata string `json:"picdata"`  //base64 string
	Type string `json:"type"`  //jpg, png, bmp
}
type PersonInfo struct {
	Id string `json:"_id"`
	Category string `json:"category"`
	Pictures []PersonPic `json:"pictures"`
}
type PersonDelete struct {
	Id string `json:"_id"`
	Category string `json:"category"`
}

//Sub lib category
type LibCategory struct {
	Category string `json:"category"`
}

//Image Recognition/////////////////////////////
type FaceRect struct {
	Left int `json:"left"`
	Top int  `json:"top"`
	Right int  `json:"right"`
	Bottom int  `json:"bottom"`
}
type Picture struct {
	Type string `json:"type"`  //jpg, png, bmp, yuv(NV21)
	Width int `json:"width"`
	Height int `json:"height"`
	Picdata string `json:"picdata"`  //base64 string
	////Face_rect FaceRect `json:"face_rect"`
}
type Image4Recog struct {
	DeviceId string `json:"device_id"`
	TrackId string `json:"track_id"`
	Category string `json:"category"`
	ReqType int `json:"req_type"`
	Pictures []Picture `json:"pictures"`
	TopK int `json:"top_k"`
}
////////////////////////////////////////////////
type Result struct {
	Id string `json:"_id"`
	Dist float32 `json:"dist"`
	Score float32 `json:"score"`
}
type RecogResp struct {
	Details []Result `json:"details"`
	IsRecord int `json:"is_record"`
	IsSuccess int `json:"is_success"`
	QryImgs  []string `json:"qry_imgs"`
	Result int `json:"result"`
}
////////////////////////////////////////////////
