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

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "remote exec bash(集成ssh相关功能)",
	Long:  `ssh远程工具`,
}

//远程执行bash命令
var sshBash = &cobra.Command{
	Use:   "bash",
	Short: "remote exec bash(远程执行bash指令)",
	Long:  `ssh远程工具`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: bash ")
			os.Exit(1)
		}
		for _, ssh := range jenkinsMod.SshCs {
			jenkinsMod.Wg.Add(1)
			go func(ssh *jenkins.SshC) {
				defer jenkinsMod.Wg.Done()
				//sshclient := ssh.SshClient()
				//sftpc := jenkins.NewSftpC(sshclient)
				if !ssh.Disbash {
					if ssh.Cmd != "" {
						err := ssh.Execbash(ssh.Cmd)
						if err != nil {
							log.Println(err)
							return
						}
					} else {
						err := ssh.Execbash(args[0])
						if err != nil {
							log.Println(err)
							return
						}
					}
				}
			}(ssh)
		}
		jenkinsMod.Wg.Wait()

	},
}

var sshTask = &cobra.Command{
	Use:   "task",
	Short: "remote exec bash script(远程执行shell脚本,会将本地脚本先上传再执行脚本)",
	Long:  `ssh远程执行脚本`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: local file path ")
			os.Exit(1)
		}
		if target == "" {
			fmt.Println("❌ requires at one arguments: -t target remote path")
			os.Exit(1)
		}
		for _, ssh := range jenkinsMod.SshCs {
			jenkinsMod.Wg.Add(1)
			go func(ssh *jenkins.SshC) {
				defer jenkinsMod.Wg.Done()
				//sshclient := ssh.SshClient()
				sshclient := ssh.SshClientRsaAndSshClient()
				sftpc := jenkins.NewSftpC(sshclient)
				err := sftpc.ExecTask(args[0], target)
				if err != nil {
					log.Println(err)
					return
				}
			}(ssh)
		}
		jenkinsMod.Wg.Wait()

	},
}

func init() {
	sshCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "Specify a target path or remote path(指定远程目录)")
	rootCmd.AddCommand(sshCmd)
	sshCmd.AddCommand(sshBash)
	sshCmd.AddCommand(sshTask)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
