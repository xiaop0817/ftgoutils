package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xiaop0817/ftgoutils/c"
	"io/ioutil"
	"log"
	"net/http"
	"unsafe"
)

var Debug = true
var prefix = c.C(fmt.Sprintf("%s", "[HTTP]"), c.LightBlue)

func debug(f string, v ...interface{}) {
	if Debug {
		log.Printf(f, v...)
	}
}

// PostJson POST
func PostJson(url string, body interface{}, entity interface{}, header map[string]string) error {
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
	debug("%s %s", prefix, c.C(str, c.LightGreen))
	err = json.Unmarshal([]byte(*str), entity)
	return err
}
