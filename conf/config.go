package conf

/*參數說明
app.port // 應用端口
app.upload_file_path // 圖片上傳的臨時文件夾目錄，絕對路徑！
app.cookie_key // 生成加密session
app.serve_type // 默認請使用GoServe
mysql.dsn // mysql 連接地址dsn
*/

var AppJsonConfig = []byte(`
{
  "app": {
    "port": "8322",
    "upload_file_path": "/Users/royale/go/src/go-gofram-chat/uploads",
    "cookie_key": "4238uihfieh49r3453kjdfg",
    "serve_type": "GoServe"
  },
  "mysql": {
    "dsn": "root:hanshans@tcp(127.0.0.1:3306)/go_frame_chat?charset=utf8&parseTime=True&loc=Local"
  }
}
`)
