package limit

import "hcm/pkg/dal/table/limit"

// LimitRuleResult is the result of limit rule
type LimitRuleResult struct {
	Count   uint64
	Details []*limit.LimitRuleModel
}

// LimitRuleCountResult DiskCountResult is the result of disk count
type LimitRuleCountResult struct {
	Count uint64
}
