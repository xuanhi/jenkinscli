package jenkins

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bndr/gojenkins"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

//jiekins 连接对象
type Jenkins struct {
	Instance    *gojenkins.Jenkins
	Server      string
	JenkinsUser string
	Token       string
	Context     context.Context

	//邮箱
	MailSmpt   string
	MailPort   int
	MailUser   string
	MailToken  string
	MailFrom   string   //发送邮箱
	MailTo     []string //主送
	MailCc     []string //抄送
	MailBcc    []string //密送
	MailSub    string   //主题标题
	MailBody   string   //邮箱内容
	MailAttach string   //附件路径

}

//配置被集中在json文件中
type Config struct {
	Server         string `mapstructure:"Server"`
	JenkinsUser    string `mapstructure:"JenkinsUser"`
	Token          string `mapstructure:"Token"`
	ConfigPath     string
	ConfigFileName string
	ConfigFullPath string

	MailSmpt  string   `mapstructure:"MailSmpt"`
	MailPort  int      `mapstructure:"MailPort"`
	MailUser  string   `mapstructure:"MailUser"`
	MailToken string   `mapstructure:"MailToken"`
	MailFrom  string   `mapstructure:"MailFrom"` //发送邮箱
	MailTo    []string `mapstructure:"MailTo"`   //主送
	MailCc    []string `mapstructure:"MailCc"`   //抄送
	MailBcc   []string `mapstructure:"MailBcc"`  //密送
	MailSub   string   `mapstructure:"MailSub"`  //主题标题
	//	MailBody   string   `mapstructure:"MailBody"`   //邮箱内容
	MailAttach string `mapstructure:"MailAttach"` //附件路径
}

//设置默认配置路径
func (j *Config) SetConfigPath(path string) {
	dir, file := filepath.Split(path)
	//fmt.Println("当前文件路径", dir)
	j.ConfigPath = dir
	j.ConfigFileName = file
	j.ConfigFullPath = j.ConfigPath + j.ConfigFileName
}

//从指定路径加载配置文件
func (j *Config) LoadConfig() (config Config, err error) {
	viper.AddConfigPath(j.ConfigPath)
	viper.SetConfigName(j.ConfigFileName)
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	//fmt.Println(viper.GetString("Server"))
	//	fmt.Printf("打印Config结构体%v\n", config)

	return
}

//init 将会初始化连接jenkins server
func (j *Jenkins) Init(config Config) error {
	j.JenkinsUser = config.JenkinsUser
	j.Server = config.Server
	j.Token = config.Token
	j.Context = context.Background()

	j.Instance = gojenkins.CreateJenkins(
		nil,
		j.Server,
		j.JenkinsUser,
		j.Token,
	)
	_, err := j.Instance.Init(j.Context)

	j.MailSmpt = config.MailSmpt
	j.MailPort = config.MailPort
	j.MailUser = config.MailUser
	j.MailToken = config.MailToken
	j.MailFrom = config.MailFrom
	j.MailTo = config.MailTo
	j.MailCc = config.MailCc
	j.MailBcc = config.MailBcc
	j.MailAttach = config.MailAttach
	j.MailSub = config.MailSub
	//	j.MailBody = config.MailBody

	return err
}

//下载制品库artifacts
func (j *Jenkins) DownloadArtifacts(jobName string, buildID int64, pathToSave string) error {
	job, err := j.Instance.GetJob(j.Context, jobName)
	if err != nil {
		return errors.New("❌ unable to find the job")
	}
	build, err := job.GetBuild(j.Context, buildID)
	if err != nil {
		return errors.New("❌ unable to find the specific build id")
	}
	artifacts := build.GetArtifacts()
	if len(artifacts) <= 0 {
		fmt.Printf("No artifacts available for download\n")
		return nil
	}
	for _, a := range artifacts {
		fmt.Printf("Saving artifact %s in %s\n", a.FileName, pathToSave)
		_, err := a.SaveToDir(j.Context, pathToSave)
		if err != nil {
			return errors.New("❌ unable to download artifact")
		}
	}
	return nil
}

//显示构建队列
func (j *Jenkins) ShowBuildQueue() error {
	queue, _ := j.Instance.GetQueue(j.Context)
	totalTasks := 0
	for i, item := range queue.Raw.Items {
		fmt.Printf("Name: %s\n", item.Task.Name)
		fmt.Printf("ID: %d\n", item.ID)
		j.ShowStatus(item.Task.Color)
		fmt.Printf("Pending: %v\n", item.Pending)
		fmt.Printf("Stuck: %v\n", item.Stuck)

		fmt.Printf("Why: %s\n", item.Why)
		fmt.Printf("URL: %s\n", item.Task.URL)
		fmt.Printf("\n")
		totalTasks = i + 1
	}
	fmt.Printf("Number of tasks in the build queue: %d\n", totalTasks)

	return nil
}

