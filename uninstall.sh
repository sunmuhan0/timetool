#!/usr/bin/env bash
# 移除 timetool 注册的 GNOME 全局快捷键
set -euo pipefail

BASE="org.gnome.settings-daemon.plugins.media-keys"
SLOT="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/timetool/"

echo ">> 移除快捷键 slot ..."
existing=$(gsettings get "$BASE" custom-keybindings)
# 从列表中删掉我们的 slot（连同可能的前后逗号/空格）
newlist=$(printf '%s' "$existing" \
  | sed "s#, *'$SLOT'##; s#'$SLOT', *##; s#'$SLOT'##")
# 如果删空了，规整为空列表
[[ "$newlist" == "[]" || "$newlist" == "[ ]" ]] && newlist="@as []"
gsettings set "$BASE" custom-keybindings "$newlist"

echo "✅ 已移除 Ctrl+Shift+Alt+T 快捷键。"
