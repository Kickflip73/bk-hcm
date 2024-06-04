package limit

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	etcd3 "go.etcd.io/etcd/client/v3"
	"hcm/pkg/api/core"
	dataproto "hcm/pkg/api/data-service/limit"
	dataservice "hcm/pkg/client/data-service"
	"hcm/pkg/logs"
	"sync"
	"time"
)

var (
	limiter *Limiter
	once    sync.Once
	lease   *etcd3.LeaseGrantResponse
)

// Limiter is used to rate-limit some process
type Limiter struct {
	userReqRules      map[string]*dataproto.BaseLimitRule
	reqCloudRules     map[string]*dataproto.BaseLimitRule
	dataServiceClient *dataservice.Client
	redisClient       *redis.Client
	etcdClient        *etcd3.Client
	syncMaxTrys       int
	syncDuration      time.Duration
	lock              sync.RWMutex
}

// GetLimiter  return limiter
func GetLimiter() *Limiter {
	return limiter
}

// InitLimiterWithSyncByMysql init Limiter and sync limit rules by mysql
func InitLimiterWithSyncByMysql(dataServiceClient *dataservice.Client, redisClient *redis.Client, etcdClient *etcd3.Client) *Limiter {
	once.Do(func() {
		initLimiterWithSyncByMysql(dataServiceClient, redisClient, etcdClient)
	})
	return limiter
}

// InitLimiterWithRules init Limiter and set limit rules
func InitLimiterWithRules(rules map[string]*dataproto.BaseLimitRule, redisClient *redis.Client, etcdClient *etcd3.Client) *Limiter {
	once.Do(func() {
		initLimiterWithRules(rules, redisClient, etcdClient)
	})
	return limiter
}

func initLimiterWithSyncByMysql(dataServiceClient *dataservice.Client, redisClient *redis.Client, etcdClient *etcd3.Client) {
	// 初始化Limiter
	logs.Infof("开始初始化Limiter")
	limiter = new(Limiter)
	limiter.dataServiceClient = dataServiceClient
	limiter.redisClient = redisClient
	limiter.etcdClient = etcdClient
	limiter.userReqRules = make(map[string]*dataproto.BaseLimitRule)
	limiter.reqCloudRules = make(map[string]*dataproto.BaseLimitRule)
	limiter.syncMaxTrys = 5
	limiter.syncDuration = 1 * time.Second

	// 从数据库同步限流规则
	backendKit := core.NewBackendKit()
	err := limiter.SyncLimiterRules(backendKit)
	if err != nil {
		return
	}
}

func initLimiterWithRules(rules map[string]*dataproto.BaseLimitRule, redisClient *redis.Client, etcdClient *etcd3.Client) {
	// 初始化Limiter
	logs.Infof("开始初始化Limiter")
	limiter = new(Limiter)
	limiter.redisClient = redisClient
	limiter.etcdClient = etcdClient
	limiter.userReqRules = make(map[string]*dataproto.BaseLimitRule)
	limiter.reqCloudRules = make(map[string]*dataproto.BaseLimitRule)
	limiter.syncMaxTrys = 5
	limiter.syncDuration = 1 * time.Second

	// 从数据库同步限流规则
	for _, rule := range rules {
		if rule.Scene == UserRequest {
			mapKey := fmt.Sprintf("user_req_map_key")
			limiter.userReqRules[mapKey] = rule
		} else {
			mapKey := fmt.Sprintf("%s.%s", rule.Account, rule.Identify)
			limiter.reqCloudRules[mapKey] = rule
		}
	}
}

func (l *Limiter) checkReady() error {
	if l.redisClient == nil {
		return fmt.Errorf("limiter is not fully initialized，redisClient is nil")
	}

	if l.etcdClient == nil {
		return fmt.Errorf("limiter is not fully initialized，etcdClient is nil")
	}

	return nil
}
