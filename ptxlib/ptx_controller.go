package ptxlib

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//PTXController PTX Controller
type PTXController struct {
	appID  string
	appKey string
}

// NewPTXController 新建一個 PTXController
func NewPTXController(id, key string) (p *PTXController) {
	p = &PTXController{}
	p.appID = id
	p.appKey = key
	return
}

// APIGet 取的特定API的資料
func (p *PTXController) APIGet(url string) (data []byte, err error) {
	xdate, auth := p.authGenerator()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("x-date", xdate)
	req.Header.Set("Authorization", auth)
	rsp, err := client.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	data, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	return
}

// JSONGet 取的特定API的資料，並將回傳資料轉換
func (p *PTXController) JSONGet(apiURL string, queryParams url.Values, convertSturct interface{}) (err error) {
	if queryParams == nil {
		queryParams = url.Values{}
	}
	queryParams.Set("$format", "json")

	resBytes, err := p.APIGet(apiURL + "?" + queryParams.Encode())
	if err != nil {
		return
	}

	err = json.Unmarshal(resBytes, &convertSturct)
	return
}

// AuthGenerator 生成 ptx hmac 認證用的 header (Authorization, x-date)
func (p *PTXController) authGenerator() (string, string) {
	xdate, sign := signGenerator(p.appKey)
	auth := "hmac username=\"" + p.appID + "\", algorithm=\"hmac-sha1\", headers=\"x-date\", signature=\"" + sign + "\""

	return xdate, auth
}

func signGenerator(APPKEY string) (string, string) {
	xdate := getServerTime()
	encryptXdate := "x-date: " + xdate
	encryptSign := hmacSha1Generator(encryptXdate, APPKEY)

	return xdate, encryptSign
}

// getServerTime 獲取現在時間，回傳 ptx hmac 認證所使用的時間格式
func getServerTime() string {
	//ptx platform time is GMT 0.
	return time.Now().UTC().Format(http.TimeFormat)
}

func hmacSha1Generator(encXdate, appkey string) string {
	key := []byte(appkey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(encXdate))
	macEncrypted := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return macEncrypted
}
