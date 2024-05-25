package services

//// 存储来自示波器硬件的数据 -> ws发送到前端
//var (
//	dataBuffer []byte
//	mu         sync.Mutex
//)
//
//// ReceiveHardwareData 处理来自硬件的数据
//func ReceiveHardwareData() error {
//	// 初始化 logger
//	logger := logf.Open("HardWare")
//	// 创建 KCP 会话
//	sess, err := kcp.DialWithOptions()
//	if err != nil {
//		logger.Error("Failed to connect to KCP server:", zap.Error(err))
//		return err
//	}
//	// 函数结束时关闭连接
//	defer func() {
//		if err := sess.Close(); err != nil {
//			logger.Error("Error closing sess: %v", zap.Error(err))
//		}
//	}()
//	// 分配内存
//	buf := make([]byte, 4096)
//	for {
//		n, err := sess.Read(buf)
//		if err != nil {
//			logger.Error("Failed to read from KCP server:", zap.Error(err))
//			return err
//		}
//
//		mu.Lock()
//		// 保存数据到 buffer
//		dataBuffer = append(dataBuffer, buf[:n]...)
//		mu.Unlock()
//		return nil
//	}
//}
