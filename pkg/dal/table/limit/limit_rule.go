package limit

import (
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/dal/table/types"
	"hcm/pkg/dal/table/utils"
)

// LimitRuleColumns represents the columns of the "limiting_rule" table.
var LimitRuleColumns = utils.MergeColumns(nil, LimitRuleColumnsDescriptor)

// LimitRuleColumnsDescriptor represents the columns of the "limiting_rule" table.
var LimitRuleColumnsDescriptor = utils.ColumnDescriptors{
	{Column: "id", NamedC: "id", Type: enumor.String},
	{Column: "rule_name", NamedC: "rule_name", Type: enumor.String},
	{Column: "scene", NamedC: "scene", Type: enumor.String},
	{Column: "account", NamedC: "account", Type: enumor.String},
	{Column: "identify", NamedC: "identify", Type: enumor.String},
	{Column: "max_limit", NamedC: "max_limit", Type: enumor.Numeric},
	{Column: "windows_size", NamedC: "windows_size", Type: enumor.Numeric},
	{Column: "reject_policy", NamedC: "reject_policy", Type: enumor.String},
	{Column: "retry_interval", NamedC: "retry_interval", Type: enumor.Numeric},
	{Column: "retry_max_timeout", NamedC: "retry_max_timeout", Type: enumor.Numeric},
	{Column: "retry_max_count", NamedC: "retry_max_count", Type: enumor.Numeric},
	{Column: "deny_all", NamedC: "deny_all", Type: enumor.Boolean},
	{Column: "enabled", NamedC: "enabled", Type: enumor.String},
	{Column: "creator", NamedC: "creator", Type: enumor.String},
	{Column: "reviser", NamedC: "reviser", Type: enumor.String},
	{Column: "created_at", NamedC: "created_at", Type: enumor.Time},
	{Column: "updated_at", NamedC: "updated_at", Type: enumor.Time},
}

// LimitRuleModel represents the limiting_rule table storing the configuration of rate limiting rules.
type LimitRuleModel struct {
	ID              int64      `db:"id" json:"id"`
	RuleName        string     `db:"rule_name" json:"rule_name"`
	Scene           string     `db:"scene" json:"scene"`
	Account         string     `db:"account" json:"account"`
	Identify        string     `db:"identify" json:"identify"`
	MaxLimit        int64      `db:"max_limit" json:"max_limit"`
	WindowSize      int        `db:"windows_size" json:"windows_size"`
	RejectPolicy    string     `db:"reject_policy" json:"reject_policy"`
	RetryInterval   int        `db:"retry_interval" json:"retry_interval"`
	RetryMaxTimeout int        `db:"retry_max_timeout" json:"retry_max_timeout"`
	RetryMaxCount   int        `db:"retry_max_count" json:"retry_max_count"`
	DenyAll         bool       `db:"deny_all" json:"deny_all"`
	Enabled         bool       `db:"enabled" json:"enabled"`
	Creator         string     `db:"creator" json:"creator"`
	Reviser         string     `db:"reviser" json:"reviser"`
	CreatedAt       types.Time `db:"created_at" json:"created_at"`
	UpdatedAt       types.Time `db:"updated_at" json:"updated_at"`
}
