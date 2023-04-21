/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 混合云管理平台 (BlueKing - Hybrid Cloud Management System) available.
 * Copyright (C) 2022 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 *
 * to the current version of the project delivered to anyone in the future.
 */

// Package subnet defines subnet service.
package subnet

import (
	syncsubnet "hcm/cmd/hc-service/logics/sync/subnet"
	"hcm/pkg/adaptor/types"
	adcore "hcm/pkg/adaptor/types/core"
	"hcm/pkg/api/core"
	dataservice "hcm/pkg/api/data-service"
	"hcm/pkg/api/data-service/cloud"
	hcservice "hcm/pkg/api/hc-service"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/dal/dao/tools"
	"hcm/pkg/logs"
	"hcm/pkg/rest"
)

// HuaWeiSubnetCreate create huawei subnet.
func (s subnet) HuaWeiSubnetCreate(cts *rest.Contexts) (interface{}, error) {
	req := new(hcservice.SubnetCreateReq[hcservice.HuaWeiSubnetCreateExt])
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}
	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	cli, err := s.ad.HuaWei(cts.Kit, req.AccountID)
	if err != nil {
		return nil, err
	}

	huaweiCreateOpt := &types.HuaWeiSubnetCreateOption{
		Name:       req.Name,
		Memo:       req.Memo,
		CloudVpcID: req.CloudVpcID,
		Extension: &types.HuaWeiSubnetCreateExt{
			Region:     req.Extension.Region,
			Zone:       req.Extension.Zone,
			IPv4Cidr:   req.Extension.IPv4Cidr,
			Ipv6Enable: req.Extension.Ipv6Enable,
			GatewayIp:  req.Extension.GatewayIp,
		},
	}
	huaweiCreateRes, err := cli.CreateSubnet(cts.Kit, huaweiCreateOpt)
	if err != nil {
		return nil, err
	}

	// create hcm subnet
	syncOpt := &syncsubnet.SyncHuaWeiOption{
		AccountID:  req.AccountID,
		Region:     req.Extension.Region,
		CloudVpcID: req.CloudVpcID,
	}
	createReqs := []cloud.SubnetCreateReq[cloud.HuaWeiSubnetCreateExt]{convertHuaWeiSubnetCreateReq(huaweiCreateRes,
		req.AccountID, req.BkBizID)}
	res, err := syncsubnet.BatchCreateHuaWeiSubnet(cts.Kit, createReqs, s.cs.DataService(), s.ad, syncOpt)
	if err != nil {
		logs.Errorf("sync huawei subnet failed, err: %v, reqs: %+v, rid: %s", err, createReqs, cts.Kit.Rid)
		return nil, err
	}

	return core.CreateResult{ID: res.IDs[0]}, nil
}

func convertHuaWeiSubnetCreateReq(data *types.HuaWeiSubnet, accountID string,
	bizID int64) cloud.SubnetCreateReq[cloud.HuaWeiSubnetCreateExt] {

	subnetReq := cloud.SubnetCreateReq[cloud.HuaWeiSubnetCreateExt]{
		AccountID:  accountID,
		CloudVpcID: data.CloudVpcID,
		CloudID:    data.CloudID,
		Name:       &data.Name,
		Region:     data.Extension.Region,
		Ipv4Cidr:   data.Ipv4Cidr,
		Ipv6Cidr:   data.Ipv6Cidr,
		Memo:       data.Memo,
		BkBizID:    bizID,
		Extension: &cloud.HuaWeiSubnetCreateExt{
			Status:       data.Extension.Status,
			DhcpEnable:   data.Extension.DhcpEnable,
			GatewayIp:    data.Extension.GatewayIp,
			DnsList:      data.Extension.DnsList,
			NtpAddresses: data.Extension.NtpAddresses,
		},
	}

	return subnetReq
}

