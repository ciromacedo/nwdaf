package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/logger_conf"
	"github.com/free5gc/logger_util"
)


var (
	log         *logrus.Logger
	AppLog      *logrus.Entry
	InitLog     *logrus.Entry
	CfgLog      *logrus.Entry
	HandlerLog  *logrus.Entry
	DataRepoLog *logrus.Entry
	UtilLog     *logrus.Entry
	HttpLog     *logrus.Entry
	ConsumerLog *logrus.Entry
	GinLog      *logrus.Entry
)


func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	free5gcLogHook, err := logger_util.NewFileHook(logger_conf.Free5gcLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Hooks.Add(free5gcLogHook)
	}

	selfLogHook, err := logger_util.NewFileHook(logger_conf.NfLogDir+"nwdaf.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "CFG"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Handler"})
	DataRepoLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "DataRepo"})
	UtilLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Util"})
	HttpLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "HTTP"})
	GinLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "GIN"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
