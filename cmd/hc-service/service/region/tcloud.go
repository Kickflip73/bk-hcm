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

// Package region defines region service.
package region

import (
	dataservice "hcm/pkg/api/data-service"
	protoDsRegion "hcm/pkg/api/data-service/cloud/region"
	protoHcRegion "hcm/pkg/api/hc-service/region"
	"hcm/pkg/criteria/constant"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/dal/dao/tools"
	"hcm/pkg/logs"
	"hcm/pkg/rest"
)

// TCloudRegionSync sync tcloud region.
func (r region) TCloudRegionSync(cts *rest.Contexts, vendor enumor.Vendor) (interface{}, error) {
	req := new(protoHcRegion.TCloudRegionSyncReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	cli, err := r.ad.TCloud(cts.Kit, req.AccountID)
	if err != nil {
		logs.Errorf("RegionSyncTcloud:ad.TCloud:Err, accountID: %s, err: %v, rid: %s",
			req.AccountID, err, cts.Kit.Rid)
		return nil, err
	}

	cloudResp, err := cli.ListRegion(cts.Kit)
	if err != nil {
		logs.Errorf("RegionSyncTcloud:cli.ListRegion:Err, err: %v, rid: %s", err, cts.Kit.Rid)
		return nil, err
	}

	tmpRegions := make([]protoDsRegion.RegionBatchCreate, 0)
	for _, item := range cloudResp.Details {
		// 不可用的地区，不录入
		if item.RegionState != constant.TCloudAvailbleState {
			continue
		}

		tmpRegions = append(tmpRegions, protoDsRegion.RegionBatchCreate{
			Vendor:     vendor,
			RegionID:   item.RegionID,
			RegionName: item.RegionName,
		})
	}

	if len(tmpRegions) == 0 {
		return nil, errf.New(errf.RecordNotFound, "cloudapi has not available region")
	}

	// batch forbidden tcloud region state.
	updateStateReq := &protoDsRegion.RegionBatchUpdateReq{
		Regions: []protoDsRegion.RegionBatchUpdate{{IsAvailable: constant.AvailableNo}},
	}
	err = r.cs.DataService().TCloud.Region.BatchForbiddenRegionState(cts.Kit.Ctx, cts.Kit.Header(), updateStateReq)
	if err != nil {
		logs.Errorf("RegionSyncTcloud:BatchForbiddenRegionState:Err, err: %v, rid: %s", err, cts.Kit.Rid)
		return nil, err
	}

	// batch create tcloud region.
	createReq := &protoDsRegion.RegionCreateReq{
		Regions: tmpRegions,
	}
	resp, err := r.cs.DataService().TCloud.Region.BatchCreate(cts.Kit.Ctx, cts.Kit.Header(), createReq)
	if err != nil {
		logs.Errorf("RegionSyncTcloud:BatchCreate:Err, err: %v, rid: %s", err, cts.Kit.Rid)
		return nil, err
	}

	// batch delete tcloud region.
	deleteReq := &dataservice.BatchDeleteReq{
		Filter: tools.EqualExpression("is_available", constant.AvailableNo),
	}
	err = r.cs.DataService().TCloud.Region.BatchDelete(cts.Kit.Ctx, cts.Kit.Header(), deleteReq)
	if err != nil {
		logs.Errorf("RegionSyncTcloud:BatchDelete:Err, err: %v, rid: %s", err, cts.Kit.Rid)
		return nil, err
	}

	return resp, nil
}
