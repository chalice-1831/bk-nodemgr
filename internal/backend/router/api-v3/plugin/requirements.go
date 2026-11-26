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

package plugin

import (
	"errors"

	"github.com/TencentBlueKing/bk-nodemgr/internal/backend/auth"
	"github.com/TencentBlueKing/bk-nodemgr/internal/backend/manager"
	"github.com/TencentBlueKing/bk-nodemgr/pkg/dao/mongo/base"
	"github.com/TencentBlueKing/bk-nodemgr/pkg/logger"
	protoBackend "github.com/TencentBlueKing/bk-nodemgr/pkg/proto/backend/api/v3"
	resterrf "github.com/TencentBlueKing/bk-nodemgr/pkg/rest/errf"
	restserver "github.com/TencentBlueKing/bk-nodemgr/pkg/rest/server"
)

// GetRequirements returns the declarations of the requested plugin.
func (h *handler) GetRequirements(rCtx restserver.IContext) (interface{}, error) {
	req := new(protoBackend.PluginGetRequirementsReq)
	if err := rCtx.BindJSON(req); err != nil {
		logger.G.Biz(rCtx).WithErr(err).Error("failed to bind plugin requirements query")
		return nil, resterrf.ErrWrap(resterrf.InvalidParameter, err)
	}

	pluginName := req.GetPluginName()
	log := logger.G.Biz(rCtx).With("plugin-name", pluginName)
	if err := h.authorizeRequirementPlugin(rCtx, auth.ActionPluginView, pluginName); err != nil {
		log.WithErr(err).Error("failed to authorize plugin requirements query")
		return nil, resterrf.ErrWrap(resterrf.PermissionDenied, err)
	}

	plugin, err := h.daoPlugin.GetPlugin(rCtx, pluginName)
	if err != nil {
		log.WithErr(err).Error("failed to get plugin requirements")
		if errors.Is(err, base.ErrRecordNoFound()) {
			return nil, resterrf.ErrWrap(resterrf.RecordNotFound, err)
		}

		return nil, resterrf.ErrWrap(resterrf.DBExecCmdFailed, err)
	}

	resp := new(protoBackend.PluginGetRequirementsResp)
	resp.ConvertRequirementsFromTypes(plugin.Requirements)

	return resp.GetData(), nil
}

// SetRequirements replaces all declarations of the requested plugin.
func (h *handler) SetRequirements(rCtx restserver.IContext) (interface{}, error) {
	req := new(protoBackend.PluginSetRequirementsReq)
	if err := rCtx.BindJSON(req); err != nil {
		logger.G.Biz(rCtx).WithErr(err).Error("failed to bind plugin requirements replacement")
		return nil, resterrf.ErrWrap(resterrf.InvalidParameter, err)
	}

	pluginName := req.GetPluginName()
	log := logger.G.Biz(rCtx).With("plugin-name", pluginName)
	if err := h.authorizeRequirementPlugin(rCtx, auth.ActionPluginOperate, pluginName); err != nil {
		log.WithErr(err).Error("failed to authorize plugin requirements replacement")
		return nil, resterrf.ErrWrap(resterrf.PermissionDenied, err)
	}

	if err := h.pluginMgrIface.SetPluginRequirements(rCtx, pluginName, req.ConvertRequirementsToTypes()); err != nil {
		log.WithErr(err).Error("failed to set plugin requirements")
		switch {
		case errors.Is(err, manager.ErrPluginRequirementReference):
			return nil, resterrf.ErrWrap(resterrf.InvalidParameter, err)
		case errors.Is(err, base.ErrRecordNoFound()):
			return nil, resterrf.ErrWrap(resterrf.RecordNotFound, err)
		default:
			return nil, resterrf.ErrWrap(resterrf.BackendOperateFailed, err)
		}
	}

	return nil, nil //nolint:nilnil // A successful replacement has no response data.
}
