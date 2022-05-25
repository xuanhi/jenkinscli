/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	pathx "path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xuanhi/jenkinscli/jenkins"
)

var (
	//一般用于指定目标路径或远端路径
	target string
	//正则表达式语法
	upfilRegexp string
)

// sftpCmd represents the sftp command
var sftpCmd = &cobra.Command{
	Use:   "sftp",
	Short: "Upload or download files or folders to a remote host",
	Long: `Upload or download files or folders to a remote host
	       向远程主机上传或下载文件或文件夹`,
}

//上传文件数据
var sftpUpFile = &cobra.Command{
	Use:   "upfile",
	Short: "Upload files to a remote host",
	Long: `Upload files to a remote host
	       向远程主机上传文件，-t是远程路径,参数为本地路径`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: Local Path")
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
				sshclient := ssh.SshClient()
				sftpc := jenkins.NewSftpC(sshclient)
				sftpc.UploadFile(args[0], target)
			}(ssh)
		}
		jenkinsMod.Wg.Wait()

	},
}

//上传多个文件数据，使用正则表达式匹配
var sftpUpFileRegexp = &cobra.Command{
	Use:   "upfilereg",
	Short: "Upload files to a remote host use Regexp",
	Long: `Upload files to a remote host
	       用正则表达式RE2语法向远程主机上传多个文件,-t是远程路径,参数为本地路径`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: Local Path")
			os.Exit(1)
		}
		if target == "" {
			fmt.Println("❌ requires at one arguments: -t target remote path")
			os.Exit(1)
		}
		if upfilRegexp == "" {
			fmt.Println("❌ requires at one arguments: -R Regular expressions")
			os.Exit(1)
		}

		for _, ssh := range jenkinsMod.SshCs {
			jenkinsMod.Wg.Add(1)
			go func(ssh *jenkins.SshC) {
				defer jenkinsMod.Wg.Done()
				sshclient := ssh.SshClient()
				sftpc := jenkins.NewSftpC(sshclient)
				sftpc.UploadFileRegep(args[0], target, upfilRegexp)
			}(ssh)
		}
		jenkinsMod.Wg.Wait()

	},
}

//测试正则表达是的工具，经常与上传文件配合使用
var Regexptest = &cobra.Command{
	Use:   "regexp",
	Short: "测试正则表达式语法用于配合上传文件",
	Long:  `指定本地目录，通过-R正则表达是来筛选匹配的文件,用于配合上传文件upfilereg进行测试的正则语法检验工具`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: Local Path")
			os.Exit(1)
		}
		if upfilRegexp == "" {
			fmt.Println("❌ requires at one arguments: -R Regular expressions")
			os.Exit(1)
		}

		jenkins.UploadFileRegepTest(args[0], upfilRegexp)

	},
}

//上传目录数据
var sftpUpDir = &cobra.Command{
	Use:   "updir",
	Short: "Upload  folders to a remote host",
	Long: `Upload  folders to a remote host
	       向远程主机上传文件夹，-t是远程路径,参数为本地路径`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: Local Path")
			os.Exit(1)
		}
		if target == "" {
			fmt.Println("❌ requires at one arguments: -t target remote path")
			os.Exit(1)
		}
		if !jenkins.PathExists(args[0]) {
			log.Println("本机目录不存在或路径不是目录")
			os.Exit(1)
		}
		for _, ssh := range jenkinsMod.SshCs {
			jenkinsMod.Wg.Add(1)
			go func(ssh *jenkins.SshC) {
				defer jenkinsMod.Wg.Done()
				sshclient := ssh.SshClient()
				sftpc := jenkins.NewSftpC(sshclient)
				sftpc.UploadDirectory(args[0], target)
			}(ssh)
		}
		jenkinsMod.Wg.Wait()

	},
}

