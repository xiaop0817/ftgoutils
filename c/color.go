package c

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type color string

const (
	BgBlack     color = "\u001B[40;37m"
	BgRed       color = "\u001B[41;37m"
	BgGreen     color = "\u001B[42;37m"
	BgYellow    color = "\u001B[43;37m"
	BgBlue      color = "\u001B[44;37m"
	BgPink      color = "\u001B[45;37m"
	BgLightBlue color = "\u001B[46;37m"
	BgWhite     color = "\u001B[47;37m"
	Bold        color = "\033[1m"
	Yellow      color = "\033[33m"
	Cyan        color = "\033[36m"
	Gray        color = "\033[90m"
	Red         color = "\033[31m"
	Blue        color = "\033[34m"
	Pink        color = "\033[35m"
	Green       color = "\033[32m"
	LightRed    color = "\033[91m"
	LightGreen  color = "\033[92m"
	LightYellow color = "\033[93m"
	LightBlue   color = "\033[94m"
	LightPink   color = "\033[95m"
	LightCyan   color = "\033[96m"
	White       color = "\033[97m"
	Black       color = "\033[30m"
	Revert      color = "\033[7m"
	End         color = "\033[0m"
)

func C(s interface{}, c color) string {
	return fmt.Sprintf("%s%s%s", c, Strval(s), End)
}

func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
