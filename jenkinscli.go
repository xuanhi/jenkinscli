/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/xuanhi/jenkinscli/cmd"
	"github.com/xuanhi/jenkinscli/utils/zaplog"
)

func main() {
	zaplog.InitLogger()
	//	defer zaplog.SyncLogger()
	cmd.Execute()
}