//下载文件数据
var sftpDownFile = &cobra.Command{
	Use:   "downfile",
	Short: "download files to a remote host",
	Long: `download files to a remote host
	       向远程主机下载文件，-t是远程路径,参数为本地路径`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: Local Path")
			os.Exit(1)
		}
		if target == "" {
			fmt.Println("❌ requires at one arguments: -t target remote path")
			os.Exit(1)
		}
		if !jenkins.PathExists(args[0]) {
			log.Println("本机目录不存在或路径不是目录")
			os.Exit(1)
		}
		fmt.Println("主机数量：", len(jenkinsMod.SshCs))
		if len(jenkinsMod.SshCs) != 1 {
			for _, ssh := range jenkinsMod.SshCs {
				jenkinsMod.Wg.Add(1)
				go func(ssh *jenkins.SshC) {
					defer jenkinsMod.Wg.Done()
					sshclient := ssh.SshClient()
					sftpc := jenkins.NewSftpC(sshclient)
					localIPpath := pathx.Join(args[0], strings.Split(sftpc.SshClient.RemoteAddr().String(), ":")[0])
					err := os.Mkdir(localIPpath, 0755)
					if err != nil {
						log.Printf("%s:创建ip目录有错误\n", sftpc.SshClient.RemoteAddr().String())
						log.Fatal(err)
					}
					sftpc.DownLoadFile(localIPpath, target)
				}(ssh)
			}
			jenkinsMod.Wg.Wait()
		} else {
			sshclient := jenkinsMod.SshCs[0].SshClient()
			sftpc := jenkins.NewSftpC(sshclient)
			sftpc.DownLoadFile(args[0], target)
		}
	},
}

//下载目录数据
var sftpDownDir = &cobra.Command{
	Use:   "downdir",
	Short: "download folders to a remote host",
	Long: `download folders to a remote host
	       向远程主机下载文件夹，-t是远程路径,参数为本地路径`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: Local Path")
			os.Exit(1)
		}
		if target == "" {
			fmt.Println("❌ requires at one arguments: -t target remote path")
			os.Exit(1)
		}
		if !jenkins.PathExists(args[0]) {
			log.Println("本机目录不存在或路径不是目录")
			os.Exit(1)
		}
		fmt.Println("主机数量：", len(jenkinsMod.SshCs))
		if len(jenkinsMod.SshCs) != 1 {
			for _, ssh := range jenkinsMod.SshCs {
				jenkinsMod.Wg.Add(1)
				go func(ssh *jenkins.SshC) {
					defer jenkinsMod.Wg.Done()
					sshclient := ssh.SshClient()
					sftpc := jenkins.NewSftpC(sshclient)
					localIPpath := pathx.Join(args[0], strings.Split(sftpc.SshClient.RemoteAddr().String(), ":")[0])
					err := os.Mkdir(localIPpath, 0755)
					if err != nil {
						log.Printf("%s:创建ip目录有错误\n", sftpc.SshClient.RemoteAddr().String())
						log.Fatal(err)
					}
					sftpc.DownLoadDir(localIPpath, target)
				}(ssh)
			}
			jenkinsMod.Wg.Wait()
		} else {
			sshclient := jenkinsMod.SshCs[0].SshClient()
			sftpc := jenkins.NewSftpC(sshclient)
			sftpc.DownLoadDir(args[0], target)
		}
	},
}

func init() {
	sftpCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "Specify a target path or remote path")
	sftpCmd.PersistentFlags().StringVarP(&upfilRegexp, "regexp", "R", "", "Specify a Regexp expressions  with Use RE2 syntax")

	rootCmd.AddCommand(sftpCmd)
	sftpCmd.AddCommand(sftpUpFile)
	sftpCmd.AddCommand(sftpUpDir)
	sftpCmd.AddCommand(sftpDownFile)
	sftpCmd.AddCommand(sftpDownDir)
	sftpCmd.AddCommand(sftpUpFileRegexp)
	sftpCmd.AddCommand(Regexptest)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sftpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sftpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
