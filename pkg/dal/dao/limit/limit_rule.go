package limit

import (
	"fmt"
	"hcm/pkg/api/core"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/dal/dao/audit"
	idgenerator "hcm/pkg/dal/dao/id-generator"
	"hcm/pkg/dal/dao/orm"
	"hcm/pkg/dal/dao/tools"
	"hcm/pkg/dal/dao/types"
	"hcm/pkg/dal/dao/types/limit"
	"hcm/pkg/dal/table"
	limit_table "hcm/pkg/dal/table/limit"
	"hcm/pkg/kit"
	"hcm/pkg/logs"
	"hcm/pkg/runtime/filter"
)

// LimitRule is the interface for limit rule
type LimitRule interface {
	List(kt *kit.Kit, opt *types.ListOption) (*limit.LimitRuleResult, error)
}

var _ LimitRule = new(LimitRuleDao)

// LimitRuleDao is the dao for limit rule
type LimitRuleDao struct {
	Orm   orm.Interface
	IDGen idgenerator.IDGenInterface
	Audit audit.Interface
}

// List is used to list limit rule
func (limitRuleDao LimitRuleDao) List(kt *kit.Kit, opt *types.ListOption) (*limit.LimitRuleResult, error) {
	if opt == nil {
		return nil, errf.New(errf.InvalidParameter, "limitRule options is nil")
	}

	columnTypes := limit_table.LimitRuleColumns.ColumnTypes()
	columnTypes["extension.resource_group_name"] = enumor.String
	columnTypes["extension.self_link"] = enumor.String
	columnTypes["extension.zones"] = enumor.Json
	if err := opt.Validate(filter.NewExprOption(filter.RuleFields(columnTypes)),
		core.NewDefaultPageOption()); err != nil {
		return nil, err
	}

	whereOpt := tools.DefaultSqlWhereOption
	whereExpr, whereValue, err := opt.Filter.SQLWhereExpr(whereOpt)
	if err != nil {
		return nil, err
	}

	if opt.Page.Count {
		// this is a count request, then do count operation only.
		sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, table.LimitRuleTable, whereExpr)
		count, err := limitRuleDao.Orm.Do().Count(kt.Ctx, sql, whereValue)
		if err != nil {
			logs.ErrorJson("count limitRule failed, err: %v, filter: %s, rid: %s", err, opt.Filter, kt.Rid)
			return nil, err
		}
		return &limit.LimitRuleResult{Count: count}, nil
	}
	pageExpr, err := types.PageSQLExpr(opt.Page, types.DefaultPageSQLOption)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(`SELECT %s FROM %s %s %s`, limit_table.LimitRuleColumns.FieldsNamedExpr(opt.Fields), table.LimitRuleTable,
		whereExpr, pageExpr)

	details := make([]*limit_table.LimitRuleModel, 0)
	if err = limitRuleDao.Orm.Do().Select(kt.Ctx, &details, sql, whereValue); err != nil {
		return nil, err
	}

	result := &limit.LimitRuleResult{Details: details}

	return result, nil
}
