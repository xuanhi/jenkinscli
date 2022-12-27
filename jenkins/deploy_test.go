package jenkins

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/xuanhi/jenkinscli/utils/zaplog"
)

var sshC = NewSshC("root", "sqkjp#ssw0rd", "192.168.31.201", "22")

var sshclient = sshC.SshClient()

// var sshC = NewSshC2("192.168.100.34", "22", "root")
// var sshclient = sshC.SshClientRsa()
var sftpC = NewSftpC(sshclient)
var localpath string = "/root/aaa"

func TestUploadFile(t *testing.T) {
	zaplog.InitLogger()
	defer func() {
		err := recover()
		if err != nil {
			t.Errorf("测试遇到错误:%v", err)
		}
	}()
	fmt.Println(sftpC.Host)
	sftpC.UploadFile("/root/go1.19.3.linux-amd64.tar.gz", "/root/test")
}
func TestUploadFilebuf(t *testing.T) {
	sftpC.UploadFilebuf("/root/asitcn-module-system-2.4.6.jar", "/root/test")
}

func TestUploadDirectory(t *testing.T) {
	zaplog.InitLogger()
	sftpC.UploadDirectory("/root/aaa", "/root/test")
}

func TestDownLoadFile(t *testing.T) {
	zaplog.InitLogger()
	fmt.Println(sftpC.Host)
	sftpC.DownLoadFile("/root", "/root/test/a.txt")
}

func TestDownLoadDir(t *testing.T) {
	zaplog.InitLogger()
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
	zaplog.InitLogger()
	sftpC.UploadFileRegep("/root/aaa", "/root/test", ".*txt$")
}
func TestUploadFileRegepTest(t *testing.T) {
	zaplog.InitLogger()
	UploadFileRegepTest("/root/aaa", ".*txt$")
}

func TestExecbash(t *testing.T) {
	zaplog.InitLogger()
	sshC.Execbash("systemctl status nginx", "")
}

func TestExecTask(t *testing.T) {
	zaplog.InitLogger()
	sshC := NewSshC("xuanhi", "xianhuaihai", "192.168.20.129", "22")

	sshclient := sshC.SshClient()

	sftpC := NewSftpC(sshclient)
	sftpC.ExecTask("/root/test.sh", "/xuanhi", "", "")
}

func TestMapFormat(t *testing.T) {
	b := MapFormat("xhhtext:remote text=aa")
	fmt.Println(b)
}

func TestArgstoMap(t *testing.T) {
	aa := []string{"pc-system", "xhhradio:aaa:ccc", "xhhrext:bbb"}
	bb := ArgstoMap(aa)
	fmt.Println(bb)
}

func TestTextTemplate(t *testing.T) {
	data := map[string]string{
		"number": "2",
		"image":  "mysql",
	}
	TextTemplate("/root/tmpl.yaml", "nil", data)
}
