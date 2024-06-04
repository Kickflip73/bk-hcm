package global

import (
	"context"
	"hcm/pkg/api/data-service/limit"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/rest"
	"net/http"
)

// LimitRuleClient is data service LimitRuleClient api client.
type LimitRuleClient struct {
	client rest.ClientInterface
}

// NewLimitRuleClient create a new LimitRule api client.
func NewLimitRuleClient(client rest.ClientInterface) *LimitRuleClient {
	return &LimitRuleClient{
		client: client,
	}
}

// List ...
func (a *LimitRuleClient) List(ctx context.Context, h http.Header, request *limit.LimitRuleListReq) (
	*limit.LimitRuleListResult, error,
) {
	resp := new(limit.LimitRuleListResp)

	err := a.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef("/limit_rule/list").
		WithHeaders(h).
		Do().
		Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != errf.OK {
		return nil, errf.New(resp.Code, resp.Message)
	}

	return resp.Data, nil
}
