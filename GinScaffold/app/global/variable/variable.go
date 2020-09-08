package variable

import (
	"goskeleton/app/global/my_errors"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
)

var (
	BasePath           string       // 定义项目的根目录
	EventDestroyPrefix = "Destroy_" // 程序退出时需要销毁的时间前缀
	// 上传文件保存路径
	UploadFileField    = "files"
	UploadFileSavePath = "/storage/app/uploaaded/"

	ZapLog *zap.Logger

	WebsocketHub              interface{}
	WebsocketHandshakeSuccess = "Websocket Handshake+OnOpen Success"
	WebsocketServerPingMsg    = "Server->Ping->Client"

	// 用户自行定义其他全局变量 ↓
)

func init() {
	// 初始化程序根目录
	if path, err := os.Getwd(); err != nil {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		BasePath = strings.Replace(strings.Replace(path, `\test`, "", 1), `/test`, "", 1)
	} else {
		log.Fatal(my_errors.ErrorsBasePath)
	}
}
