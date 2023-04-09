package test

import (
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"testing"
)

func TestBar(t *testing.T) {

	file, _ := os.OpenFile("D:\\download\\upgrade2500.jtl", os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()

	//创建进度条
	bar := progressbar.DefaultBytes(
		-1, //总大小，-1表示未知
		"uploading",
	)

	//将文件copy到新的文件中，同时显示进度条
	newfile, _ := os.OpenFile("test2.jtl", os.O_CREATE|os.O_RDWR, 0644)

	io.Copy(io.MultiWriter(newfile, bar), file)
}

func TestIO(t *testing.T) {

	n, _ := io.WriteString(os.Stdout, "hello world")

	t.Log(n)
}
