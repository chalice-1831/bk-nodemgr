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

package v3_test

import (
	"encoding/json"
	"testing"

	protoBackend "github.com/TencentBlueKing/bk-nodemgr/pkg/proto/backend/api/v3"
	"github.com/stretchr/testify/require"
)

func TestPluginRequirementReplacementJSON(t *testing.T) {
	t.Parallel()

	const fileID = "0b6f0f18-6f1f-4c88-9d3f-2a4f1f3d1a11"
	const requirement = `{"requirement_file_id":"` + fileID + `","target_dir":"conf","target_file_name":"rules.json"}`
	for _, tc := range []struct {
		name   string
		fields string
		err    string
	}{
		{name: "omitted", err: "explicit array"},
		{name: "null", fields: `,"requirements":null`, err: "explicit array"},
		{name: "clear", fields: `,"requirements":[]`},
		{name: "replace", fields: `,"requirements":[` + requirement + `]`},
		{name: "traversal", fields: `,"requirements":[{"requirement_file_id":"` + fileID +
			`","target_dir":"conf/../extra","target_file_name":"rules.json"}]`, err: "target_dir"},
		{name: "duplicate destination", fields: `,"requirements":[` + requirement + `,{"requirement_file_id":"` + fileID +
			`","target_dir":"Conf","target_file_name":"RULES.JSON"}]`, err: "duplicate"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var req protoBackend.PluginSetRequirementsReq
			require.NoError(t, json.Unmarshal([]byte(`{"plugin_name":"example"`+tc.fields+`}`), &req))
			req.AutoConvert()
			if tc.err != "" {
				require.ErrorContains(t, req.Validate(), tc.err)
				return
			}
			require.NoError(t, req.Validate())
			require.NotNil(t, req.ConvertRequirementsToTypes())
		})
	}
}

func TestPluginRequirementResponseEmptyArray(t *testing.T) {
	t.Parallel()

	var resp protoBackend.PluginGetRequirementsResp
	resp.ConvertRequirementsFromTypes(nil)
	body, err := json.Marshal(resp.GetData())
	require.NoError(t, err)
	require.JSONEq(t, `{"requirements":[]}`, string(body))
}
