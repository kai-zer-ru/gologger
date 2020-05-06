package gologger

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/kaizer666/log4go"
)

var (
	logIsClosed    = true
	locker         sync.Mutex
	ErrLogIsClosed = fmt.Errorf("log is closed")
)

const (
	TRACE   = log4go.TRACE
	DEBUG   = log4go.DEBUG
	INFO    = log4go.INFO
	WARNING = log4go.WARNING
	ERROR   = log4go.ERROR
	FATAL   = log4go.FATAL
)

var logLevelArr = map[log4go.Level]log4go.Level{
	TRACE:   TRACE,
	DEBUG:   DEBUG,
	INFO:    INFO,
	WARNING: WARNING,
	ERROR:   ERROR,
	FATAL:   FATAL,
}

type Logger struct {
	logger                *log4go.Logger
	logFileName           string
	logLevel              log4go.Level
	logToConsole          bool
	hasWatched            bool
	telegram              telegram
	telegramMessagePrefix string
	telegramConnected     bool
}

func (log *Logger) EnableTelegram(botToken string, channel int64, messagePrefix string) {
	log.telegramConnected = true
	log.telegramMessagePrefix = messagePrefix
	log.telegram = telegram{
		botToken: botToken,
		channel:  channel,
	}
}

func (log *Logger) UpdateLogLevel(logLevel int) error {
	l := log4go.Level(logLevel)
	_, ok := logLevelArr[l]
	if !ok {
		return fmt.Errorf("logLevel %s is wrong", log4go.LevelName(l))
	}
	log.GetLogger().GetLoggerLog4Go().SetLevel(l)
	return nil
}

func (log *Logger) SetLogLevel(logLevel int) error {
	l := log4go.Level(logLevel)
	_, ok := logLevelArr[l]
	if !ok {
		return fmt.Errorf("logLevel %s is wrong", log4go.LevelName(l))
	}
	log.logLevel = l
	return nil
}

func (log *Logger) SetHasWatched() {
	log.hasWatched = true
}

func (log *Logger) SetLogFileName(logFileName string) {
	log.logFileName = logFileName
}

func (log *Logger) GetLogger() *Logger {
	return log
}

func (log *Logger) GetLoggerLog4Go() *log4go.Logger {
	return log.logger
}

func (log *Logger) Init() error {
	var handler *log4go.StreamHandler
	var err error
	var writer io.Writer
	var watchFile = false
	var fileAppend = false
	var handlers []log4go.Handler
	if log.logFileName != "" {
		if log.hasWatched {
			handlerW, err := log4go.NewWatchedFileHandler(log.logFileName, true, true)
			if err != nil {
				return fmt.Errorf("error while create NewWatchedFileHandler: %v", err)
			}
			handlers = append(handlers, handlerW)
		} else {
			handler, err = log4go.NewFileHandler(log.logFileName, true, true)
			if err != nil {
				return fmt.Errorf("error while create NewFileHandler: %v", err)
			}
			handlers = append(handlers, handler)
		}
		watchFile = true
		fileAppend = true
	} else {
		handler, err = log4go.NewStreamHandler(os.Stdout)
		if err != nil {
			return fmt.Errorf("error while create NewStreamHandler: %v", err)
		}
		handlers = append(handlers, handler)
		writer = os.Stdout
	}
	config := log4go.BasicConfigOpts{
		FileName:         log.logFileName,
		FileAppend:       fileAppend,
		WatchFile:        watchFile,
		Format:           "{time} {name<20} {level<8} {message}",
		Writer:           writer,
		Level:            log.logLevel,
		Handlers:         handlers,
		WriteStartHeader: true,
	}
	err = log4go.BasicConfig(config)
	if err != nil {
		return fmt.Errorf("error while create logger: %v", err)
	}
	logIsClosed = false
	log.logger = log4go.GetLogger()
	return nil
}

func (log *Logger) Error(arg0 string, args ...interface{}) {
	if logIsClosed {
		return
	}
	log.sendToTlg("Error", fmt.Sprintf(arg0, args...))
	log.logger.Error(arg0, args...)
}

func (log *Logger) Finest(arg0 string, args ...interface{}) {
	if logIsClosed {
		return
	}
	log.logger.Log(TRACE, arg0, args...)
}

func (log *Logger) Fine(arg0 string, args ...interface{}) {
	if logIsClosed {
		return
	}
	log.logger.Log(DEBUG, arg0, args...)
}

func (log *Logger) Debug(arg0 string, args ...interface{}) {
	if logIsClosed {
		return
	}
	log.logger.Debug(arg0, args...)
}

func (log *Logger) Trace(arg0 string, args ...interface{}) {
	if logIsClosed {
		return
	}
	log.logger.Log(TRACE, arg0, args...)
}

func (log *Logger) Info(arg0 string, args ...interface{}) {
	if logIsClosed {
		return
	}
	log.logger.Info(arg0, args...)
}

func (log *Logger) Warn(arg0 string, args ...interface{}) {
	if logIsClosed {
		return
	}
	log.sendToTlg("Warning", fmt.Sprintf(arg0, args...))
	log.logger.Warning(arg0, args...)
}

func (log *Logger) Critical(arg0 string, args ...interface{}) error {
	if logIsClosed {
		return nil
	}
	log.sendToTlg("Critical", fmt.Sprintf(arg0, args...))
	log.logger.Fatal(arg0, args...)
	return fmt.Errorf("%v", arg0)
}

func (log *Logger) Close() error {
	locker.Lock()
	defer locker.Unlock()
	if !logIsClosed {
		logIsClosed = true
		log4go.Shutdown()
		return nil
	}
	return ErrLogIsClosed
}

func (log *Logger) sendToTlg(level, text string) {
	if log.telegramConnected {
		log.telegram.message = fmt.Sprintf("%s %s: %s", log.telegramMessagePrefix, level, text)
		err := log.telegram.SendMessageToTelegram()
		if err != nil {
			panic(err)
		}
	}
}
