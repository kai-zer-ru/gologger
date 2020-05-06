package gologger

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestSetLogLevelError(t *testing.T) {
	log := Logger{}
	_ = os.Remove("main.log")
	log.SetLogFileName("main.log")
	err := log.SetLogLevel(213)
	if err == nil {
		t.Error("error is nil")
	}
	_ = log.SetLogLevel(0)
	err = log.Init()
	if err != nil {
		t.Error(err)
	}
	err = log.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestSetLogLevelNoError(t *testing.T) {
	log := Logger{}
	err := log.SetLogLevel(0)
	if err == nil {
		t.Errorf("error is nil")
	}
	err = log.SetLogLevel(1)
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}
	err = log.SetLogLevel(2)
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}
	err = log.SetLogLevel(3)
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}
	err = log.SetLogLevel(4)
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}
	err = log.SetLogLevel(5)
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}
	err = log.SetLogLevel(6)
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}
	_ = os.Remove("main.log")
	log.SetLogFileName("main.log")
	err = log.Init()
	if err != nil {
		t.Error(err)
	}

	err = log.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestSetLogFileName(t *testing.T) {
	_ = os.Remove("main.log")
	fileName := "main.log"
	log := Logger{}
	log.SetLogFileName(fileName)
	err := log.Init()
	if err != nil {
		t.Error(err)
	}
	err = log.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestLogOnClosedLogger(t *testing.T) {
	_ = os.Remove("main.log")
	fileName := "main.log"
	log := Logger{}
	log.SetLogFileName(fileName)
	err := log.SetLogLevel(1)
	if err != nil {
		t.Error(err)
	}
	log.Error("Error")
	log.Debug("Error")
	log.Info("Error")
	log.Fine("Error")
	log.Finest("Error")
	log.Trace("Error")
	log.Warn("Error")
	_ = log.Critical("Error")
	err = log.Close()
	if err != nil {
		if err != ErrLogIsClosed {
			t.Error(err)
		}
	}
	_, err = os.Open("main.log")
	if err == nil {
		t.Errorf("no error on open file")
	}
	if !os.IsNotExist(err) {
		t.Errorf("file exist")
	}
}

func TestLogNoPermissions(t *testing.T) {
	fileName := "/main.log"
	log := Logger{}
	log.SetLogFileName(fileName)
	err := log.SetLogLevel(1)
	if err != nil {
		t.Error(err)
	}
	err = log.Init()
	if err == nil {
		t.Errorf("no error")
	}
}

func TestLogStdOut(t *testing.T) {
	log := Logger{}
	err := log.SetLogLevel(1)
	if err != nil {
		t.Error(err)
	}
	err = log.Init()
	if err != nil {
		t.Error(err)
	}
	log.Debug("Debug")
	err = log.Close()
	if err != nil {
		t.Error(err)
	}
}
func TestLogLogger(t *testing.T) {
	_ = os.Remove("main.log")
	fileName := "main.log"
	log := Logger{}
	log.SetLogFileName(fileName)
	err := log.SetLogLevel(1)
	if err != nil {
		t.Error(err)
	}
	err = log.Init()
	if err != nil {
		t.Error(err)
	}
	log.Error("Error")
	log.Debug("Error")
	log.Info("Error")
	log.Fine("Error")
	log.Finest("Error")
	log.Trace("Error")
	log.Warn("Error")
	err = log.Close()
	if err != nil {
		if err != ErrLogIsClosed {
			t.Error(err)
		}
	}
	_, err = os.Open("main.log")
	if err != nil {
		t.Error(err)
	}
}

func TestChangeLogLevel(t *testing.T) {
	_ = os.Remove("main.log")
	fileName := "main.log"
	log := Logger{}
	log.SetLogFileName(fileName)
	err := log.SetLogLevel(2)
	if err != nil {
		t.Error(err)
	}
	err = log.Init()
	if err != nil {
		t.Error(err)
	}
	log.Trace("Trace string NO LOG")
	log.Debug("Debug string LOG")
	err = log.UpdateLogLevel(1)
	if err != nil {
		t.Error(err)
	}
	log.Trace("Trace string LOG")
	err = log.Close()
	if err != nil {
		t.Error(err)
	}
	file, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
	}
	scanner := bufio.NewScanner(file)

	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, "NO LOG") {
			found = true
		}
	}

	if found {
		t.Errorf("last message not found")
	}
}

func TestTelegram(t *testing.T) {
	token := os.Getenv("token")
	channel := os.Getenv("channel")
	channelInt64 := int64(0)
	t.Log(token)
	t.Log(channel)
	if channel != "" {
		var err error
		channelInt64, err = strconv.ParseInt(channel, 10, 0)
		if err != nil {
			t.Error(err)
			return
		}
	}
	fileName := "main.log"
	log := Logger{}
	log.SetLogFileName(fileName)
	err := log.SetLogLevel(2)
	if err != nil {
		t.Error(err)
		return
	}
	err = log.Init()
	if err != nil {
		t.Error(err)
		return
	}
	if token != "" && channelInt64 != 0 {
		log.EnableTelegram(token, channelInt64, "Test gologger:")
		log.Error("test error")
	} else {
		t.Error("Error get ENV vars")
	}

}

func TestRemoveLogFile(_ *testing.T) {
	_ = os.Remove("main.log")
}
