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
	"github.com/xuanhi/jenkinscli/utils/zaplog"
)

// 脚本的位置参数
var canshu string

// sudo密码
var sudo string

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "remote exec bash(集成ssh相关功能)",
	Long:  `ssh远程工具`,
}

// 远程执行bash命令
var sshBash = &cobra.Command{
	Use:   "bash",
	Short: "remote exec bash(远程执行bash指令)",
	Long:  `ssh远程工具`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ requires at one arguments: bash ")
			os.Exit(1)
		}
		//声明切片用于记录执行成功和失败的主机
		hostsuccess := make([]string, 0)
		hosterr := make([]string, 0)
		//判断是否指定主机组
		if hostgroup != "" {
			sshc, ok := jenkinsMod.Extend[hostgroup]
			if !ok {
				log.Println("没有找到主机组，请检查是否配置这个主机组")
				return
			}
			for _, ssh := range sshc {
				jenkinsMod.Wg.Add(1)
				go func(ssh *jenkins.SshC) {
					defer jenkinsMod.Wg.Done()
					//sshclient := ssh.SshClient()
					//sftpc := jenkins.NewSftpC(sshclient)
					if !ssh.Disbash {
						//优先执行配置文件的cmd命令
						if ssh.Cmd != "" {
							err := ssh.Execbash(ssh.Cmd)
							if err != nil {
								//log.Println(err)
								zaplog.Sugar.Errorf("主机: %s,执行命令出错: %v", ssh.Host, err)
								hosterr = append(hosterr, ssh.Host)
								return
							}
							hostsuccess = append(hostsuccess, ssh.Host)
						} else {
							err := ssh.Execbash(args[0])
							if err != nil {
								//log.Println(err)
								zaplog.Sugar.Errorf("主机: %s,执行命令出错: %v", ssh.Host, err)
								hosterr = append(hosterr, ssh.Host)
								return
							}
							hostsuccess = append(hostsuccess, ssh.Host)
						}
					}
				}(ssh)
			}
			jenkinsMod.Wg.Wait()
		} else {
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
								//log.Println(err)
								zaplog.Sugar.Errorf("主机: %s,执行命令出错: %v", ssh.Host, err)
								hosterr = append(hosterr, ssh.Host)
								return
							}
							hostsuccess = append(hostsuccess, ssh.Host)
						} else {
							err := ssh.Execbash(args[0])
							if err != nil {
								//log.Println(err)
								zaplog.Sugar.Errorf("主机: %s,执行命令出错: %v", ssh.Host, err)
								hosterr = append(hosterr, ssh.Host)
								return
							}
							hostsuccess = append(hostsuccess, ssh.Host)
						}
					}
				}(ssh)
			}
			jenkinsMod.Wg.Wait()
			zaplog.Sugar.Infof("统计: success: %d \t error: %d \t 错误主机：%v", len(hostsuccess), len(hosterr), hosterr)
		}

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
		if hostgroup != "" {
			sshc, ok := jenkinsMod.Extend[hostgroup]
			if !ok {
				log.Println("没有找到主机组，请检查是否配置这个主机组")
				return
			}
			for _, ssh := range sshc {
				jenkinsMod.Wg.Add(1)
				go func(ssh *jenkins.SshC) {
					defer jenkinsMod.Wg.Done()
					//sshclient := ssh.SshClient()
					sshclient := ssh.SshClientRsaAndSshClient()
					sftpc := jenkins.NewSftpC(sshclient)
					err := sftpc.ExecTask(args[0], target, canshu, sudo)
					if err != nil {
						log.Println(err)
						return
					}
				}(ssh)
			}
			jenkinsMod.Wg.Wait()
		} else {
			for _, ssh := range jenkinsMod.SshCs {
				jenkinsMod.Wg.Add(1)
				go func(ssh *jenkins.SshC) {
					defer jenkinsMod.Wg.Done()
					//sshclient := ssh.SshClient()
					sshclient := ssh.SshClientRsaAndSshClient()
					sftpc := jenkins.NewSftpC(sshclient)
					err := sftpc.ExecTask(args[0], target, canshu, sudo)
					if err != nil {
						log.Println(err)
						return
					}
				}(ssh)
			}
			jenkinsMod.Wg.Wait()
		}

	},
}

func init() {
	sshCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "Specify a target path or remote path(指定远程目录)")
	sshCmd.PersistentFlags().StringVarP(&hostgroup, "hosts", "H", "", "Select the host group you want to activate(选择主机组,不选择默认用Sshs组)")
	sshTask.Flags().StringVarP(&canshu, "canshu", "c", "", "Specify location parameters for the script(为bash脚本键入位置参数,是一个字符串类型多个参数用空格隔开)")
	sshTask.Flags().StringVarP(&sudo, "sudo", "s", "", "Specify the sudo password when executing the script(对于ubuntu系统执行脚本时指定sudo密码)")

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
