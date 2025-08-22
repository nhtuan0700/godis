package core

import (
	"time"

	"github.com/nhtuan0700/godis/internal/constant"
)

func ActiveDeleteExpiredKeys() {
	for {
		var expiredCount = 0
		var sampleCountRemain = constant.ActiveExpireSampleSized

		for key, expiredTime := range dictStore.GetExpiredDictStore() {
			sampleCountRemain--
			// get first ActiveExpireSampleSized elements
			if sampleCountRemain < 0 {
				break
			}

			// if expired then delete and increase expiredcount
			if time.Now().UnixMilli() > int64(expiredTime) {
				dictStore.Del(key)
				expiredCount++
			}
		}
		if float64(expiredCount)/float64(constant.ActiveExpireSampleSized) <= constant.ActiveExpireThreshold {
			break
		}
	}
}
