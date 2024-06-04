package limit

import (
	"fmt"
	protolimit "hcm/pkg/api/data-service/limit"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/dal/dao/types"
	"hcm/pkg/logs"
	"hcm/pkg/rest"
)

// ListLimitRule query LimitRule list
func (svc *service) ListLimitRule(cts *rest.Contexts) (interface{}, error) {
	req := new(protolimit.LimitRuleListReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, err
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	opt := &types.ListOption{
		Filter: req.Filter,
		Page:   req.Page,
	}
	daoLimitRuleResp, err := svc.dao.LimitRule().List(cts.Kit, opt)
	if err != nil {
		logs.Errorf("list limitRule failed, err: %v, rid: %s", err, cts.Kit.Rid)
		return nil, fmt.Errorf("list limitRule failed, err: %v", err)
	}

	return daoLimitRuleResp, nil
}
