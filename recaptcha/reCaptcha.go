package recaptcha

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type conf struct {
	Addr   string `json:"addr"`
	Secret string `json:"secret"`
}

var cfg conf

func init() {
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Printf("ioutil.ReadFile err: %v", err)
	}

	if err = json.Unmarshal(b, &cfg); err != nil {
		log.Printf("json.Unmarshal err: %v", err)
	}
}

type reCaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func CheckRecaptcha(response string) (err error) {
	var u = make(url.Values)
	u.Add("secret", cfg.Secret)
	u.Add("response", response)

	var myReq = Request{
		Method: http.MethodPost,
		URL:    cfg.Addr + "?" + u.Encode(),
	}

	err = myReq.Send(nil)
	if err != nil {
		err = errors.New("myReq.Send error: " + err.Error())
		return
	}

	var result reCaptchaResponse
	err = myReq.json2Struct(&result)
	if err != nil {
		err = errors.New("myReq.json2Struct error: " + err.Error())
		return
	}

	// Check recaptcha verification success.
	if !result.Success {
		err = errors.New("unsuccessful recaptcha verify request")
		return
	}

	// Check response score.
	if result.Score < 0.5 {
		err = errors.New("lower received score than expected")
		return
	}

	return
}
