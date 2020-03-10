package blink

//#include "blink.h"
import "C"

import (
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/lxn/win"
	"github.com/wsteel/blink/internal/devtools"
	"os"
	"path/filepath"
	"runtime"
)

//任务队列,保证所有的API调用都在痛一个线程
var jobQueue = make(chan func())

//初始化blink,释放并加载dll,启动调用队列
func InitBlink() error {

	errdir := os.MkdirAll(TempPath, 0644)
	if errdir != nil {
		return fmt.Errorf("无法创建临时目录：%s, err: %s", TempPath, errdir)
	}

	//定义dll的路径
	dllPath := filepath.Join(TempPath, "blink_"+runtime.GOARCH+".dll")

	path, errx := os.Executable()
	if errx != nil {
		return fmt.Errorf("无法获取当前目录：%s, err: %s", "now", errx)
	}

	ThisPath := filepath.Dir(path)
	dllPath = filepath.Join(ThisPath, "blink_"+runtime.GOARCH+".dll")

	// 判断文件是否存在
	_, err := os.Stat(dllPath) //os.Stat获取文件信息
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("无法获取dll文件：%s", dllPath)
	}

	//启动一个新的协程来处理blink的API调用
	go func() {
		//将这个协程锁在当前的线程上
		runtime.LockOSThread()
		//初始化
		C.initBlink(
			C.CString(dllPath),
			C.CString(TempPath),
			C.CString(filepath.Join(TempPath, "cookie.dat")),
		)

		if isDebug {
			//注册DevTools工具到虚拟文件系统
			RegisterFileSystem("__devtools__", &assetfs.AssetFS{
				Asset:     devtools.Asset,
				AssetDir:  devtools.AssetDir,
				AssetInfo: devtools.AssetInfo,
			})
		}

		//消费API调用,同时处理好windows消息
		for {
			select {
			case job := <-jobQueue:
				job()
			default:
				//消息循环
				msg := &win.MSG{}
				if win.GetMessage(msg, 0, 0, 0) != 0 {
					win.TranslateMessage(msg)
					//是否传递下去
					next := true
					//拿到对应的webview
					view := getWebViewByHandle(msg.HWnd)
					if view != nil {
						next = view.processMessage(msg)
					}
					if next {
						win.DispatchMessage(msg)
					}
				}
			}
		}
	}()

	logger.Println("blink初始化完毕")

	return nil
}
