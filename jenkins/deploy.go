package jenkins

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"github.com/xuanhi/jenkinscli/utils/color"
	"github.com/xuanhi/jenkinscli/utils/zaplog"
	"golang.org/x/crypto/ssh"
)

const (
	xx string = "==>"
)

// 保存了ssh连接的基本信息
type SshC struct {
	User       string        `mapstructure:"User"`
	Password   string        `mapstructure:"Password"`
	Host       string        `mapstructure:"Host"`
	Port       string        `mapstructure:"Port"`
	Timeout    time.Duration `mapstructure:"Timeout"`
	Cmd        string        `mapstructure:"Cmd"`
	Disbash    bool          `mapstructure:"Disbash"`
	Privatekey string        `mapstructure:"Privatekey"`
}

type SftpC struct {
	Host string
	//ssh 客户端句柄
	SshClient *ssh.Client
	//sftp 客户端句柄
	Client *sftp.Client
}

// 初始化SshC对象
func NewSshC(user, password, host, port string) *SshC {
	return &SshC{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Timeout:  10 * time.Second,
	}
}

// 初始化SshC对象
func NewSshC2(host, port, user string) *SshC {
	return &SshC{
		Host: host,
		Port: port,
		User: user,
	}
}

// 创建一个sshclient客户端句柄
func (s *SshC) SshClient() *ssh.Client {
	if s.Timeout == 0 {
		s.Timeout = 10 * time.Second
	}
	if s.Port == "" {
		s.Port = "22"
	}
	sshConfig := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         s.Timeout,
	}
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port), sshConfig)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	return sshClient
}

// 通过密钥创建一个sshclient客户端句柄
func (s *SshC) SshClientRsa() *ssh.Client {
	if s.Timeout == 0 {
		s.Timeout = 10 * time.Second
	}
	if s.Port == "" {
		s.Port = "22"
	}
	if s.Privatekey == "" {
		dirname, err := os.UserHomeDir() //获取家目录
		if err != nil {
			log.Fatalf("unable to read UserHomeDir: %v", err)
		}
		s.Privatekey = dirname + "/.ssh/id_rsa"
	}
	//var hostKey ssh.PublicKey

	key, err := ioutil.ReadFile(s.Privatekey)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		//	HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         s.Timeout,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port), config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	return client
}

// SshClient 和SshClientRsa 两个方法合为一个
func (s *SshC) SshClientRsaAndSshClient() *ssh.Client {
	if s.Password == "" {
		return s.SshClientRsa()
	}
	return s.SshClient()

}

// 远程执行bash命令
func (s *SshC) Execbash(cmd string) error {
	//client := s.SshClient()
	client := s.SshClientRsaAndSshClient()
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		log.Println("Failed to create session: ", err)
		return err
	}
	defer session.Close()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(fmt.Sprintf("/usr/bin/bash -c \"%s\"", cmd)); err != nil {
		//log.Println("Failed to run: " + err.Error())
		//	zaplog.Sugar.Errorw("exec bash remote server failed!", "host", client.RemoteAddr().String(), "stderr", stderr.String())
		//zaplog.Sugar.Errorf("status [FAILED] host [%s] stderr:\n %v", client.RemoteAddr().String(), stderr.String())
		//zaplog.Sugar.Errorf("status \033[1;31;40m[FAILED]\033[0m  host \033[1;35;40m [%s]\033[0m  stdout:\n\n %v\n", client.RemoteAddr().String(), stderr.String())
		zaplog.Sugar.Errorf("status %s  host %s  stdout:\n\n %v\n", color.GreenB("FAILED"), color.PurpleB(client.RemoteAddr().String()), color.Red(stderr.String()))
		return err
	}
	//log.Printf("%sremote host:%s exec bash remote server finished!", xx, client.RemoteAddr().String())
	//zaplog.Sugar.Infof("status \033[1;32;40m[SUCCESS]\033[0m  host \033[1;35;40m [%s]\033[0m  stdout:\n\n %v\n\n", client.RemoteAddr().String(), stdout.String())
	zaplog.Sugar.Infof("status %s  host %s stdout:\n\n %v\n\n", color.GreenB("SUCCESS"), color.PurpleB(client.RemoteAddr().String()), stdout.String())
	//fmt.Println(b.String())
	return nil
}

