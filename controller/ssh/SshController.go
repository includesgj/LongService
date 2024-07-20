package sshtro

import (
	"GinProject12/model"
	"GinProject12/response"
	sshs "GinProject12/serverce/cmd/ssh"
	utssh "GinProject12/util/ssh"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"time"
)

var server sshs.SSHService

// SshService 修改ssh服务
// @Summary      修改ssh服务
// @Description  修改ssh服务 operate有两种值 1, enable (开启) 2, disable (关闭)
// @Tags         ssh
// @Accept       json
// @Produce      json
// @Param        operate query string ture "operate"
// @Success      200
// @Failure      201
// @Failure      404
// @Failure      500
// @Router       /ssh/operate [GET]
func SshService(c *gin.Context) {
	operate := c.Query("operate")
	if operate != "disable" && operate != "enable" {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数不正确")
		return
	}

	if err := server.OperationSsh(operate); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, nil, fmt.Sprintf("%s成功", operate))
}

// SshConnect 修改ssh服务
// @Summary      连接远程终端
// @Description  连接远程终端
// @Tags         ssh
// @Accept       json
// @Produce      json
// @Param        host query string ture "ip地址"
// @Param 	     port query string ture "端口"
// @Param  		 password query string ture "密码"
// @Param		 username query string ture "用户名"
// @Success      200
// @Failure      201
// @Failure      404
// @Failure      500
// @Router       /ssh/connect [GET]
func SshConnect(c *gin.Context) {

	host := c.Query("host")
	port := c.Query("port")
	password := c.Query("password")
	username := c.Query("username")
	fmt.Println("host", host, "port", port, "pwd", password, "user", username)

	vmInfo := model.ConnectRequest{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}

	config := &ssh.ClientConfig{
		User: vmInfo.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(vmInfo.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//测试连接
	flag, connectID := utssh.VerifyConnect(&vmInfo)
	if flag != true {
		response.Response(c, http.StatusInternalServerError, 500, nil, "测试连接失败")
		return
	}
	//连接成功

	//创建客户端与会话
	conn, err := ssh.Dial("tcp", vmInfo.Host+":"+vmInfo.Port, config)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "unable to create session"+err.Error())
		return
	}
	defer session.Close()

	//设置流
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	defer stdinPipe.Close()

	// Prepare pipes for capturing output
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	stderrPipe, err := session.StderrPipe()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	//升级为websocket
	ws, err := model.UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "ws升级失败"+err.Error())
		return
	}
	defer ws.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// Request a pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "request for pseudo terminal failed: "+err.Error())
		return
	}

	// Start a shell
	if err := session.Shell(); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err)
		return
	}

	// 测试
	commands := []string{" "}
	for _, cmd := range commands {
		_, err = stdinPipe.Write([]byte(cmd))
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		}
	}

	/// Custom buffers
	var stdoutBuf sshs.OutputDataBuffer
	// Capture stdout and stderr
	go io.Copy(&stdoutBuf, stdoutPipe)
	go io.Copy(&stdoutBuf, stderrPipe)

	//将虚拟机的初始化信息输出到ws上
	//延时
	time.Sleep(500 * time.Millisecond)

	stdoutBuf.Flush(ws, "cmd", connectID)

	model.ConnHeartBeat[ws] = time.Now()

	// 接受前端的数据

	var inputBuffer sshs.OutputDataBuffer

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "read:"+err.Error())
			break
		}

		var cmdData model.ShellRequest
		err = json.Unmarshal(message, &cmdData)
		if err != nil {
			response.Fail(c, nil, err.Error())
			continue // 解析失败时继续监听下一条消息
		}

		switch cmdData.Type {
		case model.Cmd:
			// 将接收到的消息放入缓冲区
			inputBuffer.Write([]byte(cmdData.Command))

			//将缓冲区的数据写入服务器
			inputBuffer.WriteToHost(stdinPipe)

			//延时
			time.Sleep(50 * time.Millisecond)

			// 虚拟机的数据写入websocket
			stdoutBuf.Flush(ws, "cmd", connectID)
		case model.HeartBeat:

			model.ConnHeartBeat[ws] = time.Now()
			stdoutBuf.Write([]byte("pong"))
			stdoutBuf.Flush(ws, "pong", connectID)
		}

	}

}
