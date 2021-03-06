/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
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

//免疫程序连接Jenkins服务器初始化报错
var immunity bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "", "", "Path to config file(指定配置文件路径)")
	rootCmd.PersistentFlags().BoolVarP(&immunity, "immunity", "I", false, "Used to prevent the jenkins server from exiting with an initialized error(免疫Jenkins初始化保存导致的程序退出)")
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
		jenkinsConfig.SetConfigPath(dirname + "/.config/jenkinscli/config.yaml")
	}
	config, err := jenkinsConfig.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//	fmt.Printf("打印Config结构体2%v\n", config)
	jenkinsMod = jenkins.Jenkins{}
	err = jenkinsMod.Init(config) //初始化jenkins对象
	if err != nil {
		if !immunity {
			fmt.Println("❌ jenkins server unreachable: " + jenkinsMod.Server)
			log.Println("连接Jenkins出错,退出程序：", err)
			os.Exit(1)
		} else {
			log.Println("❌ jenkins server unreachable: " + jenkinsMod.Server)
			log.Println("连接Jenkins出错,忽略程序继续执行:", err)
		}

	}
}
