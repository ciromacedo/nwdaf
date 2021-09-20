package main

import (
	"fmt"
	"github.com/free5gc/NFs/nwdaf/logger"
	nwdaf_service "github.com/free5gc/NFs/nwdaf/service"
	"github.com/free5gc/NFs/nwdaf/version"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var NWDAF = &nwdaf_service.NWDAF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {

	app := cli.NewApp()
	app.Name = "nwdaf"
	fmt.Print(app.Name, "\n")
	appLog.Infoln("NWDAF version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -nwdafcfg nwdaf configuration file"
	app.Action = action
	app.Flags = NWDAF.GetCliCmd()

	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Warnf("Error args: %v", err)
	}
}

func action(c *cli.Context) error{

	//app.AppInitializeWillInitialize(c.String("free5gccfg"))
	//NWDAF.Initialize(c)
	//NWDAF.Start()

	if err := NWDAF.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	NWDAF.Start()

	return nil
}
