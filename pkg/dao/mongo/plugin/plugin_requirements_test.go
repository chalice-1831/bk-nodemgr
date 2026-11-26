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
	"testing"

	"github.com/TencentBlueKing/bk-nodemgr/pkg/types"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestPluginRequirementUpsertPreservesDeclaration(t *testing.T) {
	t.Parallel()

	requirements := []types.PluginRequirement{{
		RequirementFileID: "0b6f0f18-6f1f-4c88-9d3f-2a4f1f3d1a11",
		TargetDir:         "conf",
		TargetFileName:    "rules.json",
	}}
	plugin := &types.Plugin{Name: "example", Group: "updated", Requirements: requirements}
	models := buildUpsertManyParams([]*Plugin{convPluginFromTypes(plugin)})
	require.Len(t, models, 1)

	model, ok := models[0].(*mongo.UpdateOneModel)
	require.True(t, ok)
	require.NotNil(t, model.Upsert)
	require.True(t, *model.Upsert)
	encoded, err := bson.Marshal(model.Update)
	require.NoError(t, err)
	fields := bson.Raw(encoded).Lookup("$set").Document()
	require.Equal(t, "updated", fields.Lookup(FieldKeyGroup).StringValue())
	for _, key := range []string{"data", FieldKeyRequirements} {
		_, err := fields.LookupErr(key)
		require.Error(t, err, "ordinary upsert must preserve existing dependency declarations")
	}
}
