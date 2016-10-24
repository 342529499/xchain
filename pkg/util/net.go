package util

import "net"

var Local_Address_IP string

func SetLocalIP(ip string) error {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback then display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if ip == ipnet.IP.String() {
					Local_Address_IP = ip
					return nil
				}
			}
		}
	}

	return nil
}

func GetLocalIP() string {
	if Local_Address_IP != "" {
		return Local_Address_IP
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback then display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