// HuaWeiSubnetUpdate update huawei subnet.
func (s subnet) HuaWeiSubnetUpdate(cts *rest.Contexts) (interface{}, error) {
	id := cts.PathParameter("id").String()

	req := new(hcservice.SubnetUpdateReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}
	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	getRes, err := s.cs.DataService().HuaWei.Subnet.Get(cts.Kit.Ctx, cts.Kit.Header(), id)
	if err != nil {
		return nil, err
	}

	cli, err := s.ad.HuaWei(cts.Kit, getRes.AccountID)
	if err != nil {
		return nil, err
	}

	updateOpt := &types.HuaWeiSubnetUpdateOption{
		SubnetUpdateOption: types.SubnetUpdateOption{
			ResourceID: getRes.CloudID,
			Data:       &types.BaseSubnetUpdateData{Memo: req.Memo},
		},
		Name:   getRes.Name,
		Region: getRes.Region,
		VpcID:  getRes.CloudVpcID,
	}
	err = cli.UpdateSubnet(cts.Kit, updateOpt)
	if err != nil {
		return nil, err
	}

	updateReq := &cloud.SubnetBatchUpdateReq[cloud.HuaWeiSubnetUpdateExt]{
		Subnets: []cloud.SubnetUpdateReq[cloud.HuaWeiSubnetUpdateExt]{{
			ID: id,
			SubnetUpdateBaseInfo: cloud.SubnetUpdateBaseInfo{
				Memo: req.Memo,
			},
		}},
	}
	err = s.cs.DataService().HuaWei.Subnet.BatchUpdate(cts.Kit.Ctx, cts.Kit.Header(), updateReq)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// HuaWeiSubnetDelete delete huawei subnet.
func (s subnet) HuaWeiSubnetDelete(cts *rest.Contexts) (interface{}, error) {
	id := cts.PathParameter("id").String()

	getRes, err := s.cs.DataService().HuaWei.Subnet.Get(cts.Kit.Ctx, cts.Kit.Header(), id)
	if err != nil {
		return nil, err
	}

	cli, err := s.ad.HuaWei(cts.Kit, getRes.AccountID)
	if err != nil {
		return nil, err
	}

	delOpt := &types.HuaWeiSubnetDeleteOption{
		BaseRegionalDeleteOption: adcore.BaseRegionalDeleteOption{
			BaseDeleteOption: adcore.BaseDeleteOption{ResourceID: getRes.CloudID},
			Region:           getRes.Region,
		},
		VpcID: getRes.CloudVpcID,
	}
	err = cli.DeleteSubnet(cts.Kit, delOpt)
	if err != nil {
		return nil, err
	}

	deleteReq := &dataservice.BatchDeleteReq{
		Filter: tools.EqualExpression("id", id),
	}
	err = s.cs.DataService().Global.Subnet.BatchDelete(cts.Kit.Ctx, cts.Kit.Header(), deleteReq)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// HuaWeiSubnetCountIP count huawei subnets' available ips.
func (s subnet) HuaWeiSubnetCountIP(cts *rest.Contexts) (interface{}, error) {
	id := cts.PathParameter("id").String()
	if len(id) == 0 {
		return nil, errf.New(errf.InvalidParameter, "id is required")
	}

	getRes, err := s.cs.DataService().HuaWei.Subnet.Get(cts.Kit.Ctx, cts.Kit.Header(), id)
	if err != nil {
		return nil, err
	}

	cli, err := s.ad.HuaWei(cts.Kit, getRes.AccountID)
	if err != nil {
		return nil, err
	}

	listOpt := &types.HuaWeiVpcIPAvailGetOption{
		Region:   getRes.Region,
		SubnetID: getRes.CloudID,
	}
	availabilities, err := cli.GetSubnetIPAvailabilities(cts.Kit, listOpt)
	if err != nil {
		return nil, err
	}

	return &hcservice.SubnetCountIPResult{
		AvailableIPv4Count: uint64(availabilities.TotalIps - availabilities.UsedIps),
	}, nil
}
