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

// Package plugin this package is used to store the need data for plugin.
package plugin

import (
	"time"

	"github.com/TencentBlueKing/bk-nodemgr/pkg/contextx"
	"github.com/TencentBlueKing/bk-nodemgr/pkg/dao/mongo/base"
	"github.com/TencentBlueKing/bk-nodemgr/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func newDao(client *mongo.Database, tableName string) *dao {
	d := &dao{
		client:    client.Collection(tableName),
		tableName: tableName,
		IOrm:      nil,
	}
	d.IOrm = base.NewOrm[*Plugin, Plugin](d)

	return d
}

type dao struct {
	client    *mongo.Collection
	tableName string

	base.IOrm[*Plugin, Plugin]
}

// GetClient get the dao's client.
func (d *dao) GetClient() *mongo.Collection {
	return d.client
}

// GetTableName get the dao's table name.
func (d *dao) GetTableName() string {
	return d.tableName
}

// GetIndexes get the dao's indexes.
func (d *dao) GetIndexes() []mongo.IndexModel {
	var indexes []mongo.IndexModel

	return indexes
}

func (d *dao) upsertMany(nCtx contextx.IContext, plugins ...*Plugin) error {
	models := buildUpsertManyParams(plugins)

	result, err := d.client.BulkWrite(nCtx, models)
	if err != nil {
		return err
	}

	if result.UpsertedCount > 0 {
		logger.G.Sys().With("inserted-count", result.UpsertedCount).Info("upserted config policy template")
	}

	if result.MatchedCount > 0 {
		logger.G.Sys().With("matched-count", result.MatchedCount).Info("upserted config policy template")
	}

	return nil
}

func buildUpsertManyParams(plugins []*Plugin) []mongo.WriteModel {
	models := make([]mongo.WriteModel, 0, len(plugins))
	now := time.Now()
	for _, plugin := range plugins {
		filter := bson.D{
			bson.E{Key: FieldKeyName, Value: plugin.Name},
		}

		// Requirements are managed by the dedicated declaration update.
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: base.FieldKeyIsDeleted, Value: false},
				{Key: base.FieldKeyUpdatedAt, Value: now},
				{Key: "data.tenant_id", Value: plugin.TenantID},
				{Key: FieldKeyName, Value: plugin.Name},
				{Key: FieldKeyGroup, Value: plugin.Group},
				{Key: FieldKeyPkgName, Value: plugin.PkgName},
				{Key: FieldKeyMemo, Value: plugin.Memo},
				{Key: FieldKeyVisibleBizIDs, Value: plugin.VisibleBizIDs},
			}},
			{Key: "$setOnInsert", Value: bson.D{{Key: base.FieldKeyCreatedAt, Value: now}}},
		}

		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}

	return models
}

func (d *dao) updateRequirements(nCtx contextx.IContext, pluginName string, requirements []pluginRequirement) error {
	filter := WithName(pluginName)(base.AliveFilter())
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: FieldKeyRequirements, Value: requirements},
		{Key: base.FieldKeyUpdatedAt, Value: time.Now()},
	}}}
	// Check the match so replacing declarations for a missing plugin fails.
	result, err := d.client.UpdateOne(nCtx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return base.ErrRecordNoFound()
	}

	return nil
}
