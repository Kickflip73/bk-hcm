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

package cvm

import (
	"hcm/cmd/hc-service/logics/sync/cvm"
	typecore "hcm/pkg/adaptor/types/core"
	typecvm "hcm/pkg/adaptor/types/cvm"
	"hcm/pkg/api/hc-service/sync"
	"hcm/pkg/criteria/constant"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/logs"
	"hcm/pkg/rest"
)

// SyncHuaWeiCvm ...
func (svc *syncCvmSvc) SyncHuaWeiCvm(cts *rest.Contexts) (interface{}, error) {
	req := new(sync.HuaWeiSyncReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	syncOpt := &cvm.SyncHuaWeiCvmOption{
		AccountID: req.AccountID,
		Region:    req.Region,
	}

	_, err := cvm.SyncHuaWeiCvm(cts.Kit, syncOpt, svc.adaptor, svc.dataCli)
	if err != nil {
		logs.Errorf("request to sync huawei cvm all rel failed, err: %v, rid: %s", err, cts.Kit.Rid)
		return nil, err
	}

	return nil, nil
}

// SyncHuaWeiCvmWithRelResource ...
func (svc *syncCvmSvc) SyncHuaWeiCvmWithRelResource(cts *rest.Contexts) (interface{}, error) {
	req := new(sync.HuaWeiSyncReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	cli, err := svc.adaptor.HuaWei(cts.Kit, req.AccountID)
	if err != nil {
		return nil, err
	}

	listOpt := &typecvm.HuaWeiListOption{
		Region:   req.Region,
		CloudIDs: nil,
		Page: &typecore.HuaWeiCvmOffsetPage{
			Offset: int32(0),
			Limit:  int32(constant.BatchOperationMaxLimit),
		},
	}
	for {
		cvms, err := cli.ListCvm(cts.Kit, listOpt)
		if err != nil {
			logs.Errorf("request adaptor list huawei cvm failed, err: %v, opt: %v, rid: %s", err, listOpt, cts.Kit.Rid)
			return nil, err
		}

		if cvms == nil || len(*cvms) == 0 {
			break
		}

		cloudIDs := make([]string, 0, len(*cvms))
		for _, one := range *cvms {
			cloudIDs = append(cloudIDs, one.Id)
		}

		syncOpt := &cvm.SyncHuaWeiCvmOption{
			AccountID: req.AccountID,
			Region:    req.Region,
			CloudIDs:  cloudIDs,
		}

		_, err = cvm.SyncHuaWeiCvmWithRelResource(cts.Kit, syncOpt, svc.adaptor, svc.dataCli)
		if err != nil {
			logs.Errorf("request to sync huawei cvm all rel failed, err: %v, rid: %s", err, cts.Kit.Rid)
			return nil, err
		}

		if len(*cvms) < typecore.TCloudQueryLimit {
			break
		}

		listOpt.Page.Offset += typecore.TCloudQueryLimit
	}

	return nil, nil
}
