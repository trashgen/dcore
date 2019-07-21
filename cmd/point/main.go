package main

import (
	"log"

	dcconf "dcore/codebase/module/config"
	"dcore/codebase/module/http/persistence"
	dchttpsrv "dcore/codebase/module/http/server"
	dcutil "dcore/codebase/util"
)

func main() {
	c, err := dcutil.LoadJSONConfig(dcconf.NewPointConfig(dcconf.NewMetaConfig()))
	if err != nil {
		log.Fatal(err.Error())
	}

	config, ok := c.(*dcconf.PointConfig)
	if !ok {
		log.Fatal("Config: type mismatch")
	}

	dchttpsrv.NewPoint(config, persistence.NewMockPersistModule()).Start()
}
