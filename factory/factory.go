/*
 * NWDAF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/free5gc/nwdaf/logger"
)

var NwdafConfig Config

func checkErr(err error) {
	if err != nil {
		err = fmt.Errorf("[Configuration] %s", err.Error())
		logger.AppLog.Fatal(err)
	}
}

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		NwdafConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &NwdafConfig); yamlErr != nil {
			return yamlErr
		}
	}

	return nil
}

/*
func InitConfigFactory(f string) {
	content, err := ioutil.ReadFile(f)
	checkErr(err)

	NwdafConfig = Config{}

	err = yaml.Unmarshal([]byte(content), &NwdafConfig)
	checkErr(err)

	logger.InitLog.Infof("Successfully initialize configuration %s", f)
}
*/


func CheckConfigVersion() error {
	currentVersion := NwdafConfig.GetVersion()

	if currentVersion != NWDAF_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("config version is [%s], but expected is [%s].",
			currentVersion, NWDAF_EXPECTED_CONFIG_VERSION)
	}

	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}
