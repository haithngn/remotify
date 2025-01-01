package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"notification-deployer/internal/data"
	"time"
)

type Pairing struct {
	maxAttempt int
	duration   int
}

func NewPairing(maxAttempt, duration int) *Pairing {
	return &Pairing{
		maxAttempt: maxAttempt,
		duration:   duration,
	}
}

type PairResult struct {
	DeviceType  string `json:"device_type"`
	DeviceToken string `json:"device_token"`
}

type SucceedCallback func(pairing *Pairing, pairId string, r PairResult, err error)
type ObserveFunc func(pairId string) (string, error)

func (p *Pairing) StartObserveDeviceTokenParing(pairId string, observeFunc ObserveFunc, callback SucceedCallback) {
	p.observeMobileDeviceToken(pairId, 3, 3, observeFunc, callback)
}

func (p *Pairing) observeMobileDeviceToken(pairId string, attempts int, duration int, observeFunc ObserveFunc, callback SucceedCallback) {
	if attempts < 0 {
		callback(p, pairId, PairResult{}, errors.New("pairing expired"))

		return
	}

	var res PairResult
	resp, err := observeFunc(pairId)
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		log.Error("Cannot unmarshal response")
	} else {
		callback(p, pairId, res, nil)
		return
	}

	time.Sleep(time.Duration(duration) * time.Second)
	p.observeMobileDeviceToken(pairId, attempts-1, duration, observeFunc, callback)
}

func (p *Pairing) GetHTTP(endpoint, header, params string) string {
	var parameters url.Values
	err := json.Unmarshal([]byte(params), &parameters)
	if err != nil {
		return fmt.Sprintf(`{
				"status_code": %s,
				"body": {
					"message": "Bad Request",
				}
			}`, data.AppErrorCode_Invalid_Header)
	}

	u := "https://dummyjson.com" + endpoint + parameters.Encode()
	log.Debug("HTTP GET: ", u)

	// Make the GET request
	resp, err := http.Get(u)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to post HTTP (%s)", err.Error()))
		return fmt.Sprintf(`{
				"status_code": %s,
				"body": {
					"message": "Failed to post HTTP",
				}
			}`, data.AppErrorCode_Cannot_Make)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Failed to close response responseData")
		}
	}(resp.Body)

	log.Debug("Response Status:", resp.Status)
	log.Debug("Response Headers:", resp.Header)
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response responseData")
		return fmt.Sprintf(`{
				"status_code": %s,
				"body": {
					"message": "Failed to read response responseData",
				}
			}`, data.AppErrorCode_Invalid_Response)
	}

	return fmt.Sprintf(`{
				"status_code": 200,
				"body": %s
			}`, string(responseData))
}
