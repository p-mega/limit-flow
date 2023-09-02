package filter

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/p-mega/limit-flow/utils"
)

var ipMap *sync.Map
var limitedIpMap *sync.Map

func init() {
	ipMap = new(sync.Map)
	limitedIpMap = new(sync.Map)
}

/*
MIN_SAFE_TIME: 用户访问最小安全时间，在该时间内如果访问次数大于阀值，则记录为恶意IP，否则视为正常访问
LIMIT_NUMBER: 用户连续访问最高阀值，超过该值则认定为恶意操作的IP，进行限制
LIMITED_TIME_MILLIS: 默认限制时间（单位：ms）3600000,3600(s)
*/
func GinLimitFlow(MIN_SAFE_TIME, LIMIT_NUMBER, LIMITED_TIME_MILLIS int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if IsNotIllegalRequest(ctx.Request) {
			log.Println("have a illegal request")
			ctx.String(500, "request is not illegal")
			ctx.Abort()
		}
		filterLimitedIpMap()
		ip, err := utils.GetRemoteIP(ctx.Request)
		if err != nil {
			log.Println("request ip not found")
			ctx.String(500, "ip not found")
			ctx.Abort()
		}
		if isLimitedIP(ip) {
			if limitedTime, ok := limitedIpMap.Load(ip); ok {
				limitedTime = limitedTime.(int64) - time.Now().UnixMicro()
				remainingTime := limitedTime.(int64) / 1000
				if limitedTime.(int64)%1000 > 0 {
					remainingTime += 1
				}
				ctx.String(500, fmt.Sprintf("request is frequently, please waiting for %ds", remainingTime))
				ctx.Abort()
			}
		}
		if info, ok := ipMap.Load(ip); ok {
			ipInfo := info.([]int64)
			ipInfo[0] = ipInfo[0] + 1
			if ipInfo[0] > LIMIT_NUMBER {
				ipAccessTime := ipInfo[1]
				currentTimeMillis := time.Now().UnixMicro()
				if currentTimeMillis-ipAccessTime <= MIN_SAFE_TIME {
					limitedIpMap.Store(ip, currentTimeMillis+LIMITED_TIME_MILLIS)
					ctx.String(500, fmt.Sprintf("request is frequently, please waiting for %ds", LIMITED_TIME_MILLIS))
					ctx.Abort()
				} else {
					initIpVisitsNumber(ip)
				}
			}
		} else {
			initIpVisitsNumber(ip)
		}
	}
}

func initIpVisitsNumber(ip string) {
	ipInfo := make([]int64, 2)
	ipInfo[0] = 0
	ipInfo[1] = time.Now().UnixMicro()
	ipMap.Store(ip, ipInfo)
}

func isLimitedIP(ip string) bool {
	_, ok := limitedIpMap.Load(ip)
	return ok
}

func filterLimitedIpMap() {
	currentTimeMillis := time.Now().UnixMicro()
	limitedIpMap.Range(func(key, value any) bool {
		if expireTimeMillis, ok := limitedIpMap.Load(key); ok {
			if expireTimeMillis.(int64) < currentTimeMillis {
				limitedIpMap.Delete(key)
			}
		}
		return true
	})
}

func IsNotIllegalRequest(request *http.Request) bool {
	domain := request.Header.Get("referer")
	return strings.HasPrefix(domain, request.URL.Scheme+"://"+request.TLS.ServerName)
}
