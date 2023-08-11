package log

import (
	"bytes"
	e "comm/error"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CarLogger struct {
	*zap.Logger
}

func New(level int8) *CarLogger {
	var cfg zap.Config
	cfg = zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(level))
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.DisableStacktrace = true
	logger := zap.Must(cfg.Build())
	return &CarLogger{Logger: logger}
}

func (log *CarLogger) PrintError(err error) {
	var buffer bytes.Buffer
	log.print(&buffer, err)
	log.Error(buffer.String())
}

func (log *CarLogger) PrintWarn(err error) {
	var buffer bytes.Buffer
	log.print(&buffer, err)
	log.Warn(buffer.String())
}

func (log *CarLogger) PrintInfo(err error) {
	var buffer bytes.Buffer
	log.print(&buffer, err)
	log.Info(buffer.String())
}

func (log *CarLogger) print(siteBuffer *bytes.Buffer, err error) {
	me, ok := err.(e.Err)
	if ok {
		if siteBuffer.Len() == 0 {
			siteBuffer.WriteString(me.CodeName())
			siteBuffer.WriteString(":")
			siteBuffer.WriteString(err.Error())
			siteBuffer.WriteString(":")
		}
		siteBuffer.WriteString(me.Site)
		siteBuffer.WriteString(" -> ")
		if me.Err != nil {
			log.print(siteBuffer, me.Err)
		} else {
			siteBuffer.WriteString(":  " + me.Msg)
		}
	} else {
		if siteBuffer.Len() == 0 {
			siteBuffer.WriteString(e.GetSite(3) + ":  " + err.Error())
		} else {
			siteBuffer.WriteString(":  " + err.Error())
		}
	}
}

func (log *CarLogger) Write(p []byte) (n int, err error) {
	log.Logger.Info(string(p))
	return len(p), nil
}
