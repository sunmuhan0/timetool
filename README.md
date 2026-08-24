# timetool —— 选中文本时间戳/时间识别小工具

用鼠标选中一段文本，按 **Ctrl+Shift+Alt+T**，自动识别其中的时间戳或时间，
弹窗展示对应的 **时间、时间戳、星期**。

## 支持的输入

| 输入 | 识别为 |
|------|--------|
| `1704153845`（10 位）| 秒级时间戳 |
| `1704153845123`（13 位）| 毫秒级时间戳 |
| `2016-01-02 15:04:05` | 时间字符串（本地时区）|

也兼容 `2016/01/02 15:04:05`、`2016-01-02T15:04:05`、`2016-01-02`、`2016-01-02 15:04` 等；
输入两侧的引号、逗号、空白会被自动清除。

输出内容包含：
- 标准时间 `2006-01-02 15:04:05`
- 中文星期
- 10 位秒级时间戳
- 13 位毫秒级时间戳
- 时区

## 实现说明

- **读取选中文本**：纯 Go 通过 X11（XWayland）协议读取鼠标选中的 `PRIMARY selection`，
  用的是 `github.com/jezek/xgb`，**不依赖 xclip / xsel / wl-clipboard** 等外部程序。
- **展示**：默认用 `zenity` 弹窗（内容可选中复制）；也支持 `notify-send` 桌面通知。
- **全局快捷键**：Wayland 下低层全局热键不可用，采用 GNOME 自定义快捷键
  （`gsettings`）在按下 Ctrl+Shift+Alt+T 时启动本程序，这是 Wayland 下的标准做法。

## 安装

```bash
cd ~/projects/time
./install.sh
```

脚本会编译出 `timetool`，并注册 GNOME 全局快捷键 `Ctrl+Shift+Alt+T`。
安装后直接用鼠标选中文本再按快捷键即可。

> 若快捷键无响应，注销重新登录一次，或到「设置 → 键盘 → 查看及自定义快捷键 →
> 自定义快捷键」确认 “timetool 时间戳识别” 已存在。

## 卸载

```bash
./uninstall.sh
```

移除全局快捷键（二进制文件请自行删除）。

## 命令行用法（调试）

```bash
./timetool                        # 读取选中文本，zenity 弹窗展示（默认）
./timetool -notify                # 读取选中文本，notify-send 通知展示
./timetool -text "1704153845"     # 直接解析指定文本
./timetool -stdout -text "..."    # 解析并打印到标准输出（便于测试/管道）
```

### 改用通知展示

编辑快捷键命令为 `timetool -notify`：

```bash
gsettings set \
  "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/timetool/" \
  command "$(pwd)/timetool -notify"
```

## 目录结构

```
projects/time/
├── main.go         程序入口、命令行参数、流程编排
├── selection.go    X11 PRIMARY selection 读取（纯 Go）
├── timeparse.go    时间戳/时间解析与格式化
├── display.go      zenity 弹窗 / notify-send 通知
├── install.sh      编译 + 注册全局快捷键
├── uninstall.sh    移除全局快捷键
└── README.md
```
