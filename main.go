package main

import (
	"go-web/server"
	"os"
	"os/exec"
	"os/signal"
)

func main() {

	chChromeDie := make(chan struct{})  //ChromeDie
	chBackendDie := make(chan struct{}) //后端Die
	go server.Run()
	go startBrowser(chChromeDie, chBackendDie) //启动chrome
	chsignal := listenToInterrupt()            //监听中断信号
	for{
	select {                                   //等待中断信号
	case <-chsignal: 
		chBackendDie <-struct{}{}
	case <-chChromeDie:
		os.Exit(0) //退出程序
	}
}
}

func startBrowser(chChromeDie chan struct{}, chBackendDie chan struct{}) {

	chromePath := "C:\\Users\\lenovo\\AppData\\Local\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://localhost:27149/static/index.html")
	cmd.Start()
	go func() {
		<-chBackendDie
		cmd.Process.Kill()
	}()
	go func() {
		cmd.Wait()
		chChromeDie <- struct{}{}
	}()
}

func listenToInterrupt() chan os.Signal {
	// 初始化一个os.Signal类型的channel
	// 我们必须使用缓冲通道，否则在信号发送时如果还没有准备好接收信号，就有丢失信号的风险。
	chsignal := make(chan os.Signal, 1) //监听中断信号
	//监听信号(notify) 用户中断进程 ，就往频道里发一个信息
	signal.Notify(chsignal, os.Interrupt)
	return chsignal
}
