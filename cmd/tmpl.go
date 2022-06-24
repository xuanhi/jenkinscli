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

//模板文件路径
var tmpl string

//定义权限
var mod uint32

// tmplCmd represents the tmpl command
var tmplCmd = &cobra.Command{
	Use:   "tmpl",
	Short: "For generating dynamic files(用于生成动态文件)",
	Long: `用于动态文件生成,第一个参数为输出路径，其它参数为替换变量,name:vaule形式,
文件里使用{{.name}}来替换变量name的值,name为输入参数name:value的name一致`,
}

var tmpllocal = &cobra.Command{
	Use:   "local",
	Short: "For generating dynamic files(用于本机生成动态文件)",
	Long: `用于动态文件生成,第一个参数为输出路径，如果第一个参数为'nil'表示输出到终端，其它参数为替换变量,
name:vaule形式,文件里使用{{.name}}来替换变量name的值,name为输入参数name:value的name一致`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("tmpl tmpllocal")
		if len(args) < 1 {
			fmt.Println("❌ requires at one arguments: path")
			os.Exit(1)
		}
		//注入命令行参数到map
		mapv := jenkins.ArgstoMap(args)
		jenkins.TextTemplate(tmpl, args[0], mapv)
	},
}
var tmplremote = &cobra.Command{
	Use:   "remote",
	Short: "For generating dynamic files(用于生成动态文件并发送到远端)",
	Long: `用于动态文件生成,第一个参数为输出路径，其它参数为替换变量,name:vaule形式,
文件里使用{{.name}}来替换变量name的值,name为输入参数name:value的name一致`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("tmpl tmplremote")
		if len(args) < 1 {
			fmt.Println("❌ requires at one arguments: path")
			os.Exit(1)
		}

		//选择主机组
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
					defer sshclient.Close()
					sftpc := jenkins.NewSftpC(sshclient)
					defer sftpc.Client.Close()
					//获取map参数
					data := jenkins.ArgstoMap(args)
					err := sftpc.TmplandUploadFile(tmpl, args[0], mod, data)
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
					defer sshclient.Close()
					sftpc := jenkins.NewSftpC(sshclient)
					defer sftpc.Client.Close()
					//获取map参数
					data := jenkins.ArgstoMap(args)
					err := sftpc.TmplandUploadFile(tmpl, args[0], mod, data)
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
	tmplCmd.PersistentFlags().StringVarP(&tmpl, "tmpl", "l", "", "Specify the template file(指定模板文件路径)")
	tmplremote.Flags().StringVarP(&hostgroup, "hosts", "H", "", "Select the host group you want to activate(选择主机组,不选择默认用Sshs组)")
	tmplremote.Flags().Uint32VarP(&mod, "mod", "M", 0, "Set file permissions(设置文件权限:例如0777)")

	rootCmd.AddCommand(tmplCmd)
	tmplCmd.AddCommand(tmpllocal)
	tmplCmd.AddCommand(tmplremote)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tmplCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tmplCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
