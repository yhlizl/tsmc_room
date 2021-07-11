package controller

import (
	"log"
	"os"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func ImgKrUpload(r *ghttp.Request) {

	file := r.GetUploadFile("file")
	filepath := g.Cfg().GetString("conf.upload_file_path")

	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	filename, err := file.Save(filepath, true)
	if err != nil {
		log.Println(err)
	}

	r.Response.WriteJson(g.Map{
		"code": 0,
		"data": map[string]interface{}{
			"url": "uploads/" + filename,
		},
	})
}
