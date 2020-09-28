package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const format = "2006-01-02 15:04:05"
const formatD = "2006-01-02"
const formatH = "2006-01-02-15"
const formatT = "15:04:05"
const MB = 1048576

var reg, _ = regexp.CompilePOSIX(".log$")

func (writer *FileWriter) StartRotate() {

	switch writer.options.rotatePolicy {
	case ROTATE_POLICY_SIZE:
		startSizeRotate(writer)
	case ROTATE_POLICY_HOUR:
		startTimeHourRotate(writer)
	default:
		startTimeDayRotate(writer)
	}
}

//按照大小滚动
func startSizeRotate(l *FileWriter) {
	file, err := l.logFile.Stat()
	if err != nil {
		return
	}
	if file.Size()/MB >= l.options.rotateSizeMB {
		max := findMax(l.logFile.Name())
		logRotate(l, fmt.Sprintf("%d", max))
	}
}

//按照时间滚动，每日一个日志
func startTimeDayRotate(l *FileWriter) {
	now := time.Now()
	nextS := now.Add(time.Hour * 24).Format(formatD)
	nextT, _ := time.ParseInLocation(formatD, nextS, time.Local)
	delay := nextT.Sub(now)
	time.AfterFunc(delay, func() {
		logRotate(l, now.Format(formatD))
		//延迟一分钟后重新定时rotate任务
		time.AfterFunc(time.Minute, func() {
			startTimeDayRotate(l)
		})
	})
}

//按照时间滚动，每日一个日志
func startTimeHourRotate(l *FileWriter) {
	now := time.Now()
	nextS := now.Add(time.Hour).Format(formatH)
	nextT, _ := time.ParseInLocation(formatH, nextS, time.Local)
	hour := nextT.Sub(now)
	time.AfterFunc(hour, func() {
		logRotate(l, now.Format(formatH))
		//延迟一分钟后重新定时rotate任务
		time.AfterFunc(time.Minute, func() {
			startTimeHourRotate(l)
		})
	})
}

//放在定时器里执行所以需要加锁
func logRotate(l *FileWriter, suffix string) {
	l.Lock()
	defer l.Unlock()
	log.Println("start log file rotate", suffix)
	l.logFile = rotateFile(l.logFile, suffix)
}

//切换文件，关闭旧文件并且按照规则重命名，打开新文件
func rotateFile(f *os.File, suffix string) *os.File {
	if err := f.Sync(); err != nil {
		log.Println("[ERROR] logger_rotate:", err)
	}
	if err := f.Close(); err != nil {
		log.Println("[ERROR] logger_rotate:", err)
	}
	oldName := f.Name()
	newName := reg.ReplaceAllString(oldName, "") + "_" + suffix + ".log"
	if err := os.Rename(oldName, newName); err != nil {
		log.Println("[ERROR] logger_rotate:", err)
	}
	newFile, _ := newFile(oldName)
	return newFile
}

func findMax(logPath string) int {
	dir := filepath.Dir(logPath)
	files, _ := ioutil.ReadDir(dir)
	i := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		i++
	}
	return i
}

func newFile(logPath string) (*os.File, error) {
	checkDir(logPath)
	return os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

func checkDir(filename string) {
	dir := filepath.Dir(filename)
	_, err := os.Stat(dir)
	if !os.IsExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}
}
