package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/riba2534/webhook/utils"
	log "github.com/sirupsen/logrus"
)

func GitHubEvent(c *gin.Context) {
	hook, _ := github.New(github.Options.Secret("hpc"))
	payload, err := hook.Parse(c.Request, github.PingEvent, github.PushEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
			log.Errorf("解析失败,错误信息为: %s", utils.MarshalAny2String(payload))
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("github.ErrEventNotFound"))
			return
		}
	}
	switch p := payload.(type) {
	case github.PingPayload:
		log.Infof("接收到 ping 事件，内容为: %s", utils.MarshalAny2String(p))
	case github.PushPayload:
		log.Infof("接收到 push 事件，内容为: %s", utils.MarshalAny2String(p))
		go UpdateBlog()
	default:
		log.Info("没有匹配事件")
	}
	c.JSON(http.StatusOK,gin.H{
		"message": "pong",
	})
}

func UpdateBlog() {
	resp, err := exec_shell("/home/ubuntu/hugo-blog/auto_deploy.sh")
	if err != nil {
		log.Errorf("err=%+v", err)
	}
	log.Infof("执行结果为: \n%s\n\n", resp)
}

// 错误处理函数
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// 阻塞式的执行外部 shell 命令的函数，等待执行完毕并返回标准输出
func exec_shell(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行 name 指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)

	//读取 io.Writer 类型的 cmd.Stdout，再通过 bytes.Buffer(缓冲 byte 类型的缓冲器) 将 byte 类型转化为 string 类型 (out.String():这是 bytes 类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run 执行 c 包含的命令，并阻塞直到完成。这里 stdout 被取出，cmd.Wait() 无法正确获取 stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	checkErr(err)
	return out.String(), err
}
