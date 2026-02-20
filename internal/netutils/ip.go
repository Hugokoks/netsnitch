package netutils

import "net"

func GetLocalIP(target net.IP) (net.IP, error) {
	// make fake UDP connection with target to obtain right ip of responsible network card
	conn, err := net.Dial("udp", target.String()+":80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
