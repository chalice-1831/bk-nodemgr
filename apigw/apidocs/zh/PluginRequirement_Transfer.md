### 描述

- 该接口提供版本：v3.0.1-alpha.119+。
- 该接口所需权限：plugin_operate（操作插件）；按目标主机节点角色另需 agent_operate（操作 Agent）或 proxy_operate（操作 Proxy）。
- 该接口功能描述：使用记录中的下发文件名，将依赖文件传输到指定主机，返回已发起的传输任务信息。

### URL

POST /api/v3/plugin/requirement/transfer

### 输入参数

| 参数名称            | 参数类型 | 必选 | 描述                           |
| ------------------- | -------- | ---- | ------------------------------ |
| requirement_file_id | string   | 是   | 依赖文件记录的 UUID            |
| target_dir          | string   | 是   | 目标主机上的文件目录，不得为空 |
| target_host_id      | int64    | 是   | 目标主机 ID，不能为负数        |

**参数说明**：

- 接口对记录归属插件检查 `plugin_operate` 权限。
- 按目标主机节点角色检查 operate 权限：Agent 或空白角色需 `agent_operate`，Proxy 角色需 `proxy_operate`，未知角色需同时具备两者。
- `requirement_file_id` 必须是合法 UUID；`target_dir` 不能为空；`target_host_id` 不能为负数。
- 使用记录中保存的下发文件名发起单主机传输；响应中的 `task_id` 表示传输已发起，传输结果需另行查询。

### 调用示例

```json
{
  "requirement_file_id": "0b6f0f18-6f1f-4c88-9d3f-2a4f1f3d1a11",
  "target_dir": "/etc/bkmonitor",
  "target_host_id": 10001
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "request_id": "req-20260324-000016",
  "error": null,
  "data": {
    "task_id": "task-20260324-0001",
    "target_file_name": "agent.conf",
    "file_md5": "6f1ed002ab5595859014ebf0951522d9",
    "file_size": 18
  }
}
```

### 响应参数说明

| 参数名称   | 参数类型 | 描述                   |
| ---------- | -------- | ---------------------- |
| code       | int32    | 状态码，0表示成功      |
| message    | string   | 请求信息               |
| request_id | string   | 请求ID                 |
| error      | object   | 错误信息，成功时为null |
| data       | object   | 响应数据               |

#### data

| 参数名称         | 参数类型 | 描述                           |
| ---------------- | -------- | ------------------------------ |
| task_id          | string   | 传输任务 ID，用于查询传输状态  |
| target_file_name | string   | 下发到目标主机的文件名         |
| file_md5         | string   | 本次传输文件的 MD5             |
| file_size        | int64    | 本次传输文件的大小，单位为字节 |

#### error

| 参数名称 | 参数类型 | 描述         |
| -------- | -------- | ------------ |
| system   | string   | 错误系统标识 |
| message  | string   | 错误消息     |
| details  | array    | 错误详情列表 |

#### error.details[n]

| 参数名称 | 参数类型 | 描述     |
| -------- | -------- | -------- |
| code     | string   | 错误代码 |
| message  | string   | 错误消息 |
