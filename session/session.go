package session

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Session struct {
	headers map[string]string
	verify  bool
	client  http.Client
}

type Options func(session *Session)

func WithHeaders(headers map[string]string) Options {
	return func(session *Session) {
		session.headers = headers
	}
}

func WithVerify(verify bool) Options {
	return func(session *Session) {
		session.verify = verify
	}
}

func New(opts ...Options) *Session {
	session := &Session{headers: nil, verify: false, client: http.Client{}}
	for _, opt := range opts {
		opt(session)
	}
	if session.verify == false {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		session.client = http.Client{Transport: tr}
	} else {
		session.client = http.Client{}
	}
	if session.headers == nil {
		session.headers = make(map[string]string)
	}
	return session
}

func (s *Session) resetHeaders(res http.Response) {
	s.headers["Cookie"] = res.Header.Get("Cookie")
}

func (s *Session) putHeaders(req *http.Request) {
	for key_, value_ := range s.headers {
		req.Header.Set(key_, value_)
	}
}

func (s *Session) Headers() map[string]string {
	return s.headers
}

func (s *Session) AddHeader(key string, value string) {
	s.headers[key] = value
}

func (s *Session) templateRequest(url string, data map[string]interface{}, method string) *Handler {
	reqBody, err := json.Marshal(data)
	if err != nil {
		return &Handler{json: ""}
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return &Handler{json: ""}
	}
	s.putHeaders(req)
	res, err := s.client.Do(req)
	if err != nil {
		return &Handler{json: ""}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return &Handler{json: string(body)}
}

func (s *Session) Get(url string, params map[string]interface{}) *Handler {
	res := s.templateRequest(url, params, "GET")
	return res
}

func (s *Session) Post(url string, data map[string]interface{}) *Handler {
	res := s.templateRequest(url, data, "POST")
	return res
}

func (s *Session) UploadFile(url string, field string, filePath string, data map[string]string) *Handler {
	// 如果读取错误直接将错误返回
	file, err := os.Open(filePath)
	if err != nil {
		return &Handler{json: ""}
	}
	defer file.Close()
	// 创建写入
	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile(field, filePath)
	if err != nil {
		return &Handler{json: ""}
	}
	// 拷贝文件内容
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return &Handler{json: ""}
	}
	// 写入其他字段
	for key_, val := range data {
		_ = bodyWriter.WriteField(key_, val)
	}
	err = bodyWriter.Close()
	contentType := bodyWriter.FormDataContentType()
	req, _ := http.NewRequest("POST", url, bodyBuf)
	s.headers["Content-Type"] = contentType
	s.putHeaders(req)
	res, err := s.client.Do(req)
	if err != nil {
		return &Handler{json: ""}
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return &Handler{json: string(body)}
}

func (s *Session) Put(url string, data map[string]interface{}) *Handler {
	res := s.templateRequest(url, data, "PUT")
	return res
}

func (s *Session) Patch(url string, data map[string]interface{}) *Handler {
	res := s.templateRequest(url, data, "PATCH")
	return res
}

func (s *Session) Delete(url string, data map[string]interface{}) *Handler {
	res := s.templateRequest(url, data, "DELETE")
	return res
}
