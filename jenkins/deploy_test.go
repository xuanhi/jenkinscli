package jenkins

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"
)

var sshC = NewSshC("root", "ADMIN12345", "192.168.100.34", "22")
var sshclient = sshC.SshClient()
var sftpC = NewSftpC(sshclient)
var localpath string = "/root/aaa"

func TestUploadFile(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Errorf("测试遇到错误")
		}
	}()
	fmt.Println(sftpC.Host)
	sftpC.UploadFile("/root/1.jpg", "/root")
}

func TestUploadDirectory(t *testing.T) {
	sftpC.UploadDirectory("/root/hello", "/root/test")
}

func TestDownLoadFile(t *testing.T) {
	fmt.Println(sftpC.Host)
	sftpC.DownLoadFile("/root", "/root/2.jpg")
}

func TestDownLoadDir(t *testing.T) {
	sftpC.DownLoadDir("/root/aaa", "/root/test")
}

func TestDownLoadDirP(t *testing.T) {
	localIPpath := path.Join(localpath, strings.Split(sftpC.SshClient.RemoteAddr().String(), ":")[0])
	err := os.Mkdir(localIPpath, 0755)
	if err != nil {
		log.Printf("%s:创建ip目录有错误\n", sftpC.SshClient.RemoteAddr().String())
		log.Fatal(err)
	}
	sftpC.DownLoadDir(localIPpath, "/root/test")

}

func TestDownLoadFileP(t *testing.T) {
	sftpC.DownLoadFileP("/root", "/root/2.jpg")
}

func TestUploadFileRegep(t *testing.T) {
	sftpC.UploadFileRegep("/root/aaa", "/root/test", ".*txt$")
}
func TestUploadFileRegepTest(t *testing.T) {
	UploadFileRegepTest("/root/aaa", ".*jpg$")
}
