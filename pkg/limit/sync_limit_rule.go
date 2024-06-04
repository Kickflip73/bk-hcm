package limit

import (
	"fmt"
	"hcm/pkg/api/core"
	dataproto "hcm/pkg/api/data-service/limit"
	"hcm/pkg/kit"
	"hcm/pkg/logs"
	"hcm/pkg/runtime/filter"
	"time"
)

// SyncLimiterRules sync the api limiter rules from mysql
func (l *Limiter) SyncLimiterRules(kit *kit.Kit) error {
	logs.Infof("begin SyncLimiterRules")
	var i int
	var err error
	for i = 1; i <= l.syncMaxTrys; i++ {
		err = l.syncLimiterRules(kit)
		if err == nil {
			logs.Infof("SyncLimiterRules successful")
			return nil
		}
		logs.Errorf("fail to syncLimiterRules, err: %s, retries: %d/%d", err, i, l.syncMaxTrys)
		time.Sleep(l.syncDuration)
	}

	return fmt.Errorf("fail to syncLimiterRules, err:%s", err)
}

func (l *Limiter) syncLimiterRules(kit *kit.Kit) error {
	if l.dataServiceClient == nil {
		return fmt.Errorf("dataServiceClient has not been set up yet")
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	// call data-service api：/api/v1/data/limit-rule/list to get limiter rules
	result, err := l.dataServiceClient.Global.LimitRule.List(
		kit.Ctx,
		kit.Header(),
		limitRuleListReq(),
	)
	if err != nil {
		return err
	}

	for _, limitRule := range result.Details {
		switch limitRule.Scene {
		case UserRequest:
			mapKey := fmt.Sprintf("user_req_map_key")
			l.userReqRules[mapKey] = limitRule
		case RequestCloud:
			mapKey := fmt.Sprintf("%s.%s", limitRule.Account, limitRule.Identify)
			l.reqCloudRules[mapKey] = limitRule
		default:
			logs.Errorf("scene of limiting rule is an illegal value: %s", limitRule.Scene)
		}
	}

	return nil
}

// limitRuleListReq return expression for querying the enabled limiting rule list
func limitRuleListReq() *dataproto.LimitRuleListReq {
	return &dataproto.LimitRuleListReq{
		Filter: &filter.Expression{
			Op: filter.And,
			Rules: []filter.RuleFactory{
				&filter.AtomRule{
					Field: "enabled",
					Op:    filter.OpFactory(filter.Equal),
					Value: "1",
				},
			},
		},
		Page: &core.BasePage{
			Count: false,
			Start: 0,
			Limit: 500,
		},
	}
}

const (
	// UserRequest User requested scenario
	UserRequest string = "user_req"
	// RequestCloud Requesting cloud services scenario
	RequestCloud string = "req_cloud"
)
