### Description

- API Version: v3.0.1-alpha.119+.
- Required Permission: plugin_operate (Operate Plugin); depending on the target host node role, also agent_operate (Operate Agent) or proxy_operate (Operate Proxy).
- Function: Transfer the dependency file to the specified host using the record's delivery filename and return the launched transfer task information.

### URL

POST /api/v3/plugin/requirement/transfer

### Input Parameters

| Parameter Name      | Parameter Type | Required | Description                                     |
| ------------------- | -------------- | -------- | ----------------------------------------------- |
| requirement_file_id | string         | Yes      | UUID of the dependency file record              |
| target_dir          | string         | Yes      | Target directory on the host, must not be empty |
| target_host_id      | int64          | Yes      | Target host ID, must not be negative            |

**Parameter Notes**:

- The interface checks the `plugin_operate` permission for the plugin that owns the record.
- It checks operate permissions by the target host node role: Agent or blank role requires `agent_operate`; Proxy role requires `proxy_operate`; an unknown role requires both.
- `requirement_file_id` must be a valid UUID; `target_dir` must not be empty; `target_host_id` must not be negative.
- The transfer starts a single-host delivery using the record's delivery filename; `task_id` in the response means the transfer has been launched, and the result must be queried separately.

### Request Example

```json
{
  "requirement_file_id": "0b6f0f18-6f1f-4c88-9d3f-2a4f1f3d1a11",
  "target_dir": "/etc/bkmonitor",
  "target_host_id": 10001
}
```

### Response Example

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

### Response Parameters

| Parameter Name | Parameter Type | Description                        |
| -------------- | -------------- | ---------------------------------- |
| code           | int32          | Status code, 0 indicates success   |
| message        | string         | Request message                    |
| request_id     | string         | Request ID                         |
| error          | object         | Error information, null on success |
| data           | object         | Response data                      |

#### data

| Parameter Name   | Parameter Type | Description                                |
| ---------------- | -------------- | ------------------------------------------ |
| task_id          | string         | Transfer task ID, used to query the status |
| target_file_name | string         | Delivery filename on the target host       |
| file_md5         | string         | MD5 of the transferred file                |
| file_size        | int64          | Size of the transferred file in bytes      |

#### error

| Parameter Name | Parameter Type | Description             |
| -------------- | -------------- | ----------------------- |
| system         | string         | Error system identifier |
| message        | string         | Error message           |
| details        | array          | Error detail list       |

#### error.details[n]

| Parameter Name | Parameter Type | Description   |
| -------------- | -------------- | ------------- |
| code           | string         | Error code    |
| message        | string         | Error message |
