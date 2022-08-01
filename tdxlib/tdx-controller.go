package tdxlib

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// updateBufferDuration token 要失效前的更新緩衝時間
const updateBufferDuration = 1 * time.Hour

// authFailedRetryInterval token 更新失敗時，重新嘗試的間隔時間
const authFailedRetryInterval = 3 * time.Minute

// apiBasicBaseURL tdx 基礎服務 API 的 Base URL
const apiBasicBaseURL string = "https://tdx.transportdata.tw/api/basic"

//TDXController TDX API Controller
type TDXController struct {
	sync.RWMutex
	cID           string // Client Id
	cSEC          string // Client Secret
	authorization string // authorization header value
	refreshTimer  *time.Timer
}

// NewTDXController 新建一個 TDXController
func NewTDXController(cID, cSEC string) (resTC *TDXController, err error) {
	tc := &TDXController{
		cID:  cID,
		cSEC: cSEC,
	}

	var refreshAfter time.Duration
	refreshAfter, err = tc.getAccessToken()
	if err != nil {
		err = errors.New("init auth failed, err= " + err.Error())
		return
	}
	tc.setUpdateTimer(refreshAfter)

	resTC = tc
	return
}

// GetAuthorization 取得 authorization
func (tc *TDXController) getAuthorization() (authorization string) {
	tc.RLock()
	authorization = tc.authorization
	tc.RUnlock()
	return
}

// APIGet 取的特定API的資料
func (tc *TDXController) APIGet(url string) (respBody io.ReadCloser, err error) {
	authorization := tc.getAuthorization()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", authorization)
	rsp, err := client.Do(req)
	if err != nil {
		return
	}

	respBody = rsp.Body
	return
}

// JSONGet 取的特定API的資料，並將回傳資料轉換
func (tc *TDXController) JSONGet(apiURL string, queryParams url.Values, convertSturct interface{}) (err error) {
	if queryParams == nil {
		queryParams = url.Values{}
	}
	queryParams.Set("$format", "json")

	respBody, err := tc.APIGet(apiURL + "?" + queryParams.Encode())
	if err != nil {
		return
	}
	defer respBody.Close()

	jDecoder := json.NewDecoder(respBody)
	err = jDecoder.Decode(&convertSturct)
	return
}
