package jenkins

import (
	"fmt"
	"testing"
)

var sshC = NewSshC("root", "ADMIN12345", "192.168.100.34", "22")
var sshclient = sshC.SshClient()
var sftpC = NewSftpC(sshclient)

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
