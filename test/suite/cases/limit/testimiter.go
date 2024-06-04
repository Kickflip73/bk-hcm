package testlimit

import (
	"fmt"
	"hcm/pkg/limit"
	"sync"
	"sync/atomic"
	"time"
)

// TestLimit test limit allow
func TestLimit() {
	var successCount int64 = 0
	var failCount int64 = 0
	var allowTimeSum int64 = 0
	var maxAllowTime int64 = 0
	var minAllowTime int64 = 1000
	wg := sync.WaitGroup{}
	fn := func() {
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				allowStart := time.Now()
				err := limit.GetLimiter().CheckReqCloudAllow("account1", "/api/test1")
				allowTime := time.Now().Sub(allowStart).Milliseconds()
				if allowTime > maxAllowTime {
					maxAllowTime = allowTime
				}
				if allowTime < minAllowTime {
					minAllowTime = allowTimeSum
				}
				atomic.AddInt64(&allowTimeSum, allowTime)
				if err != nil {
					fmt.Printf("not allow: %s\n", err)
					atomic.AddInt64(&failCount, 1)
				} else {
					fmt.Println("allow")
					atomic.AddInt64(&successCount, 1)
				}
				wg.Done()
			}()
			time.Sleep(100 * time.Millisecond)
		}
		wg.Done()
	}

	start := time.Now()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go fn()
	}

	wg.Wait()
	timeMillisecond := time.Now().Sub(start).Milliseconds()
	fmt.Printf("总耗时（预期5000ms）：%d\n", timeMillisecond)
	fmt.Printf("allow处理总耗时：%d, 单次allow平均耗时：%f，最大allow时间：%d，最小allow时间：%d\n", allowTimeSum, float64(allowTimeSum)/250, maxAllowTime, minAllowTime)
	//fmt.Printf("EtcdGet总时间（ms）：%d，总次数：%d，平均时间：%f\n", limit.EtcdGetTime, limit.EtcdGetNum, float64(limit.EtcdGetTime)/float64(limit.EtcdGetNum))
	//fmt.Printf("EtcdSet总时间（ms）：%d，总次数：%d，平均时间：%f\n", limit.EtcdSetTime, limit.EtcdSetNum, float64(limit.EtcdSetTime)/float64(limit.EtcdSetNum))
	//fmt.Printf("EtcdPut总时间（ms）：%d，总次数：%d，平均时间：%f\n", limit.EtcdPutTime, limit.EtcdPutNum, float64(limit.EtcdPutTime)/float64(limit.EtcdPutNum))
	fmt.Printf("RedisLuaRun总时间（ms）：%d，总次数：%d，平均时间：%f\n", limit.RedisLuaRunTime, limit.RedisLuaRunNum, float64(limit.RedisLuaRunTime)/float64(limit.RedisLuaRunNum))
	fmt.Printf("总允许：%d，总拒绝：%d，平均每秒允许：%f，每秒拒绝：%f", successCount, failCount, float64(successCount)*1000/float64(timeMillisecond), float64(failCount)*1000/float64(timeMillisecond))
}
