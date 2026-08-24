package main

import (
	"os/exec"
)

// showDialog 用 zenity 弹出可复制的信息窗口
func showDialog(title, body string) error {
	cmd := exec.Command("zenity", "--info",
		"--title="+title,
		"--width=360",
		"--text="+body)
	return cmd.Run()
}

// showError 用 zenity 弹出错误窗口
func showError(msg string) {
	cmd := exec.Command("zenity", "--error", "--title=时间工具", "--width=360", "--text="+msg)
	_ = cmd.Run()
}

// showNotify 用 notify-send 弹出桌面通知
func showNotify(title, body string) error {
	cmd := exec.Command("notify-send", "-a", "时间工具", "-i", "appointment-new", title, body)
	return cmd.Run()
}

// notifyError 用 notify-send 弹出错误通知
func notifyError(msg string) {
	cmd := exec.Command("notify-send", "-a", "时间工具", "-u", "critical", "时间工具 - 出错", msg)
	_ = cmd.Run()
}
