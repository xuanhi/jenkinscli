/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/xuanhi/jenkinscli/jenkins"
	"github.com/xuanhi/jenkinscli/utils/zaplog"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a resource in Jenkins(启动Jenkins流水线,相对于launch功能少很多且实现方法不一样)",
}

var enableJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Enable job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("❌ requires at least one argument")
		}
		zaplog.Sugar.Infof("⏳ Enabling job %s...\n", args[0])
		job, err := jenkinsMod.Instance.GetJob(jenkinsMod.Context, args[0])
		if err != nil {
			zaplog.Sugar.Errorf("unable to find the job: %s - err: %s \n", args[0], err)
			os.Exit(1)
		}
		mapv := jenkins.ArgstoMap(args)
		// mapv := map[string]string{
		// 	"xhh":       "123456789",
		// 	"xhhstring": "remote string",
		// 	"xhhtext":   "remote text",
		// 	"xhhradio":  "三",
		// }
		//job.Enable(jenkinsMod.Context)
		queueid, err := job.InvokeSimple(jenkinsMod.Context, mapv)
		if err != nil {
			panic(err)
		}
		build, err := jenkinsMod.Instance.GetBuildFromQueueID(jenkinsMod.Context, queueid)
		if err != nil {
			panic(err)
		}
		//wait for build to finish
		for build.IsRunning(jenkinsMod.Context) {
			time.Sleep(5000 * time.Millisecond)
			build.Poll(jenkinsMod.Context)
		}
		zaplog.Sugar.Infof("build number %d with result: %v\n", build.GetBuildNumber(), build.GetResult())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)
	enableCmd.AddCommand(enableJobCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
