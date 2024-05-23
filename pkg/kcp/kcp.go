// 响应KCP

package kcp

import (
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/xtaci/kcp-go/v5"
)

// DialWithOptions 使用配置参数建立 KCP 连接
func DialWithOptions() (*kcp.UDPSession, error) {
	kcpConfig := config.G().Server.Kcp
	// 创建连接
	sess, err := kcp.DialWithOptions(kcpConfig.Addr, nil, kcpConfig.Sndwnd, kcpConfig.Rcvwnd)
	if err != nil {
		return nil, err
	}
	// 配置 KCP 参数
	sess.SetWindowSize(kcpConfig.Sndwnd, kcpConfig.Rcvwnd)
	sess.SetMtu(kcpConfig.Mtu)
	sess.SetNoDelay(kcpConfig.NoDelay, kcpConfig.Interval, kcpConfig.Resend, kcpConfig.NC)

	return sess, nil
}
