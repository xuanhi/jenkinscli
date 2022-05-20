package jenkins

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//保存了ssh连接的基本信息
type SshC struct {
	User     string `mapstructure:"User"`
	Password string `mapstructure:"Password"`
	Host     string `mapstructure:"Host"`
	Port     string `mapstructure:"Port"`
	Timeout  time.Duration
}

type SftpC struct {
	Host string
	//ssh 客户端句柄
	SshClient *ssh.Client
	//sftp 客户端句柄
	Client *sftp.Client
}

//初始化SshC对象
func NewSshC(user, password, host, port string) *SshC {
	return &SshC{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Timeout:  10 * time.Second,
	}
}

//创建一个sshclient客户端句柄
func (s *SshC) SshClient() *ssh.Client {
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
		log.Fatalln("sshclient", err.Error())
	}
	return sshClient
}

//初始化SftpC对象 同时创建了sftpclient 客户端句柄
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

//创建一个sftpclient客户端句柄
func (f *SftpC) SftpClient() {
	sftpClient, err := sftp.NewClient(f.SshClient)
	if err != nil {
		log.Fatalln("sftpclient", err.Error())
	}
	f.Client = sftpClient
}

//上传文件 指定文件路径到远程目录下
func (f *SftpC) UploadFile(localFilePath, remoteFilePath string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := f.Client.Create(path.Join(remoteFilePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remoteFilePath, remoteFileName))
		log.Fatal(err)
	}
	defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		log.Fatal(err)
	}
	dstFile.Write(ff)
	fmt.Println(localFilePath + " copy file to remote server finished!")
}

//上传目录 上传本地目录下所有文件到远程目录下，不会将指定目录的父目录上传到远程目录下，只会上传内容
func (f *SftpC) UploadDirectory(localPath string, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal("read dir list fail ", err)
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
	fmt.Println(localPath + " copy directory to remote server finished!")
}

//下载文件 指定本地目录，指定远程文件下载目录和文件
func (f *SftpC) DownLoadFile(localpath, remotepath string) {
	srcFile, err := f.Client.Open(remotepath)
	if err != nil {
		log.Fatal("文件读取失败", err)
	}
	defer srcFile.Close()
	localFilename := path.Base(remotepath)
	dstFile, err := os.Create(path.Join(localpath, localFilename))
	if err != nil {
		log.Fatalln("文件创建失败", err)
	}
	defer dstFile.Close()
	if _, err := srcFile.WriteTo(dstFile); err != nil {
		log.Fatal("文件写入失败", err)
	}
	fmt.Println(remotepath, "文件下载成功")
}

//下载目录，将远端目录下载到本地目录下
func (f *SftpC) DownLoadDir(localpath, remotepath string) {
	remotefiles, err := f.Client.ReadDir(remotepath)
	if err != nil {
		log.Fatal("remote read dir list fail ", err)
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
}
