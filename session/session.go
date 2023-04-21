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

type Session interface {
	Get(url string, params map[string]interface{}) (map[string]interface{}, error)
	Post(url string, data map[string]interface{}) (map[string]interface{}, error)
	Put(url string, data map[string]interface{}) (map[string]interface{}, error)
	Delete(url string, data map[string]interface{}) (map[string]interface{}, error)
	Patch(url string, data map[string]interface{}) (map[string]interface{}, error)
	UploadFile(url string, field string, filePath string, data map[string]string) (map[string]interface{}, error)
	AddHeader(key string, value string)
}

type session struct {
	headers map[string]string
	verify  bool
	client  http.Client
}

func New(headers map[string]string, verify bool) *session {
	var client http.Client
	if verify == false {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = http.Client{Transport: tr}
	} else {
		client = http.Client{}
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	return &session{
		headers: headers,
		verify:  true,
		client:  client,
	}
}

func (s *session) resetHeaders(res http.Response) {
	s.headers["Cookie"] = res.Header.Get("Cookie")
}

func (s *session) putHeaders(req *http.Request) {
	for key_, value_ := range s.headers {
		req.Header.Set(key_, value_)
	}
}

func (s *session) Headers() map[string]string {
	return s.headers
}

func (s *session) AddHeader(key string, value string) {
	s.headers[key] = value
}

func (s *session) templateRequest(url string, data map[string]interface{}, method string) (map[string]interface{}, error) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	s.putHeaders(req)
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	var resMap map[string]interface{}
	_ = json.Unmarshal(body, &resMap)
	return resMap, nil
}

func (s *session) Get(url string, params map[string]interface{}) (map[string]interface{}, error) {
	res, err := s.templateRequest(url, params, "GET")
	return res, err
}

func (s *session) Post(url string, data map[string]interface{}) (map[string]interface{}, error) {
	res, err := s.templateRequest(url, data, "POST")
	return res, err
}

func (s *session) UploadFile(url string, field string, filePath string, data map[string]string) (map[string]interface{}, error) {
	// 如果读取错误直接将错误返回
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// 创建写入
	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile(field, filePath)
	if err != nil {
		return nil, err
	}
	// 拷贝文件内容
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer res.Body.Close()
	var resMap map[string]interface{}
	body, _ := io.ReadAll(res.Body)
	_ = json.Unmarshal(body, &resMap)
	return resMap, nil
}

func (s *session) Put(url string, data map[string]interface{}) (map[string]interface{}, error) {
	res, err := s.templateRequest(url, data, "PUT")
	return res, err
}

func (s *session) Patch(url string, data map[string]interface{}) (map[string]interface{}, error) {
	res, err := s.templateRequest(url, data, "PATCH")
	return res, err
}

func (s *session) Delete(url string, data map[string]interface{}) (map[string]interface{}, error) {
	res, err := s.templateRequest(url, data, "DELETE")
	return res, err
}
