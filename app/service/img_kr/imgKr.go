package img_kr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"

	"github.com/valyala/fasthttp"
)

func Upload(uploadFile string) map[string]interface{} {

	bodyBufer := &bytes.Buffer{}
	//創建一個multipart文件寫入器，方便按照http規定格式寫入內容
	bodyWriter := multipart.NewWriter(bodyBufer)
	//從bodyWriter生成fileWriter,並將文件內容寫入fileWriter,多個文件可進行多次
	fileWriter, err := bodyWriter.CreateFormFile("file", path.Base(uploadFile))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	file, err := os.Open(uploadFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//不要忘記關閉打開的文件
	defer file.Close()
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Println(err.Error())
	}

	//關閉bodyWriter停止寫入數據
	bodyWriter.Close()

	contentType := bodyWriter.FormDataContentType()
	//構建request，發送請求
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer func() {
		// 用完需要釋放資源
		fasthttp.ReleaseResponse(response)
		fasthttp.ReleaseRequest(request)
	}()

	request.Header.SetContentType(contentType)
	//直接將構建好的數據放入post的body中
	request.SetBody(bodyBufer.Bytes())
	request.Header.SetMethod("POST")

	request.Header.SetBytesKV([]byte("Referer"), []byte("https://imgkr.com/"))

	request.SetRequestURI("https://imgkr.com/api/v2/files/upload")
	err = fasthttp.Do(request, response)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var res map[string]interface{}
	e := json.Unmarshal(response.Body(), &res)
	if e != nil {
		log.Println(e)
	}
	return res
}
