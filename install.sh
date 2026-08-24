#!/usr/bin/env bash
# 构建 timetool 并注册 GNOME 全局快捷键 Ctrl+Shift+Alt+T
set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="$DIR/timetool"

echo ">> 编译 timetool ..."
( cd "$DIR" && go build -o timetool . )
echo "   完成: $BIN"

# ── 注册 GNOME 自定义快捷键 ──────────────────────────────
# GNOME 用一个 keybinding 路径列表管理自定义快捷键，这里追加一项。
BASE="org.gnome.settings-daemon.plugins.media-keys"
KEY_PREFIX="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings"
SLOT="$KEY_PREFIX/timetool/"
SCHEMA="org.gnome.settings-daemon.plugins.media-keys.custom-keybinding"

echo ">> 注册全局快捷键 Ctrl+Shift+Alt+T ..."

# 读取现有列表
existing=$(gsettings get "$BASE" custom-keybindings)

# 如果还没包含我们的 slot，就加进去
if [[ "$existing" != *"$SLOT"* ]]; then
  if [[ "$existing" == "@as []" || "$existing" == "[]" ]]; then
    newlist="['$SLOT']"
  else
    # 去掉结尾的 ] 再追加
    newlist="${existing%]}, '$SLOT']"
  fi
  gsettings set "$BASE" custom-keybindings "$newlist"
fi
20260820

# 配置这个 slot 的名称/命令/按键
gsettings set "$SCHEMA:$SLOT" name 'timetool 时间戳识别'
gsettings set "$SCHEMA:$SLOT" command "$BIN"
gsettings set "$SCHEMA:$SLOT" binding '<Control><Shift><Alt>t'

echo "   已注册: <Control><Shift><Alt>t -> $BIN"
echo
echo "✅ 安装完成。用鼠标选中一段时间戳或时间，然后按 Ctrl+Shift+Alt+T。"
echo "   如需改用桌面通知展示，把命令改为: $BIN -notify"
echo "   （设置 -> 键盘 -> 自定义快捷键 里可查看/修改）"