//显示对象的状态
func (j *Jenkins) ShowStatus(object string) {
	switch object {
	case "blue":
		fmt.Printf("Status: ✅ Success\n")

	case "red":
		fmt.Printf("Status: ❌ Failed\n")

	case "red_anime", "blue_anime", "yellow_anime", "gray_anime", "notbuild_anime":
		fmt.Printf("Status: ⏳ In Progress\n")

	case "notbuilt":
		fmt.Printf("Status: 🚧 Not Build\n")

	default:
		if len(object) > 0 {
			fmt.Printf("Status: %s\n", object)
		}
	}
}

//显示所有views
func (j *Jenkins) ShowViews() error {
	views, err := j.Instance.GetAllViews(j.Context)
	if err != nil {
		return err
	}
	for _, view := range views {
		fmt.Printf("✅ %s\n", view.GetName())
		fmt.Printf("%s\n", view.GetUrl())
		fmt.Printf("\n")
		for _, job := range view.GetJobs() {
			fmt.Printf("✅ %s\n", job.Name)
			fmt.Printf("%s\n", job.Url)
		}
		fmt.Printf("\n")
	}
	return nil
}

//显示所有jobs
func (j *Jenkins) ShowAllJobs() error {
	jobs, err := j.Instance.GetAllJobs(j.Context)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		fmt.Printf("✅ %s\n", job.Raw.Name)
		j.ShowStatus(job.Raw.Color)
		fmt.Printf("%s\n", job.Raw.Description)
		fmt.Printf("%s\n", job.Raw.URL)
		fmt.Printf("\n")
	}
	return nil
}

//获取最后的build
func (j *Jenkins) GetLastBuild(jobName string) error {
	fmt.Printf("⏳ Collecting job information...\n")
	job, err := j.Instance.GetJob(j.Context, jobName)
	if err != nil {
		return errors.New("❌ unable to find the last build job")
	}
	build, err := job.GetLastBuild(j.Context)
	if err != nil {
		return errors.New("❌ unable to find the last build job")
	}

	if len(build.Job.Raw.LastBuild.URL) > 0 {
		fmt.Printf("✅ Last build Number: %d\n", build.Job.Raw.LastBuild.Number)
		fmt.Printf("✅ Last build URL: %s\n", build.Job.Raw.LastBuild.URL)
		fmt.Printf("✅ Parameters: %s\n", build.GetParameters())
	} else {
		fmt.Printf("No last build available for job: %s", jobName)
	}
	return nil
}

//获取最后失败的构建
func (j *Jenkins) GetLastFailedBuild(jobName string) error {
	fmt.Printf("⏳ Collecting job information...\n")
	jobObj, err := j.Instance.GetJob(j.Context, jobName)
	if err != nil {
		return errors.New("❌ unable to find the specific job")
	}
	build, err := jobObj.GetLastFailedBuild(j.Context)
	if err != nil {
		return errors.New("❌ unable to get the last failed build")
	}
	if len(build.GetUrl()) > 0 {
		fmt.Printf("Last Failed build Number: %d\n", build.GetBuildNumber())
		fmt.Printf("Last Failed build URL: %s\n", build.GetUrl())
		fmt.Printf("Parameters: %s\n", build.GetParameters())
	} else {
		fmt.Printf("No last failed build available for job")
	}
	return nil
}

//获取最后成功的构建
func (j *Jenkins) GetLastSuccessfulBuild(jobName string) error {
	fmt.Printf("⏳ Collecting job information...\n")
	jobObj, err := j.Instance.GetJob(j.Context, jobName)
	if err != nil {
		return errors.New("❌ unable to find the specific job")
	}
	build, err := jobObj.GetLastSuccessfulBuild(j.Context)
	if err != nil {
		return errors.New("❌ unable to get the last successful build")
	}
	if len(build.GetUrl()) > 0 {
		fmt.Printf("✅ Last Successful build Number: %d\n", build.GetBuildNumber())
		fmt.Printf("✅ Last Successful build URL: %s\n", build.GetUrl())
		fmt.Printf("✅ Parameters: %s\n", build.GetParameters())
	} else {
		fmt.Printf("No last successful build available for job")
	}
	return nil
}

//获取最后不稳定的构建
func (j *Jenkins) GetLastUnstableBuild(jobName string) error {
	fmt.Printf("⏳ Collecting job information...\n")
	job, err := j.Instance.GetJob(j.Context, jobName)
	if err != nil {
		return errors.New("❌ unable to find the specific job")
	}
	build, err := job.GetLastBuild(j.Context)
	if err != nil {
		return errors.New("❌ unable to find the last unstable build job")
	}

	if len(build.GetUrl()) > 0 {
		fmt.Printf("Last unstable build Number: %d\n", build.GetBuildNumber())
		fmt.Printf("Last unstable build URL: %s\n", build.GetUrl())
		fmt.Printf("Parameters: %s\n", build.GetParameters())
	} else {
		fmt.Printf("No last unstable build available for job: %s", jobName)
	}
	return nil
}

