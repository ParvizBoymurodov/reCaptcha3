package recaptcha

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type siteVerifyRequest struct {
	RecaptchaResponse string `json:"g-recaptcha-response"`
}

func RecaptchaMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}

		var body siteVerifyRequest
		err = json.Unmarshal(bodyBytes, &body)
		if err != nil {
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		if err = CheckRecaptcha(body.RecaptchaResponse); err != nil {
			err = errors.New("Wrong Captcha")
			return
		}

		next.ServeHTTP(w, r)
	}
}
