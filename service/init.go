package service

import (
	"bufio"
	"fmt"
	mongoDBLibLogger "github.com/free5gc/MongoDBLibrary/logger"
	"github.com/free5gc/nwdaf/analyticsinfo"
	"github.com/free5gc/nwdaf/eventssubscription"
	openApiLogger "github.com/free5gc/openapi/logger"
	pathUtilLogger "github.com/free5gc/path_util/logger"
	"io"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/MongoDBLibrary"
	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	"github.com/free5gc/nwdaf/consumer"
	nwdaf_context "github.com/free5gc/nwdaf/context"
	"github.com/free5gc/nwdaf/factory"
	"github.com/free5gc/nwdaf/logger"
	"github.com/free5gc/nwdaf/util"
	"github.com/free5gc/path_util"
)

type NWDAF struct{}

type (
	Config struct {
		nwdafcfg string
	}
)

var config Config

var nwdafCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "nwdafcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*NWDAF) GetCliCmd() (flags []cli.Flag) {
	return nwdafCLi
}

/*func (*NWDAF) Initialize(c *cli.Context) {

	config = Config{
		nwdafcfg: c.String("nwdafcfg"),
	}

	if config.nwdafcfg != "" {
		factory.InitConfigFactory(config.nwdafcfg)
	} else {
		DefaultNwdafConfigPath := path_util.Gofree5gcPath("free5gc/config/nwdafcfg.yaml")
		factory.InitConfigFactory(DefaultNwdafConfigPath)
	}

	if app.ContextSelf().Logger.NWDAF.DebugLevel != "" {
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.NWDAF.DebugLevel)
		if err != nil {
			initLog.Warnf("Log level [%s] is not valid, set to [info] level", app.ContextSelf().Logger.NWDAF.DebugLevel)
			logger.SetLogLevel(logrus.InfoLevel)
		} else {
			logger.SetLogLevel(level)
			initLog.Infof("Log level is set to [%s] level", level)
		}
	} else {
		initLog.Infoln("Log level is default set to [info] level")
		logger.SetLogLevel(logrus.InfoLevel)
	}

	logger.SetReportCaller(app.ContextSelf().Logger.NWDAF.ReportCaller)

}*/

func (nwdaf *NWDAF) Initialize(c *cli.Context) error {
	config = Config{
		nwdafcfg: c.String("nwdafcfg"),
	}

	if config.nwdafcfg != "" {
		if err := factory.InitConfigFactory(config.nwdafcfg); err != nil {
			return err
		}
	} else {
		DefaultAmfConfigPath := path_util.Free5gcPath("free5gc/config/nwdafcfg.yaml")
		if err := factory.InitConfigFactory(DefaultAmfConfigPath); err != nil {
			return err
		}
	}

	nwdaf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (nwdaf *NWDAF) setLogLevel() {
	if factory.NwdafConfig.Logger == nil {
		initLog.Warnln("NWDAF config without log level setting!!!")
		return
	}

	if factory.NwdafConfig.Logger.UDR != nil {
		if factory.NwdafConfig.Logger.UDR.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NwdafConfig.Logger.UDR.DebugLevel); err != nil {
				initLog.Warnf("UDR Log level [%s] is invalid, set to [info] level",
					factory.NwdafConfig.Logger.UDR.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("UDR Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("UDR Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.NwdafConfig.Logger.UDR.ReportCaller)
	}

	if factory.NwdafConfig.Logger.PathUtil != nil {
		if factory.NwdafConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NwdafConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.NwdafConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.NwdafConfig.Logger.PathUtil.ReportCaller)
	}

	if factory.NwdafConfig.Logger.OpenApi != nil {
		if factory.NwdafConfig.Logger.OpenApi.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NwdafConfig.Logger.OpenApi.DebugLevel); err != nil {
				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.NwdafConfig.Logger.OpenApi.DebugLevel)
				openApiLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				openApiLogger.SetLogLevel(level)
			}
		} else {
			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
			openApiLogger.SetLogLevel(logrus.InfoLevel)
		}
		openApiLogger.SetReportCaller(factory.NwdafConfig.Logger.OpenApi.ReportCaller)
	}

	if factory.NwdafConfig.Logger.MongoDBLibrary != nil {
		if factory.NwdafConfig.Logger.MongoDBLibrary.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NwdafConfig.Logger.MongoDBLibrary.DebugLevel); err != nil {
				mongoDBLibLogger.MongoDBLog.Warnf("MongoDBLibrary Log level [%s] is invalid, set to [info] level",
					factory.NwdafConfig.Logger.MongoDBLibrary.DebugLevel)
				mongoDBLibLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				mongoDBLibLogger.SetLogLevel(level)
			}
		} else {
			mongoDBLibLogger.MongoDBLog.Warnln("MongoDBLibrary Log level not set. Default set to [info] level")
			mongoDBLibLogger.SetLogLevel(logrus.InfoLevel)
		}
		mongoDBLibLogger.SetReportCaller(factory.NwdafConfig.Logger.MongoDBLibrary.ReportCaller)
	}
}

func (nwdaf *NWDAF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range nwdaf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (nwdaf *NWDAF) Start() {
	// get config file info
	config := factory.NwdafConfig
	mongodb := config.Configuration.Mongodb

	initLog.Infof("NWDAF Config Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	// Connect to MongoDB
	MongoDBLibrary.SetMongoDB(mongodb.Name, mongodb.Url)

	initLog.Infoln("Server started")

	router := logger_util.NewGinWithLogrus(logger.GinLog)

	// Order is important for the same route pattern.
	//datarepository.AddService(router)
	analyticsinfo.AddService(router)
	eventssubscription.AddService(router)

	nwdafLogPath := util.NwdafLogPath
	nwdafPemPath := util.NwdafPemPath
	nwdafKeyPath := util.NwdafKeyPath

	self := nwdaf_context.NWDAF_Self()
	util.InitNwdafContext(self)

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)
	profile := consumer.BuildNFInstance(self)
	var newNrfUri string
	var err error
	newNrfUri, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, profile.NfInstanceId, profile)
	if err == nil {
		self.NrfUri = newNrfUri
	} else {
		initLog.Errorf("Send Register NFInstance Error[%s]", err.Error())
	}

	server, err := http2_util.NewServer(addr, nwdafLogPath, router)
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: %+v", err)
	}

	serverScheme := factory.NwdafConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(nwdafPemPath, nwdafKeyPath)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (nwdaf *NWDAF) Exec(c *cli.Context) error {

	//NWDAF.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("nwdafcfg"))
	args := nwdaf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./nwdaf", args...)

	nwdaf.Initialize(c)

	var stdout io.ReadCloser
	if readCloser, err := command.StdoutPipe(); err != nil {
		initLog.Fatalln(err)
	} else {
		stdout = readCloser
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	var stderr io.ReadCloser
	if readCloser, err := command.StderrPipe(); err != nil {
		initLog.Fatalln(err)
	} else {
		stderr = readCloser
	}
	go func() {
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	var err error
	go func() {
		if errormessage := command.Start(); err != nil {
			fmt.Println("command.Start Fails!")
			err = errormessage
		}
		wg.Done()
	}()

	wg.Wait()
	return err
}
