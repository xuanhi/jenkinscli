/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a resource Jenkins",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("❌ requires at least one argument")
		}
		return nil
	},
}
var viewsInfo = &cobra.Command{
	Use:   "views",
	Short: "get all views",
	RunE: func(cmd *cobra.Command, args []string) error {
		jenkinsMod.ShowViews()
		return nil
	},
}

//build
var build = &cobra.Command{
	Use:   "build",
	Short: "build related commands",
}

var buildQueue = &cobra.Command{
	Use:   "queue",
	Short: "get build queue",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("⏳ Collecting build queue information...\n")
		err := jenkinsMod.ShowBuildQueue()
		if err != nil {
			fmt.Println("❌ cannot collect build queue")
			os.Exit(1)
		}
	},
}

//job
var job = &cobra.Command{
	Use:   "job",
	Short: "job related commands",
}
var jobAll = &cobra.Command{
	Use:   "all",
	Short: "get all jobs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("⏳ Collecting all job(s) information...\n")
		err := jenkinsMod.ShowAllJobs()
		if err != nil {
			fmt.Printf("❌ unable to find any job. err: %s \n", err)
			os.Exit(1)
		}
	},
}
var jobGetLastBuild = &cobra.Command{
	Use:   "lastbuild",
	Short: "get last build from a job",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at least one argument [JOB NAME]")
			os.Exit(1)
		}
		err := jenkinsMod.GetLastBuild(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
var jobGetLastSuccessfulBuild = &cobra.Command{
	Use:   "lastsuccessfulbuild",
	Short: "get last successful build from a job",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at least one argument [JOB NAME]")
			os.Exit(1)
		}
		err := jenkinsMod.GetLastSuccessfulBuild(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
var jobLastFailedBuild = &cobra.Command{
	Use:   "lastfailedbuild",
	Short: "get last failed build from a job",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at least one argument [JOB NAME]")
			os.Exit(1)
		}
		err := jenkinsMod.GetLastFailedBuild(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
var jobLastUnstableBuild = &cobra.Command{
	Use:   "lastunstablebuild",
	Short: "get last unstable build from a job",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at least one argument [JOB NAME]")
			os.Exit(1)
		}
		err := jenkinsMod.GetLastUnstableBuild(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
var jobLastStableBuild = &cobra.Command{
	Use:   "laststablebuild",
	Short: "get last stable build from a job",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at least one argument [JOB NAME]")
			os.Exit(1)
		}
		err := jenkinsMod.GetLastStableBuild(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
var jobAllBuildIds = &cobra.Command{
	Use:   "allbuilds",
	Short: "get all build of id from a job",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at least one argument [JOB NAME]")
			os.Exit(1)
		}
		err := jenkinsMod.GetAllBuildIds(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

//node commands
var node = &cobra.Command{
	Use:   "nodes",
	Short: "nodes related commands",
}

var nodesOffline = &cobra.Command{
	Use:   "offline",
	Short: "get nodes offline",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("⏳ Collecting node(s) information...\n")
		hosts, err := jenkinsMod.ShowNodes("offline")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// We must exit as failure in case we have nodes offline
		if len(hosts) > 0 {
			os.Exit(1)
		}
	},
}

var nodesOnline = &cobra.Command{
	Use:   "online",
	Short: "get nodes online",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("⏳ Collecting node(s) information...\n")
		_, err := jenkinsMod.ShowNodes("online")
		if err != nil {
			fmt.Printf("❌ unable to find nodes - err: %s \n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.AddCommand(viewsInfo)
	getCmd.AddCommand(job)
	getCmd.AddCommand(build)
	getCmd.AddCommand(node)

	build.AddCommand(buildQueue)

	job.AddCommand(jobAll)
	job.AddCommand(jobGetLastBuild)
	job.AddCommand(jobGetLastSuccessfulBuild)
	job.AddCommand(jobLastFailedBuild)
	job.AddCommand(jobLastUnstableBuild)
	job.AddCommand(jobAllBuildIds)
	job.AddCommand(jobLastStableBuild)

	node.AddCommand(nodesOffline)
	node.AddCommand(nodesOnline)
}
