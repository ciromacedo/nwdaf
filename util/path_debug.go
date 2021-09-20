//+build debug

package util

import (
	"free5gc/lib/path_util"
)

var NwdafLogPath = path_util.Gofree5gcPath("free5gc/nwdafsslkey.log")
var NwdafPemPath = path_util.Gofree5gcPath("free5gc/support/TLS/_debug.pem")
var NwdafKeyPath = path_util.Gofree5gcPath("free5gc/support/TLS/_debug.key")
var DefaultNwdafConfigPath = path_util.Gofree5gcPath("free5gc/config/nwdafcfg.yaml")
