package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var path string

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Start a Jenkins resource and you can trigger an artifact download",
	Long: `Starts a Jenkins resource and can trigger an artifact download; 
	        the default is not to trigger an artifact download if the download path is specified`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("launch called")
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: JOB_NAME")
			os.Exit(1)
		}
		qid, err := jenkinsMod.Instance.BuildJob(jenkinsMod.Context, args[0], nil) //构建并返回队列id
		fmt.Println("------queueid:", qid)
		if err != nil {
			panic(err)
		}
		build, err := jenkinsMod.Instance.GetBuildFromQueueID(jenkinsMod.Context, qid) //通过队列id返回构建对象
		if err != nil {
			panic(err)
		}
		for build.IsRunning(jenkinsMod.Context) {
			time.Sleep(5000 * time.Millisecond)
			build.Poll(jenkinsMod.Context)
		}
		fmt.Printf("build number %d with result: %v\n", build.GetBuildNumber(), build.GetResult())

		if path != "" && build.GetResult() == "SUCCESS" {
			fmt.Println("⏳ downloading artifacts...")
			err := jenkinsMod.DownloadArtifacts(args[0], build.GetBuildNumber(), path)
			if err != nil {
				fmt.Printf("cannot download artifacts: %s\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	launchCmd.Flags().StringVarP(&path, "path", "p", "", "Specify a directory for downloading artifacts")
	rootCmd.AddCommand(launchCmd)

}
