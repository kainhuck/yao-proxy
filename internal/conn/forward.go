package conn

import "net"

func Forward(c1, c2 net.Conn) {
	errChan := make(chan error, 2)
	go func() {
		errChan <- Copy(c1, c2)
	}()
	go func() {
		errChan <- Copy(c2, c1)
	}()

	select {
	case <-errChan:
		return
	}
}
