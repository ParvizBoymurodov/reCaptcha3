package recaptcha

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var netClient = &http.Client{
	Timeout: time.Hour * 40,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// Request model of http client request
type Request struct {
	Method         string
	URL            string
	Headers        map[string]string
	PostData       []byte
	RespData       []byte
	RespStatusCode int
	Username       string
	Password       string
}

// Send - sends a request;
func (request *Request) Send(client *http.Client) (err error) {

	if client == nil {
		client = netClient
	}

	var bodyReader io.Reader
	if request.Method == http.MethodPost && len(request.PostData) > 0 {
		bodyReader = bytes.NewReader(request.PostData)
	}

	// prepare request;
	httpReqs, err := http.NewRequest(request.Method, request.URL, bodyReader)
	if err != nil {
		err = errors.New("http.NewRequest error: " + err.Error())
		return
	}
	if len(request.Username) != 0 && len(request.Password) != 0 {
		httpReqs.SetBasicAuth(request.Username, request.Password)
	}
	if len(request.Headers) != 0 {
		for key, value := range request.Headers {
			httpReqs.Header.Set(key, value)
		}
	}
	// send request;
	httpResp, err := client.Do(httpReqs)
	if err != nil {
		err = errors.New("client.Do error: " + err.Error())
		return
	}
	defer httpResp.Body.Close()
	request.RespStatusCode = httpResp.StatusCode

	if httpResp.Body == nil {
		err = errors.New("empty_http_response")
		return
	}

	// read response;
	request.RespData, err = ioutil.ReadAll(httpResp.Body)
	if err != nil {
		err = errors.New("ioutil.ReadAll error: " + err.Error())
		return
	}

	return
}

func (r *Request) json2Byte(v interface{}) (err error) {
	r.PostData, err = json.Marshal(v)
	return
}

func (r Request) json2Struct(v interface{}) (err error)  {
	err = json.Unmarshal(r.RespData, v)
	return
}