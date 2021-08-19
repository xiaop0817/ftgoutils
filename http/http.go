package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"unsafe"
)

func PostJson(url string, body interface{}, entity *interface{}, header map[string]string) error {
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
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	err = json.Unmarshal([]byte(*str), entity)
	return err
}
