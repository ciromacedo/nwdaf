//+build !debug

package util

import (
	"github.com/free5gc/path_util"
)

var NwdafLogPath = path_util.Free5gcPath("free5gc/nwdafsslkey.log")
var NwdafPemPath = path_util.Free5gcPath("free5gc/support/TLS/nwdaf.pem")
var NwdafKeyPath = path_util.Free5gcPath("free5gc/support/TLS/nwdaf.key")
var DefaultNwdafConfigPath = path_util.Free5gcPath("free5gc/config/nwdafcfg.yaml")
