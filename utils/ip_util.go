package utils

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

func notFound(ip string) bool {
	return len(ip) == 0 || strings.Compare(strings.ToLower(ip), strings.ToLower("unknown")) == 0
}

func GetRealIp(req *http.Request) string {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		ip = req.RemoteAddr
	}
	if ip != "127.0.0.1" {
		return ip
	}
	xRealIP := req.Header.Get("X-Real-Ip")
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	for _, address := range strings.Split(xForwardedFor, ",") {
		address = strings.TrimSpace(address)
		if address != "" {
			return address
		}
	}
	if xRealIP != "" {
		return xRealIP
	}
	return ip
}

// func GetRemoteIP(request *http.Request) (string, error) {
// 	ip := request.Header.Get("x-forwarded-for")
// 	var err error = nil
// 	if notFound(ip) {
// 		ip = request.Header.Get("Proxy-Client-IP")
// 	}
// 	if notFound(ip) {
// 		ip = request.Header.Get("WL-Proxy-Client-IP")
// 	}
// 	if notFound(ip) {
// 		ip = request.Header.Get("HTTP_CLIENT_IP")
// 	}
// 	if notFound(ip) {
// 		ip = request.Header.Get("HTTP_X_FORWARDED_FOR")
// 	}
// 	if notFound(ip) {
// 		ip = request.RemoteAddr
// 		if ip == "127.0.0.1" || ip == "0:0:0:0:0:0:0:1" {
// 			ip, err = ips()
// 		}
// 	}
// 	return ip, err
// }

func ips() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return "", err
		}
		addresses, err := byName.Addrs()
		if err != nil {
			return "", err
		}
		for _, v := range addresses {
			if ipnet, flag := v.(*net.IPNet); flag && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", errors.New("ip not found")
}
