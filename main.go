// timetool —— 选中文本中的时间戳/时间互转小工具
//
// 用法：
//
//	timetool              读取鼠标选中文本，弹窗(zenity)展示
//	timetool -notify      读取鼠标选中文本，桌面通知(notify-send)展示
//	timetool -text "..."  直接解析给定文本（用于测试）
//	timetool -stdout      结果打印到标准输出（用于测试/管道）
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		useNotify bool
		toStdout  bool
		text      string
	)
	flag.BoolVar(&useNotify, "notify", false, "使用桌面通知(notify-send)展示，默认使用弹窗(zenity)")
	flag.BoolVar(&toStdout, "stdout", false, "结果打印到标准输出（不弹窗）")
	flag.StringVar(&text, "text", "", "直接解析指定文本，而不是读取选中内容")
	flag.Parse()

	// 1) 取得待解析文本
	raw := text
	if raw == "" {
		sel, err := GetPrimarySelection()
		if err != nil {
			fail(err.Error(), useNotify, toStdout)
			return
		}
		raw = sel
	}

	// 2) 解析
	res, err := Parse(raw)
	if err != nil {
		fail(err.Error(), useNotify, toStdout)
		return
	}

	// 3) 展示
	switch {
	case toStdout:
		fmt.Println(res.Format())
	case useNotify:
		if err := showNotify("时间识别结果", res.FormatShort()); err != nil {
			// 通知失败则退回标准输出
			fmt.Println(res.Format())
		}
	default:
		if err := showDialog("时间识别结果", res.Format()); err != nil {
			fmt.Println(res.Format())
		}
	}
}

func fail(msg string, useNotify, toStdout bool) {
	switch {
	case toStdout:
		fmt.Fprintln(os.Stderr, "错误："+msg)
	case useNotify:
		notifyError(msg)
	default:
		showError(msg)
	}
	os.Exit(1)
}
