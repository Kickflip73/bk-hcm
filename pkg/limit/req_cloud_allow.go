package limit

import (
	"context"
	"fmt"
	etcd3 "go.etcd.io/etcd/client/v3"
	"hcm/pkg/logs"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// CheckReqCloudAllow Check the current limit on the request cloud service side
func (l *Limiter) CheckReqCloudAllow(accountId string, identify string) error {
	// 判断限流器是否已经准备好，各client是否已经初始化
	err := l.checkReady()
	if err != nil {
		return err
	}

	l.lock.RLock()
	defer l.lock.RUnlock()

	// TODO 匹配限流规则
	mapKey := fmt.Sprintf("%s.%s", accountId, identify)
	rule := l.reqCloudRules[mapKey]
	if rule == nil {
		logs.Infof("no matching current limiting rules found：%s", mapKey)
		return nil
	}

	// TODO 判断是否限流
	if rule.DenyAll {
		return fmt.Errorf("此接口已拒绝所有请求")
	}

	key := fmt.Sprintf("%s.%s", ReqCloudLimitRulePrefix, mapKey)
	allow, allowErr := l.allowByRedis(key, rule.WindowSize, rule.MaxLimit)
	//allow, allowErr := l.allowByEtcd(key, rule.WindowSize, rule.MaxLimit)
	if allowErr != nil {
		return allowErr
	}
	if allow {
		return nil
	}

	switch rule.RejectPolicy {
	case Direct:
		return fmt.Errorf("此接口流量已被限制")
	case RetryTimeout:
		start := time.Now()
		for time.Now().Sub(start).Milliseconds() < int64(rule.RetryMaxTimeout) {
			allow, allowErr := l.allowByRedis(key, rule.WindowSize, rule.MaxLimit)
			//allow, allowErr := l.allowByEtcd(key, rule.WindowSize, rule.MaxLimit)
			if allowErr != nil {
				return allowErr
			}
			if allow {
				return nil
			}
			retryInterval := time.Duration(rule.RetryInterval) * time.Millisecond
			time.Sleep(retryInterval)
		}
	case RetryCount:
		for i := 1; i <= rule.RetryMaxCount; i++ {
			allow, allowErr := l.allowByRedis(key, rule.WindowSize, rule.MaxLimit)
			//allow, allowErr := l.allowByEtcd(key, rule.WindowSize, rule.MaxLimit)
			if allowErr != nil {
				return allowErr
			}
			if allow {
				return nil
			}
			retryInterval := time.Duration(rule.RetryInterval) * time.Millisecond
			time.Sleep(retryInterval)
		}
	default:
		return fmt.Errorf("current limiting rule rejection policy is an illegal value: %s", rule.RejectPolicy)
	}

	return fmt.Errorf("此接口流量已被限制")
}

// RedisLuaRunTime int64 = 0
var RedisLuaRunTime int64 = 0

// RedisLuaRunNum int64 = 0
var RedisLuaRunNum int64 = 0

func (l *Limiter) allowByRedis(key string, windowsSize int, maxLimit int64) (bool, error) {
	// 原子的，串行的执行lua脚本
	start := time.Now()
	result, err := l.redisClient.Eval(context.Background(), checkAllowScript, []string{key}, windowsSize).Result()
	if err != nil {
		logs.Errorf("执行限流脚本出错，%s", err)
		return false, err
	}

	count, ok := result.(int64)
	if !ok {
		logs.Errorf("执行限流脚本出错，key：%s, result：%s", key, result)
		return false, fmt.Errorf("执行限流脚本出错，key：%s, result：%s", key, result)
	}
	atomic.AddInt64(&RedisLuaRunTime, time.Now().Sub(start).Milliseconds())
	atomic.AddInt64(&RedisLuaRunNum, 1)

	if count > maxLimit {
		return false, nil
	}
	return true, nil
}

const checkAllowScript = `
local cnt = redis.pcall('INCR', KEYS[1]);
if type(cnt) ~= "number"
then
	return cnt
end

local rs = redis.pcall('TTL', KEYS[1]);
if type(rs) ~= "number"
then
	return rs
end

if rs == -1
then
	rs = redis.pcall('EXPIRE', KEYS[1], ARGV[1]);
	if type(rs) ~= "number"
	then
		return rs
	end
end

return cnt
`

var locked sync.Mutex

// EtcdGetTime
var EtcdGetTime int64 = 0

// EtcdGetNum
var EtcdGetNum int64 = 0

// EtcdSetTime
var EtcdSetTime int64 = 0

// EtcdSetNum
var EtcdSetNum int64 = 0

// EtcdPutTime
var EtcdPutTime int64 = 0

// EtcdPutNum
var EtcdPutNum int64 = 0

func (l *Limiter) allowByEtcd(key string, windowsSize int, maxLimit int64) (bool, error) {
	// 模拟分布式锁，串行执行所有请求
	locked.Lock()
	defer locked.Unlock()
	var count int64
	//for {
	start := time.Now()
	ctx := context.Background()
	resp, err := l.etcdClient.Get(ctx, key)
	if err != nil {
		return false, err
	}
	atomic.AddInt64(&EtcdGetTime, time.Now().Sub(start).Milliseconds())
	atomic.AddInt64(&EtcdGetNum, 1)

	if len(resp.Kvs) == 0 || resp.Kvs[0].Lease == 0 {
		//// 暂时模拟一个分布式锁，用来防止多个线程同时构建新窗口
		//if !locked.TryLock() {
		//	continue
		//}
		//defer locked.Unlock()
		start = time.Now()
		lease, _ = l.etcdClient.Grant(ctx, int64(windowsSize))
		fmt.Println("创建新窗口")
		_, err = l.etcdClient.Put(ctx, key, "1", etcd3.WithLease(lease.ID))
		if err != nil {
			return false, err
		}
		count = 1
		atomic.AddInt64(&EtcdSetTime, time.Now().Sub(start).Milliseconds())
		atomic.AddInt64(&EtcdSetNum, 1)
		//break
	} else {
		start = time.Now()
		currentValue, err := strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64)
		if err != nil {
			return false, err
		}
		count = currentValue + 1
		_, err = l.etcdClient.Put(ctx, key, strconv.FormatInt(count, 10), etcd3.WithLease(lease.ID))
		if err != nil {
			return false, err
		}
		atomic.AddInt64(&EtcdPutTime, time.Now().Sub(start).Milliseconds())
		atomic.AddInt64(&EtcdPutNum, 1)

		//// 乐观锁尝试自增请求数
		//txnResp, err := etcdClient.Txn(ctx).
		//	If(etcd3.Compare(etcd3.Version(key), "=", resp.Kvs[0].Version)).
		//	Then(etcd3.OpPut(key, strconv.FormatInt(count, 10), etcd3.WithLease(lease.ID))).
		//	Commit()
		//if err != nil {
		//	return false, err
		//}
		//if txnResp.Succeeded {
		//	break
		//}

		//// 更新失败，自旋重试
		//time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
	//}

	if count > maxLimit {
		return false, nil
	}
	return true, nil
}

const (
	// Direct direct Rejection
	Direct string = "direct"
	// RetryTimeout retry within the specified time
	RetryTimeout string = "retry_timeout"
	// RetryCount retry within the specified number of time
	RetryCount string = "retry_count"
)

const (
	// ReqCloudLimitRulePrefix prefix of request cloud limit rule
	ReqCloudLimitRulePrefix string = "limit.req_cloud"
	// UserReqLimitRulePrefix prefix of user request limit rule
	UserReqLimitRulePrefix string = "limit.user_req"
)