// 初始化SftpC对象 同时创建了sftpclient 客户端句柄
func NewSftpC(sshClient *ssh.Client) *SftpC {
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatalln("sftpclient", err.Error())
	}
	//拿到远端addr
	myhost := sshClient.RemoteAddr()
	return &SftpC{
		SshClient: sshClient,
		Client:    sftpClient,
		Host:      myhost.String(),
	}
}

// 创建一个sftpclient客户端句柄
func (f *SftpC) SftpClient() {
	sftpClient, err := sftp.NewClient(f.SshClient)
	if err != nil {
		log.Fatalln("sftpclient", err.Error())
	}
	f.Client = sftpClient
}

// 远程执行sh脚本
func (f *SftpC) ExecTask(localFilePath, remoteFilePath, arg, sudo string) error {
	// err := os.Chmod(localFilePath, 0755)
	// if err != nil {
	// 	log.Println("添加执行权限遇到错误", err)
	// }
	err := f.UploadFile(localFilePath, remoteFilePath)
	if err != nil {
		log.Println("脚本传输错误")
		return err
	}
	session, err := f.SshClient.NewSession()
	if err != nil {
		log.Println("创建ssh会话失败")
		return err
	}
	defer session.Close()
	remoteFileName := path.Base(localFilePath)
	dstFile := path.Join(remoteFilePath, remoteFileName)
	canshu := arg
	if sudo != "" {
		if err := session.Run(fmt.Sprintf("/usr/bin/bash -c \"echo %s | sudo -S %s %s\"", sudo, dstFile, canshu)); err != nil {
			//log.Println("执行脚本失败")
			log.Printf("%sremote host:%s exec bash remote server Failed!", xx, f.SshClient.RemoteAddr().String())
			return err
		}

		log.Printf("%sremote host:%s exec bash remote server finished!", xx, f.SshClient.RemoteAddr().String())
	} else {
		if err := session.Run(fmt.Sprintf("/usr/bin/sh %s %s", dstFile, canshu)); err != nil {
			//log.Println("执行脚本失败")
			log.Printf("%sremote host:%s exec bash remote server Failed!", xx, f.SshClient.RemoteAddr().String())
			return err
		}
		log.Printf("%sremote host:%s exec bash remote server finished!", xx, f.SshClient.RemoteAddr().String())
	}

	return nil

}

// 上传文件 指定文件路径到远程目录下
func (f *SftpC) UploadFile(localFilePath, remoteFilePath string) error {
	start := time.Now()
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		zaplog.Sugar.Errorf("%s:上传文件有错误之打开文件错误\n", f.SshClient.RemoteAddr().String())
		//	log.Println(err)
		return err
	}
	defer srcFile.Close()
	fs, err := os.Stat(localFilePath)
	if err != nil {
		zaplog.Sugar.Errorln("获取本地文件信息遇到错误：", err)
	}
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := f.Client.Create(path.Join(remoteFilePath, remoteFileName))
	if err != nil {
		zaplog.Sugar.Errorf("%s:上传文件有错误之远程创建文件错误: %v 远程路径: %s\n", f.SshClient.RemoteAddr().String(), err, path.Join(remoteFilePath, remoteFileName))
		//zaplog.Sugar.Errorln("sftpClient.Create error : ", path.Join(remoteFilePath, remoteFileName))
		//	log.Println(err)
		return err
	}
	err = dstFile.Chmod(fs.Mode())
	if err != nil {
		zaplog.Sugar.Errorf("设置权限遇到错误:%v", err)
	}
	defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		zaplog.Sugar.Errorf("%s:上传文件有错误之读取文件错误\n", f.SshClient.RemoteAddr().String())
		zaplog.Sugar.Errorln("ReadAll error : ", localFilePath)
		//	log.Println(err)
		return err
	}
	dstFile.Write(ff)

	//fmt.Println("mod:", fs)
	//f.Client.Chmod(path.Join(remoteFilePath, remoteFileName), fs.Mode())

	//log.Println(f.SshClient.RemoteAddr().String(),localFilePath," copy file to remote server finished!")
	//log.Printf("%sremote host:%s path:%s copy file to remote server finished!", xx, f.SshClient.RemoteAddr().String(), localFilePath)
	elapsed := time.Since(start).String()
	zaplog.Sugar.Infof("host: %s time %v status %s\n", color.PurpleB(f.SshClient.RemoteAddr().String()), color.CyanB(elapsed), color.GreenB(localFilePath+" copy finished!"))
	return nil
}

