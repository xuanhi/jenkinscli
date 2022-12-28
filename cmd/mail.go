/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/xuanhi/jenkinscli/utils/zaplog"
)

var (
	//邮箱内容
	mailbody string
)

// mailCmd represents the mail command
var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Ability to send custom content emails(发送邮件模块)",
	Long: `Ability to send custom content emails,
	You can specify the mail content file,specify the attachment, and the parameters are the message title
	(if no title is specified the title in the configuration file will be used)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			zaplog.Sugar.Errorln("❌ 没有设置邮箱标题，比如 jenkinscli mail \"hello world\"")
			os.Exit(1)
		}
		if mailattach != "" {
			jenkinsMod.MailAttach = mailattach
		}
		if mailbody != "" {
			contents := jenkinsMod.ReadFile(mailbody)
			jenkinsMod.MailBody = string(contents)
		}
		jenkinsMod.MailSub = args[0]
		err := jenkinsMod.SendMailCustom()
		if err != nil {
			zaplog.Sugar.Errorf("发送失败:%v", err)
			return
		}
	},
}

func init() {
	mailCmd.Flags().StringVarP(&mailattach, "attach", "a", "", "Adding attachments to emails(添加邮件附件)")
	mailCmd.Flags().StringVarP(&mailbody, "body", "b", "", "Specify the path to the email content(指定一个文件路径作为邮件内容)")
	//	mailCmd.Flags().StringVarP(&mailsubject, "subject", "s", "", "Set email title")
	rootCmd.AddCommand(mailCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mailCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mailCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