//获取最后一个稳定的构建
func (j *Jenkins) GetLastStableBuild(jobName string) error {
	job, err := j.Instance.GetJob(j.Context, jobName)
	if err != nil {
		return errors.New("❌ unable to find the specific job")
	}
	build, err := job.GetLastStableBuild(j.Context)
	if err != nil {
		return errors.New("❌ unable to find the last stable build job")
	}

	if len(build.GetUrl()) > 0 {
		fmt.Printf("✅ Last stable build Number: %d\n", build.GetBuildNumber())
		fmt.Printf("✅ Last stable build URL: %s\n", build.GetUrl())
		fmt.Printf("✅ Parameters: %s\n", build.GetParameters())
	} else {
		fmt.Printf("No last stable build available for job: %s", jobName)
	}
	return nil
}

//获取所有构建id
func (j *Jenkins) GetAllBuildIds(jobName string) error {
	fmt.Printf("⏳ Collecting job information...\n")
	job, err := j.Instance.GetJob(j.Context, jobName)
	if err != nil {
		return err
	}
	buildids, err := job.GetAllBuildIds(j.Context)
	if err != nil {
		return err
	}
	if len(buildids) > 0 {
		for _, build := range buildids {
			buildObj, err := j.Instance.GetBuild(j.Context, jobName, build.Number)
			if err != nil {
				return err
			}
			fmt.Printf("build Number: %d\n", build.Number)
			fmt.Printf("build URL: %s\n", build.URL)
			fmt.Printf("build resoult: %s\n", buildObj.GetResult())
		}
	} else {
		fmt.Printf("No last unstable build available for job: %s", jobName)
	}
	return nil
}

//显示所有节点实例
func (j *Jenkins) ShowNodes(showStatus string) ([]string, error) {
	var hosts []string

	nodes, err := j.Instance.GetAllNodes(j.Context)
	if err != nil {
		return hosts, err
	}
	for _, node := range nodes {
		//fetch node data
		switch showStatus {
		case "offline":
			if node.Raw.Offline || node.Raw.TemporarilyOffline {
				fmt.Printf("❌ %s - offline\n", node.GetName())
				fmt.Printf("Reason: %s\n\n", node.Raw.OfflineCauseReason)
			}
			hosts = append(hosts, node.GetName())
		case "online":
			if !node.Raw.Offline {
				fmt.Printf("✅ %s - online\n", node.GetName())
			}
			if node.Raw.Idle {
				fmt.Printf("😴 %s - idle\n", node.GetName())
			}
			hosts = append(hosts, node.GetName())
		}

	}
	return hosts, nil
}

//发送邮件
func (j *Jenkins) SendMail(number int64, result, name string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", j.MailFrom)
	m.SetHeader("To", j.MailTo...)
	m.SetHeader("Cc", j.MailCc...)
	m.SetHeader("Bcc", j.MailBcc...)
	m.SetHeader("Subject", j.MailSub)
	//	m.SetBody("text/html", fmt.Sprintf("流水线名称：%s 构建id：%d,构建结构：%s", name, number, result))
	m.SetBody("text/html", Emailpost(number, name, result))
	if j.MailAttach != "" {
		m.Attach(j.MailAttach)
	}
	fmt.Println("附件", j.MailAttach)

	d := gomail.NewDialer(j.MailSmpt, j.MailPort, j.MailUser, j.MailToken)

	err := d.DialAndSend(m)

	return err
}

//自定义发送邮件
func (j *Jenkins) SendMailCustom() error {
	m := gomail.NewMessage()
	m.SetHeader("From", j.MailFrom)
	m.SetHeader("To", j.MailTo...)
	m.SetHeader("Cc", j.MailCc...)
	m.SetHeader("Bcc", j.MailBcc...)
	m.SetHeader("Subject", j.MailSub)
	m.SetBody("text/html", j.MailBody)
	if j.MailAttach != "" {
		m.Attach(j.MailAttach)
	}
	fmt.Println("附件", j.MailAttach)
	d := gomail.NewDialer(j.MailSmpt, j.MailPort, j.MailUser, j.MailToken)
	fmt.Println(d)

	err := d.DialAndSend(m)

	return err
}

//读取文件内容
func (j *Jenkins) ReadFile(filepath string) []byte {
	if _, err := os.Stat(filepath); err != nil {
		fmt.Println("文件不存在或指定的不是文件")
		panic(err)
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File contents:%s", content)
	return content
}
