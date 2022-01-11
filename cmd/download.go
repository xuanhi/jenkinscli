/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download related commands",
}
var artifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "download artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("❌ requires at three arguments [JOB_NAME BUILD_ID PATH_TO_SAVE_ARTIFACTS]")
			os.Exit(1)
		}
		buildID, _ := strconv.ParseInt(args[1], 10, 64)
		fmt.Println("⏳ downloading artifacts...")
		err := jenkinsMod.DownloadArtifacts(args[0], buildID, args[2])
		if err != nil {
			fmt.Printf("cannot download artifacts: %s\n", err)
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
