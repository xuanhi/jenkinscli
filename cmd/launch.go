package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/xuanhi/jenkinscli/jenkins"
	"github.com/xuanhi/jenkinscli/utils/zaplog"
)

var (
	//下载工件路径
	path string
	//邮箱开关
	mail bool
	//邮箱标题
	mailsubject string
	//邮箱附件
	mailattach string
)

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Start a Jenkins resource and you can trigger an artifact download(启动Jenkins用于java可以下载工件)",
	Long: `	Starts a Jenkins resource and can trigger an artifact download; 
the default is not to trigger an artifact download if the download path is specified;
启动jenkins流水线,第一个参数为流水线名,如果是java还可以构建完成后下载工件到-p指定的本地目录`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("launch called")
		if len(args) < 1 {
			zaplog.Sugar.Errorln("❌ requires at one arguments: JOB_NAME")
			os.Exit(1)
		}
		//注入构建参数
		mapv := jenkins.ArgstoMap(args)
		qid, err := jenkinsMod.Instance.BuildJob(jenkinsMod.Context, args[0], mapv) //构建并返回队列id
		zaplog.Sugar.Infoln("------queueid:", qid)
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
		zaplog.Sugar.Infof("build number %d with result: %v\n", build.GetBuildNumber(), build.GetResult())

		//发邮件
		zaplog.Sugar.Infoln("发邮件入口", mail)
		if mail {
			fmt.Println("⏳ send email...")
			if mailsubject != "" {
				jenkinsMod.MailSub = mailsubject
			}
			if mailattach != "" {
				zaplog.Sugar.Infoln("赋值attach")
				jenkinsMod.MailAttach = mailattach
			}
			err := jenkinsMod.SendMail(build.GetBuildNumber(), build.GetResult(), args[0])
			if err != nil {
				zaplog.Sugar.Errorln("邮件发送出错：", err)

			}
			zaplog.Sugar.Infoln("邮件发送完成")

		}
		//下载工件
		if path != "" && build.GetResult() == "SUCCESS" {
			zaplog.Sugar.Infoln("⏳ downloading artifacts...")
			err := jenkinsMod.DownloadArtifacts(args[0], build.GetBuildNumber(), path)
			if err != nil {
				zaplog.Sugar.Errorf("cannot download artifacts: %s\n", err)
				os.Exit(1)
			}
		}

	},
}

func init() {
	launchCmd.Flags().StringVarP(&path, "path", "p", "", "Specify a directory for downloading artifacts(指定下载工件的路径(目录))")
	launchCmd.Flags().BoolVarP(&mail, "mail", "m", false, "Sending emails with default content(发送邮件开关,默认为fase,需要指定才能发送邮件)")
	launchCmd.Flags().StringVarP(&mailsubject, "subject", "s", "", "Set email title(设置邮箱标题)")
	launchCmd.Flags().StringVarP(&mailattach, "attach", "a", "", "Adding attachments to emails(给邮箱添加附件,指定路径)")
	rootCmd.AddCommand(launchCmd)

}
