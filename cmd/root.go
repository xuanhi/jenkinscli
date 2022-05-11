/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xuanhi/jenkinscli/jenkins"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "jenkinscli",
	Short:   "A client for jenkins",
	Version: "v0.0.1",
	Long:    `Client for jenkins, manage resources by the jenkis`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var jenkinsMod jenkins.Jenkins
var jenkinsConfig jenkins.Config
var configFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "", "", "Path to config file")
}

//加载配置文件
func initConfig() {
	dirname, err := os.UserHomeDir() //获取家目录
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if configFile != "" {
		jenkinsConfig.SetConfigPath(configFile)
	} else {
		jenkinsConfig.SetConfigPath(dirname + "/.config/jenkinscli/config.json")
	}
	config, err := jenkinsConfig.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("打印Config结构体2%v\n", config)
	jenkinsMod = jenkins.Jenkins{}
	err = jenkinsMod.Init(config) //初始化jenkins对象
	if err != nil {
		fmt.Println("❌ jenkins server unreachable: " + jenkinsMod.Server)
		os.Exit(1)
	}

}
