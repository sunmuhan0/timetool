package main

import (
	"fmt"
	"time"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

// GetPrimarySelection 通过 X11 (XWayland) 读取鼠标选中的文本 (PRIMARY selection)。
// 纯 Go 实现，不依赖 xclip/wl-clipboard 等外部程序。
func GetPrimarySelection() (string, error) {
	X, err := xgb.NewConn()
	if err != nil {
		return "", fmt.Errorf("连接 X 服务器失败（是否在图形会话中运行？）: %w", err)
	}
	defer X.Close()

	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)

	// 创建一个不可见的接收窗口，用于接收 selection 数据
	win, err := xproto.NewWindowId(X)
	if err != nil {
		return "", fmt.Errorf("分配窗口 ID 失败: %w", err)
	}
	// InputOnly 窗口必须 depth=0、visual=CopyFromParent(0)，否则会 BadMatch
	err = xproto.CreateWindowChecked(X, 0, win, screen.Root,
		0, 0, 1, 1, 0,
		xproto.WindowClassInputOnly, 0,
		xproto.CwEventMask, []uint32{uint32(xproto.EventMaskPropertyChange)}).Check()
	if err != nil {
		return "", fmt.Errorf("创建接收窗口失败: %w", err)
	}
	defer xproto.DestroyWindow(X, win)

	// 需要的 atom
	selPrimary := xproto.Atom(xproto.AtomPrimary)
	targetUTF8, err := internAtom(X, "UTF8_STRING")
	if err != nil {
		return "", err
	}
	targetString := xproto.Atom(xproto.AtomString)
	propAtom, err := internAtom(X, "TIMETOOL_SEL")
	if err != nil {
		return "", err
	}

	// 先尝试 UTF8_STRING，失败再退回 STRING
	if txt, ok := requestSelection(X, win, selPrimary, targetUTF8, propAtom); ok {
		return txt, nil
	}
	if txt, ok := requestSelection(X, win, selPrimary, targetString, propAtom); ok {
		return txt, nil
	}

	return "", fmt.Errorf("未能读取到选中文本（请先用鼠标选中一段文本）")
}

func internAtom(X *xgb.Conn, name string) (xproto.Atom, error) {
	reply, err := xproto.InternAtom(X, false, uint16(len(name)), name).Reply()
	if err != nil {
		return 0, fmt.Errorf("intern atom %q 失败: %w", name, err)
	}
	return reply.Atom, nil
}

// requestSelection 请求指定 target 类型的选中内容，并等待 SelectionNotify。
func requestSelection(X *xgb.Conn, win xproto.Window, selection, target, prop xproto.Atom) (string, bool) {
	err := xproto.ConvertSelectionChecked(X, win, selection, target, prop, xproto.TimeCurrentTime).Check()
	if err != nil {
		return "", false
	}

	deadline := time.After(3 * time.Second)
	for {
		select {
		case <-deadline:
			return "", false
		default:
		}

		ev, xerr := X.PollForEvent()
		if xerr != nil {
			return "", false
		}
		if ev == nil {
			// 无事件，稍等再轮询，避免忙等
			time.Sleep(10 * time.Millisecond)
			continue
		}

		sn, ok := ev.(xproto.SelectionNotifyEvent)
		if !ok {
			continue
		}
		if sn.Property == 0 {
			// 拥有者拒绝该 target
			return "", false
		}

		return readProperty(X, win, sn.Property)
	}
}

// readProperty 分块读取窗口属性中的选中文本。
func readProperty(X *xgb.Conn, win xproto.Window, prop xproto.Atom) (string, bool) {
	var data []byte
	var offset uint32
	for {
		reply, err := xproto.GetProperty(X, false, win, prop,
			xproto.GetPropertyTypeAny, offset, 4096).Reply()
		if err != nil {
			return "", false
		}
		if reply.Format == 0 && len(reply.Value) == 0 {
			break
		}
		data = append(data, reply.Value...)
		if reply.BytesAfter == 0 {
			break
		}
		// GetProperty 的 offset 单位是 32-bit word
		offset += uint32(len(reply.Value)) / 4
	}
	// 读取完毕后删除属性
	xproto.DeleteProperty(X, win, prop)

	if len(data) == 0 {
		return "", false
	}
	return string(data), true
}
