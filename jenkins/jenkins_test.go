package jenkins

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	var (
		f        string = "/root/test.txt"
		expected string = "hello world 你好 世界\n"
	)
	jenkins := Jenkins{}
	content := jenkins.ReadFile(f)
	if string(content) != expected {
		t.Errorf("ReadFile(%s) = %s; expected %s", f, content, expected)
	}
}

func TestReadFileStat(t *testing.T) {
	var (
		f string = "/root/a.txt"
	)
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("期望文件不存在，实际文件存在")
		}
	}()
	jenkins := Jenkins{}
	jenkins.ReadFile(f)

}