// 用于模板动态生成文件再上传至远端服务器
func (f *SftpC) TmplandUploadFile(infilepath, outfilepath string, mod uint32, data map[string]string) error {
	file, err := ioutil.ReadFile(infilepath)
	if err != nil {
		log.Fatalln("读取文件错误", err)
	}
	t := template.Must(template.New("k8s").Parse(string(file)))
	ftpfile, err := f.Client.Create(outfilepath)
	if err != nil {
		log.Fatalln(err)
	}
	if mod != 0 {
		err = ftpfile.Chmod(fs.FileMode(mod))
	}

	t.Execute(ftpfile, data)
	return err
}

// 使用缓存上传文件 指定文件路径到远程目录下
func (f *SftpC) UploadFilebuf(localFilePath, remoteFilePath string) error {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Printf("%s:上传文件有错误\n", f.SshClient.RemoteAddr().String())
		//	log.Println(err)
		return err
	}
	defer srcFile.Close()
	fs, err := os.Stat(localFilePath)
	if err != nil {
		log.Println("获取本地文件信息遇到错误：", err)
	}
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := f.Client.Create(path.Join(remoteFilePath, remoteFileName))
	dstFile.Chmod(fs.Mode())
	if err != nil {
		log.Printf("%s:上传文件有错误\n", f.SshClient.RemoteAddr().String())
		fmt.Println("sftpClient.Create error : ", path.Join(remoteFilePath, remoteFileName))
		//	log.Println(err)
		return err
	}
	defer dstFile.Close()
	//ff, err := ioutil.ReadAll(srcFile)
	buf := bufio.NewReader(srcFile)
	b := make([]byte, 4096)
	for {
		n, err := buf.Read(b)
		if err != nil || err == io.EOF {
			break
		}
		dstFile.Write(b[:n])
	}
	//	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	//	fmt.Printf("\rUploadloading... %s complete", humanize.Bytes(size))
	// if err != nil {
	// 	log.Printf("%s:上传文件有错误\n", f.SshClient.RemoteAddr().String())
	// 	fmt.Println("ReadAll error : ", localFilePath)
	// 	//	log.Println(err)
	// 	return err
	// }
	//dstFile.Write(ff)

	//fmt.Println("mod:", fs)
	//f.Client.Chmod(path.Join(remoteFilePath, remoteFileName), fs.Mode())

	//log.Println(f.SshClient.RemoteAddr().String(),localFilePath," copy file to remote server finished!")
	log.Printf("%sremote host:%s path:%s copy file to remote server finished!", xx, f.SshClient.RemoteAddr().String(), localFilePath)
	return nil
}

