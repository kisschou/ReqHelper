package ReqHelper

import (
	"encoding/json"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
)

type (
	Request struct {
		Host    string
		IpAddr  string
		Header  map[string]string
		Params  map[string]interface{}
		Get     map[string]interface{}
		Post    map[string]interface{}
		Put     map[string]interface{}
		File    *RequestFile
		IsGet   bool
		IsPost  bool
		IsPut   bool
		IsOpt   bool
		IsDel   bool
		IsPatch bool
	}

	RequestFile struct {
		Filename string
		Header   map[string][]string
		Size     int64
		Body     multipart.File
	}
)

func New(req *http.Request) {
	r = &Request{"", "", make(map[string]string, 0), make(map[string]interface{}, 0), make(map[string]interface{}, 0), make(map[string]interface{}, 0), make(map[string]interface{}, 0), new(RequestFile), false, false, false, false, false, false}

	// Header
	for k, v := range req.Header {
		r.Header[k] = v[0]
	}

	// Get
	for k, v := range req.URL.Query() {
		r.Get[k] = v[0]
	}

	// 请求参数
	if _, ok := req.Header["Content-Type"]; ok {
		if strings.Contains(req.Header["Content-Type"][0], "json") {
			jsonParams := make(map[string]interface{}, 0)
			decoder := json.NewDecoder(req.Body)
			_ = decoder.Decode(&jsonParams)
			r.Put = jsonParams
		}

		if strings.Contains(req.Header["Content-Type"][0], "x-www-form-urlencoded") {
			for k, v := range req.PostForm {
				r.Post[k] = v[0]
			}
		}

		if strings.Contains(req.Header["Content-Type"][0], "form-data") {
			if fileBody, fileHeader, err := req.FormFile("file"); err == nil {
				defer fileBody.Close()

				r.File = &RequestFile{
					Filename: fileHeader.Filename,
					Header:   fileHeader.Header,
					Size:     fileHeader.Size,
					Body:     fileBody,
				}
			}

			for k, v := range req.PostForm {
				r.Post[k] = v[0]
			}
		}
	}

	// Get|Put|Delete请求藏在地址中的参数
	if req.Method == "GET" || req.Method == "PUT" || req.Method == "DELETE" {
		for _, v := range c.Params {
			r.Get[v.Key] = v.Value
		}
	}

	// 获取请求地址
	r.Host = req.Host

	// 获取客户端ip地址
	r.IpAddr = GetIpAddr(req)

	// 判断请求类型
	r = r.checkReqMethod(req.Method)

	// 合并参数到Params
	return r.merge2Params()
}

// checkReqMethod 检测请求类型
func (r *Request) checkReqMethod(method string) *Request {
	switch method {
	case "GET":
		r.IsGet = true
		break
	case "POST":
		r.IsPost = true
		break
	case "PUT":
		r.IsPut = true
		break
	case "DELETE":
		r.IsDel = true
		break
	case "OPTIONS":
		r.IsOpt = true
		break
	case "PATCH":
		r.IsPatch = true
		break
	}
	return r
}

// Merge Merge multiple maps.
func Merge(inputs ...map[string]interface{}) map[string]interface{} {
	if len(inputs) < 1 {
		return nil
	}
	result := make(map[string]interface{}, 0)
	for _, input := range inputs {
		for k, v := range input {
			result[k] = v
		}
	}
	return result
}

// merge2Params 所有参数合并.
func (r *Request) merge2Params() *Request {
	r.Params = Merge(r.Get, r.Post, r.Put)
	return r
}

// GetIpAddr 获取IP地址
func GetIpAddr(req *http.Request) string {
	ip := strings.TrimSpace(strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0])
	if ip == "" {
		ip = strings.TrimSpace(req.Header.Get("X-Real-Ip"))
	}
	if ip == "" {
		var err error
		if ip, _, err = net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err != nil {
			ip = ""
		}
	}
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}
