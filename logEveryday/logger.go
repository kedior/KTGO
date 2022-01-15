package logEveryday

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogEveryday 让你的logger文件每天24点准时替换,log文件保存的folderPath里
//以2021-11-22.log的格式命名
func LogEveryday(logger *log.Logger, t time.Time, folderPath string) error {
	os.MkdirAll(folderPath, os.ModePerm)
	logFile, err := os.OpenFile(folderPath+t.String()[:10]+".log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open file failed: %w", err)
	}
	logger.SetOutput(logFile)
	date, month, day := t.Date()
	next := time.Date(date, month, day+1, 0, 0, 0, 0, time.Local)
	time.AfterFunc(next.Sub(time.Now())+time.Second, func() {
		err := LogEveryday(logger, next, folderPath)
		if err != nil {
			panic(err)
		}
		logFile.Close()
	})
	return nil
}