// 上传目录 上传本地目录下所有文件到远程目录下，不会将指定目录的父目录上传到远程目录下，只会上传内容
func (f *SftpC) UploadDirectory(localPath string, remotePath string) error {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Printf("%s:上传目录有错误\n", f.SshClient.RemoteAddr().String())
		//log.Fatal("read dir list fail ", err)
		return err
	}
	//	fmt.Printf("___:%v", localFiles)
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())

		remoteFilePath := path.Join(remotePath, backupDir.Name())

		if backupDir.IsDir() {
			f.Client.Mkdir(remoteFilePath)
			f.UploadDirectory(localFilePath, remoteFilePath)
		} else {
			f.UploadFile(path.Join(localPath, backupDir.Name()), remotePath)
		}
	}
	//log.Println(localPath + " copy directory to remote server finished!")
	log.Printf("%sremote host:%s path:%s copy file to remote server finished!", xx, f.SshClient.RemoteAddr().String(), localPath)
	return nil
}

// 下载文件 指定本地目录，指定远程文件下载目录和文件
func (f *SftpC) DownLoadFile(localpath, remotepath string) error {
	srcFile, err := f.Client.Open(remotepath)
	if err != nil {
		log.Printf("%s:下载文件有错误\n", f.SshClient.RemoteAddr().String())
		//log.Fatal("文件读取失败", err)
		return err
	}
	defer srcFile.Close()
	localFilename := path.Base(remotepath)
	dstFile, err := os.Create(path.Join(localpath, localFilename))
	if err != nil {
		log.Printf("%s:下载目录有错误--文件创建失败\n", f.SshClient.RemoteAddr().String())
		//log.Fatalln("文件创建失败", err)
		return err
	}
	defer dstFile.Close()
	if _, err := srcFile.WriteTo(dstFile); err != nil {
		log.Printf("%s:下载文件有错误--文件写入有错误\n", f.SshClient.RemoteAddr().String())
		//log.Fatal("文件写入失败", err)
		return err
	}
	//fmt.Println(remotepath, "文件下载成功")
	log.Printf("%sremote host:%s path:%s copy file to remote server finished!", xx, f.SshClient.RemoteAddr().String(), remotepath)
	return nil
}

// 用于多线程下载文件 指定本地目录，指定远程文件下载目录和文件
// 下载的文件内容放在会以远程ip自动创建目录里
func (f *SftpC) DownLoadFileP(localpath, remotepath string) {
	srcFile, err := f.Client.Open(remotepath)
	if err != nil {
		log.Printf("%s:下载文件有错误\n", f.SshClient.RemoteAddr().String())
		log.Fatal("文件读取失败", err)
	}
	defer srcFile.Close()
	localFilename := path.Base(remotepath)
	localIPpath := path.Join(localpath, strings.Split(f.SshClient.RemoteAddr().String(), ":")[0])
	err = os.Mkdir(localIPpath, 0755)
	if err != nil {
		log.Printf("%s:创建ip文件有错误\n", f.SshClient.RemoteAddr().String())
		log.Fatal(err)
	}
	dstFile, err := os.Create(path.Join(localIPpath, localFilename))
	if err != nil {
		log.Printf("%s:下载文件有错误\n", f.SshClient.RemoteAddr().String())
		log.Fatalln("文件创建失败", err)
	}
	defer dstFile.Close()
	if _, err := srcFile.WriteTo(dstFile); err != nil {
		log.Printf("%s:下载文件有错误\n", f.SshClient.RemoteAddr().String())
		log.Fatal("文件写入失败", err)
	}
	//fmt.Println(remotepath, "文件下载成功")
	log.Printf("%sremote host:%s path:%s copy file to remote server finished!", xx, f.SshClient.RemoteAddr().String(), remotepath)
}

// 下载目录，将远端目录下载到本地目录下
func (f *SftpC) DownLoadDir(localpath, remotepath string) error {
	remotefiles, err := f.Client.ReadDir(remotepath)
	if err != nil {
		log.Printf("%s:下载目录有错误\n", f.SshClient.RemoteAddr().String())
		//log.Fatal("remote read dir list fail ", err)
		return err
	}
	for _, backupDir := range remotefiles {
		remoteFilePath := path.Join(remotepath, backupDir.Name())
		localFilePath := path.Join(localpath, backupDir.Name())
		if backupDir.IsDir() {
			os.Mkdir(localFilePath, backupDir.Mode())
			f.DownLoadDir(localFilePath, remoteFilePath)
		} else {
			f.DownLoadFile(path.Dir(localFilePath), remoteFilePath)
		}
	}
	return nil
}

