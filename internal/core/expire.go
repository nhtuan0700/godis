package core

import (
	"time"

	"github.com/nhtuan0700/godis/internal/constant"
)

func ActiveDeleteExpiredKeys(redisDB *RedisDB) {
	for {
		var expiredCount = 0
		var sampleCountRemain = constant.ActiveExpireSampleSized

		for key, expiredTime := range redisDB.GetExpireDict() {
			sampleCountRemain--
			// get first ActiveExpireSampleSized elements
			if sampleCountRemain < 0 {
				break
			}

			// if expired then delete and increase expiredcount
			if time.Now().UnixMilli() > int64(expiredTime) {
				redisDB.Delete(key)
				expiredCount++
			}
		}
		if float64(expiredCount)/float64(constant.ActiveExpireSampleSized) <= constant.ActiveExpireThreshold {
			break
		}
	}
}
