package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/beakeyz/dadjoke-gen/pkg/api"
	"github.com/beakeyz/dadjoke-gen/pkg/server"
	"github.com/beakeyz/dadjoke-gen/pkg/setting"
)

const(
  VERSION = "1.0.0"
)

type exitCode = int

var serverFlag = flag.NewFlagSet("dadgen-server", flag.ContinueOnError)

func Bootstrap() exitCode {
	var(
		configFile = serverFlag.String("config", "", "path to config file")
		homePath   = serverFlag.String("homepath", "", "home path, defaults to working directory")
		pidFile    = serverFlag.String("pidfile", "", "path to pid file")

		//add option to debug the requests
	)

	if err := execute(*configFile, *homePath, *pidFile); err != nil {
		fmt.Fprintln(os.Stdout, "Error occoured: %s\n", err.Error())
		return 1
	}

	return 0
}

func execute(config string, home string, pidfile string) error {
  //do more initializing
	var servConf *setting.ServerConfig = &setting.ServerConfig{}
  err := servConf.LoadFromJson(setting.SETTINGS)
  if err != nil {
    fmt.Println(err.Error())
  }

  fmt.Println(servConf.HttpAddress)
  fmt.Println(servConf.HttpPort)
  // Bootstrap http serving
	httpServ, _err := api.BootstrapHttpServer(servConf)
	if (_err != nil){
		fmt.Fprintln(os.Stdout, "Error occoured: %s\n", _err.Error())
		return _err
	}

  // Init server
	s, __err := server.New(httpServ)
	if __err != nil{
		fmt.Fprintln(os.Stdout, "Error occoured: %s\n", __err.Error())
		return __err
	}

	//holy shit this sucks, perhaps use wire?
	if err := s.Run(); err != nil{
		fmt.Fprintln(os.Stdout, "Error occoured: %s\n", err.Error())
		return err
	}

	fmt.Println("Exiting...")
	return nil
}
