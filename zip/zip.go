package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/xiaop0817/ftgoutils/c"
	"log"
	"os"
)

var prefix = c.C(fmt.Sprintf("%-10s", "[ZIP]"), c.LightCyan)

// ZipItem zip内包含文件
// Name 文件名
// Content 内容byte数组
type ZipItem struct {
	Name    string
	Content []byte
}

// WriteToZip 写入内容到Zip文件
// @param fileName Zip文件路径
// @param items zip内文件列表
func WriteToZip(fileName string, items []ZipItem) {
	err := os.Remove(fileName)
	if err != nil {
		log.Printf("%s 清理文件[%s]失败", prefix, c.C(fileName, c.LightRed))
	}

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)
	for _, zi := range items {
		if zi.Name == "" {
			continue
		}
		zipItem, _ := writer.Create(zi.Name)
		i, _ := zipItem.Write(zi.Content)
		log.Printf("%s 文件[%s]共写入[%s]bytes", prefix, c.C(zi.Name, c.LightGreen), c.C(i, c.LightRed))
	}
	err = writer.Close()
	if err != nil {
		log.Printf("%s %s:%s", prefix, c.C("Writer关闭发生错误", c.LightRed), err)
		return
	}
	ff, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	_, err = buf.WriteTo(ff)
	if err != nil {
		log.Printf("%s %s:%s", prefix, c.C("写入Zip发生错误", c.LightRed), err)
		return
	}
	log.Printf("%s 文件已生成到[%s]", prefix, c.C(fileName, c.LightRed))
}
