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

package types

import "hcm/pkg/kit"

// AccountInterface defines all the account related operations in the hybrid cloud
type AccountInterface interface {
	AccountCheck(kt *kit.Kit, secret *Secret) error
}

// Secret defines the hybrid cloud's secret info.
// TODO replace with actual account secret info
type Secret struct {
	// ID is the secret id to do credential
	ID string `json:"id,omitempty"`
	// Key is the secret key to do credential
	Key string `json:"key,omitempty"`

	// Json carry a json formatted credential information for
	// GCP(Google Cloud Platform) vendor only.
	Json []byte `json:"json,omitempty"`

	// TenantID is used only for azure credential
	TenantID string `json:"tenant_id,omitempty"`
	// SubscriptionID is used only for azure credential
	SubscriptionID string `json:"subscription_id,omitempty"`

	// ProjectID is cloud vendor project id.
	ProjectID string `json:"project_id,omitempty"`
}
