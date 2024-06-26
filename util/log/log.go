package log

import (
	"bytes"

	e "github.com/carlos-yuan/cargen/core/error"
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
	if err == nil {
		return
	}
	var buffer bytes.Buffer
	log.print(&buffer, err)
}

func (log *CarLogger) PrintWarn(err error) {
	if err == nil {
		return
	}
	var buffer bytes.Buffer
	log.print(&buffer, err)
	log.Warn(buffer.String())
}

func (log *CarLogger) PrintInfo(err error) {
	if err == nil {
		return
	}
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
