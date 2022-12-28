/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xuanhi/jenkinscli/utils/zaplog"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download related commands(手动下载工件,指定3个参数分别是:流水线名,构建号,工件保存的路径)",
}
var artifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "download artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			zaplog.Sugar.Errorln("❌ requires at three arguments [JOB_NAME BUILD_ID PATH_TO_SAVE_ARTIFACTS]")
			os.Exit(1)
		}
		buildID, _ := strconv.ParseInt(args[1], 10, 64)
		zaplog.Sugar.Infoln("⏳ downloading artifacts...")
		err := jenkinsMod.DownloadArtifacts(args[0], buildID, args[2])
		if err != nil {
			zaplog.Sugar.Errorf("cannot download artifacts: %s\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.AddCommand(artifactsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
