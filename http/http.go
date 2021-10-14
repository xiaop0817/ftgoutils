package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xiaop0817/ftgoutils/c"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"unsafe"
)

var Debug = true
var prefix = c.C(fmt.Sprintf("%s", "[HTTP]"), c.LightBlue)

func debug(f string, v ...interface{}) {
	if Debug {
		log.Printf(f, v...)
	}
}

//PostJson Http.Post
//@param url
//@Param body 请求body
//@param result 返回内容
func PostJson(url string, body interface{}, result interface{}, header map[string]string) error {
	b, _ := json.Marshal(body)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))

	//添加header
	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	return readBody(result, resp)
}

//Get Http.Get
//@param url
//@param param 请求参数map
//@param result 返回内容
func Get(url string, param map[string]interface{}, result interface{}) error {
	paramString := buildParamString(param)
	resp, err := http.Get(url + paramString)
	if err != nil {
		return err
	}
	return readBody(result, resp)
}

//Delete Http.Delete
//@param url
//@param param 请求参数map
//@param result 返回内容
func Delete(url string, param map[string]interface{}, result interface{}) error {
	paramString := buildParamString(param)
	request, _ := http.NewRequest("DELETE", url+paramString, nil)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	return readBody(result, resp)
}

// buildParamString tag:构建参数字符串
func buildParamString(param map[string]interface{}) string {
	var paramString string
	var params []string
	if param != nil {
		for k, v := range param {
			params = append(params, k+"="+c.Strval(v))
		}
	}
	if len(params) > 0 {
		paramString += "?" + strings.Join(params, "&")
	}
	return paramString
}

func readBody(result interface{}, resp *http.Response) error {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	debug("%s %s", prefix, c.C(str, c.LightGreen))
	err = json.Unmarshal([]byte(*str), result)
	return err
}
