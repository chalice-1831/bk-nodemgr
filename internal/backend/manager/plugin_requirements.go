/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 节点管理 (BlueKing - Node Management) available.
 * Copyright (C) Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.

 * We undertake not to change the open source license (MIT license) applicable

 * to the current version of the project delivered to anyone in the future.
 */

package manager

import (
	"errors"
	"fmt"

	"github.com/TencentBlueKing/bk-nodemgr/pkg/contextx"
	"github.com/TencentBlueKing/bk-nodemgr/pkg/types"
)

// ErrPluginRequirementReference indicates a file reference unavailable to the plugin.
var ErrPluginRequirementReference = errors.New("requirement file is unavailable to this plugin")

// SetPluginRequirements resolves file ownership before replacing the plugin's declarations.
func (mgr *Manager) SetPluginRequirements(nCtx contextx.IContext, pluginName string, requirements []types.PluginRequirement) error {
	if err := nCtx.CheckTenantID(); err != nil {
		return err
	}
	if len(requirements) > 0 {
		if err := mgr.validatePluginRequirementReferences(nCtx, pluginName, requirements); err != nil {
			return err
		}
	}

	return mgr.conf.StoragePlugin.SetPluginRequirements(nCtx, pluginName, requirements)
}

func (mgr *Manager) validatePluginRequirementReferences(nCtx contextx.IContext, pluginName string, requirements []types.PluginRequirement) error {
	pending := make(map[string]struct{}, len(requirements))
	fileIDs := make([]string, 0, len(requirements))
	for _, requirement := range requirements {
		id := requirement.RequirementFileID
		if _, exists := pending[id]; exists {
			continue
		}
		pending[id] = struct{}{}
		fileIDs = append(fileIDs, id)
	}

	condition := &types.PluginRequirementFileCondition{
		ExactInclude: &types.PluginRequirementFileExactFields{
			PluginName:        []string{pluginName},
			RequirementFileID: fileIDs,
		},
	}
	records, _, err := mgr.conf.FileHandler.ListPluginRequirements(nCtx, types.Page{Limit: len(fileIDs)}, condition)
	if err != nil {
		return fmt.Errorf("failed to resolve plugin requirement files: %w", err)
	}

	for _, record := range records {
		if record != nil && record.PluginName == pluginName {
			delete(pending, record.RequirementFileID)
		}
	}
	for _, id := range fileIDs {
		if _, missing := pending[id]; missing {
			return fmt.Errorf("%w: %s", ErrPluginRequirementReference, id)
		}
	}

	return nil
}
