/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a resource in Jenkins",
}

var enableJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Enable job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("❌ requires at least one argument")
		}
		fmt.Printf("⏳ Enabling job %s...\n", args[0])
		job, err := jenkinsMod.Instance.GetJob(jenkinsMod.Context, args[0])
		if err != nil {
			fmt.Printf("unable to find the job: %s - err: %s \n", args[0], err)
			os.Exit(1)
		}
		//job.Enable(jenkinsMod.Context)
		queueid, err := job.InvokeSimple(jenkinsMod.Context, nil)
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
		fmt.Printf("build number %d with result: %v\n", build.GetBuildNumber(), build.GetResult())
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
