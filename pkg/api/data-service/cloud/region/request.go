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

package region

import (
	"errors"
	"fmt"
	"time"

	"hcm/pkg/api/core"
	"hcm/pkg/criteria/constant"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/validator"
	"hcm/pkg/rest"
	"hcm/pkg/runtime/filter"
)

// -------------------------- Create --------------------------

// RegionCreateReq define region create request.
type RegionCreateReq struct {
	Regions   []RegionBatchCreate `json:"regions" validate:"required"`
	AccountID string              `json:"account_id" validate:"omitempty"`
}

// RegionBatchCreate define region rule when create.
type RegionBatchCreate struct {
	Vendor     enumor.Vendor `json:"vendor" validate:"required"`
	RegionID   string        `json:"region_id" validate:"required"`
	RegionName string        `json:"region_name" validate:"required"`
	Endpoint   string        `json:"endpoint"`
}

// Validate region create request.
func (req *RegionCreateReq) Validate() error {
	if len(req.Regions) == 0 {
		return errors.New("regions is required")
	}

	if len(req.Regions) > constant.BatchOperationMaxLimit {
		return fmt.Errorf("regions count should <= %d", constant.BatchOperationMaxLimit)
	}

	return nil
}

// -------------------------- Update --------------------------

// RegionBatchUpdateReq define region batch update request.
type RegionBatchUpdateReq struct {
	Regions []RegionBatchUpdate `json:"regions" validate:"required"`
}

// RegionBatchUpdate region batch update option.
type RegionBatchUpdate struct {
	ID          string        `json:"id"`
	Vendor      enumor.Vendor `json:"vendor"`
	RegionID    string        `json:"region_id"`
	RegionName  string        `json:"region_name"`
	IsAvailable int64         `json:"is_available"`
	Endpoint    string        `json:"endpoint"`
	Creator     string        `json:"creator"`
	Reviser     string        `json:"reviser"`
}

// Validate region batch update request.
func (req *RegionBatchUpdateReq) Validate() error {
	if len(req.Regions) == 0 {
		return errors.New("regions is required")
	}

	if len(req.Regions) > constant.BatchOperationMaxLimit {
		return fmt.Errorf("regions count should <= %d", constant.BatchOperationMaxLimit)
	}

	return nil
}

// -------------------------- List --------------------------

// RegionListReq region list req.
type RegionListReq struct {
	Field  []string           `json:"field" validate:"omitempty"`
	Filter *filter.Expression `json:"filter" validate:"required"`
	Page   *core.BasePage     `json:"page" validate:"required"`
}

// Validate region list request.
func (req *RegionListReq) Validate() error {
	return validator.Validate.Struct(req)
}

// RegionListResp define region list resp.
type RegionListResp struct {
	rest.BaseResp `json:",inline"`
	Data          *RegionListResult `json:"data"`
}

// RegionListResult define region list result.
type RegionListResult struct {
	Count   uint64         `json:"count,omitempty"`
	Details []RegionDetail `json:"details,omitempty"`
}

// -------------------------- Get --------------------------

// RegionGetResp defines get region response.
type RegionGetResp struct {
	rest.BaseResp `json:",inline"`
	Data          *RegionDetail `json:"data"`
}

// RegionDetail define region detail.
type RegionDetail struct {
	ID          string        `json:"id"`
	Vendor      enumor.Vendor `json:"vendor"`
	RegionID    string        `json:"region_id"`
	RegionName  string        `json:"region_name"`
	IsAvailable int64         `json:"is_available"`
	Endpoint    string        `json:"endpoint"`
	Creator     string        `json:"creator,omitempty"`
	Reviser     string        `json:"reviser,omitempty"`
	CreatedAt   *time.Time    `json:"created_at,omitempty"`
	UpdatedAt   *time.Time    `json:"updated_at,omitempty"`
}
