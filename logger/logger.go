package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"lualsp/auxiliary"
	"os"
	"path/filepath"
)

var (
	//debug 记录所有日志
	debugLog *log.Logger
	//info 重要的信息
	infoLog *log.Logger
	//warning 需要注意的信息
	warningLog *log.Logger
	//error 非常严重的问题
	errorLog *log.Logger
)

type outWay int

const (
	Wfile = 1 << iota
	Wstd
	Wall
)

func getWriter(outputway outWay, file io.Writer, std io.Writer) io.Writer {
	switch outputway {
	case Wfile:
		return file
	case Wstd:
		return std
	case Wall:
		return io.MultiWriter(file, std)
	default:
		return ioutil.Discard
	}
}

//Init 初始化
func Init(level string, outputway outWay) {
	switch level {
	case "debug", "info", "warning", "error":
	default:
		panic("level value is err!")
	}
	debugLog = log.New(ioutil.Discard, "[D]: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog = log.New(ioutil.Discard, "[I]: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog = log.New(ioutil.Discard, "[W]: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(ioutil.Discard, "[E]: ", log.Ldate|log.Ltime|log.Lshortfile)

	errorLog.SetOutput(getWriter(outputway, getfile("log/err.txt"), os.Stderr))
	if level == "error" {
		return
	}
	warningLog.SetOutput(getWriter(outputway, getfile("log/warn.txt"), os.Stdout))
	if level == "warning" {
		return
	}
	infoLog.SetOutput(getWriter(outputway, getfile("log/info.txt"), os.Stdout))
	if level == "info" {
		return
	}
	debugLog.SetOutput(getWriter(outputway, getfile("log/debug.txt"), os.Stdout))
}

//Debugf 调试日志格式化输出
func Debugf(format string, v ...interface{}) {
	debugLog.Output(2, fmt.Sprintf(format, v...))
}

//Debugln 调试日志输出
func Debugln(v ...interface{}) {
	debugLog.Output(2, fmt.Sprintln(v...))
}

//Infof 日志格式化输出
func Infof(format string, v ...interface{}) {
	infoLog.Output(2, fmt.Sprintf(format, v...))
}

//Infoln 日志输出
func Infoln(v ...interface{}) {
	infoLog.Output(2, fmt.Sprintln(v...))
}

//Warningf 警告日志格式化输出
func Warningf(format string, v ...interface{}) {
	warningLog.Output(2, fmt.Sprintf(format, v...))
}

//Warningln 警告日志输出
func Warningln(v ...interface{}) {
	warningLog.Output(2, fmt.Sprintln(v...))
}

//Errorf 错误异常日志格式化输出
func Errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	errorLog.Output(2, s)
	panic(s)
}

//Errorln 错误异常日志输出
func Errorln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	errorLog.Output(2, s)
	panic(s)
}

func getfile(path string) io.Writer {
	dir := filepath.Dir(path)
	if exist, _ := auxiliary.PathExists(dir); !exist {
		os.MkdirAll(dir, os.ModePerm)
	}
	file, err := os.OpenFile(path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open ", path, " file:", err)
		return ioutil.Discard
	}

	return file
}
