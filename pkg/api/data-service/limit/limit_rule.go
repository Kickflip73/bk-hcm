package limit

import (
	"hcm/pkg/api/core"
	"hcm/pkg/criteria/validator"
	"hcm/pkg/dal/table/types"
	"hcm/pkg/rest"
	"hcm/pkg/runtime/filter"
)

// LimitRuleReq
type LimitRuleListReq struct {
	Filter *filter.Expression `json:"filter" validate:"required"`
	Page   *core.BasePage     `json:"page" validate:"required"`
}

// Validate ...
func (l *LimitRuleListReq) Validate() error {
	return validator.Validate.Struct(l)
}

// LimitRuleListResp
type LimitRuleListResp struct {
	rest.BaseResp `json:",inline"`
	Data          *LimitRuleListResult `json:"data"`
}

// LimitRuleListResult
type LimitRuleListResult struct {
	Count   uint64           `json:"count"`
	Details []*BaseLimitRule `json:"details"`
}

// BaseLimitRule
type BaseLimitRule struct {
	ID              int64      `json:"id"`
	RuleName        string     `json:"rule_name"`
	Scene           string     `json:"scene"`
	Account         string     `json:"account"`
	Identify        string     `json:"identify"`
	MaxLimit        int64      `json:"max_limit"`
	WindowSize      int        `json:"windows_size"`
	RejectPolicy    string     `json:"reject_policy"`
	RetryInterval   int        `json:"retry_interval"`
	RetryMaxTimeout int        `json:"retry_max_timeout"`
	RetryMaxCount   int        `json:"retry_max_count"`
	DenyAll         bool       `json:"deny_all"`
	Enabled         bool       `json:"enabled"`
	Creator         string     `json:"creator"`
	Reviser         string     `json:"reviser"`
	CreatedAt       types.Time `json:"created_at"`
	UpdatedAt       types.Time `json:"updated_at"`
}