// 给一个目录，用正则筛选目录下的文件然后上传到远程主机上
func (f *SftpC) UploadFileRegep(localPath, remoteFilePath, reg string) error {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		zaplog.Sugar.Errorf("%s:上传文件有错误之读取本地目录错误\n", f.SshClient.RemoteAddr().String())
		//log.Println("read dir list fail ", err)
		return err
	}
	for _, backupDir := range localFiles {
		if !backupDir.IsDir() {
			if MatchFile(reg, backupDir.Name()) {
				//log.Println("匹配的文件", backupDir.Name())
				localFilePath := path.Join(localPath, backupDir.Name())
				//log.Println("匹配文件路径", localFilePath)
				zaplog.Sugar.Infof("匹配的文件: %s 全路径: %s", backupDir.Name(), localFilePath)
				err := f.UploadFile(localFilePath, remoteFilePath)
				if err != nil {
					zaplog.Sugar.Errorf("%s传输文件错误:%v", localFilePath, err)
					break
				}
			}
		}
	}
	return nil
}

// 测试正则表达式的接口
func UploadFileRegepTest(localPath, reg string) {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		//	log.Fatal("read dir list fail ", err)
		zaplog.Sugar.Errorf("read dir list fail %v", err)
		return
	}
	for _, backupDir := range localFiles {
		if !backupDir.IsDir() {
			if MatchFile(reg, backupDir.Name()) {
				//	fmt.Println("匹配成功的文件：", backupDir.Name())
				zaplog.Sugar.Infoln("匹配成功的文件：", backupDir.Name())
			}
		}
	}
}

// 判断路径是否存在和是目录
func PathExists(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()

}

// 正则匹配 (?U)非贪婪模式
func MatchFile(reg, name string) bool {
	myRegex, err := regexp.Compile(reg)
	if err != nil {
		log.Println(err)
	}
	return myRegex.MatchString(name)
}

// 正则匹配 key:vaule 形式用于校验数据
func MapFormat(name string) bool {
	myRegex, err := regexp.Compile(".*:.*")
	if err != nil {
		log.Println(err)
	}
	return myRegex.MatchString(name)

}

// 参数构建传入字符串处理
func ArgstoMap(args []string) map[string]string {
	argmap := make(map[string]string)
	//var argmap map[string]string
	for k, arg := range args {
		//跳过第一个参数，因为第一个是构建job名参数
		if k == 0 {
			continue
		}
		if MapFormat(arg) {
			//切割字符串
			kv := strings.Split(arg, ":")
			log.Println("kv:", kv)
			if len(kv) == 2 {
				argmap[kv[0]] = kv[1]
			} else {
				log.Fatalf("一个参数只能有一个：字符")
			}

		} else {
			log.Printf("第%d个参数<%s>格式不正确,将忽略此项设置", k+1, arg)
		}
	}
	return argmap
}

// 动态模板文件方法，给一个输出类型和map参数
func TextTemplate(infilepath, outfilepath string, data map[string]string) {
	f, err := ioutil.ReadFile(infilepath)
	if err != nil {
		log.Fatalln("读取文件错误", err)
	}
	t := template.Must(template.New("k8s").Parse(string(f)))
	if outfilepath == "nil" {
		err := t.Execute(os.Stdout, data)
		if err != nil {
			log.Println("executing template:", err)
		}
	} else {
		f, err := os.Create(outfilepath)
		if err != nil {
			log.Fatalln(err)
		}
		t.Execute(f, data)
	}
}
