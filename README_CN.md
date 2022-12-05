# Golang SDK for Chaos Spec

## 支持的场景
* 延迟
* 异常
* 修改参数
* 修改变量
* 修改返回值
* 内存溢出
* Panic

## 如何使用
具体可参考 [business.go](./examples/business.go)
### 下载
```shell script
go get github.com/chaosblade-io/chaosblade-exec-golang
```

### 初始化
```shell script
chaos.Init()
```

### 在函数中定义故障演练目标
```shell script
c:=chaos.New()
```

### 添加匹配器
```shell script
c.AddMatchers(
    matcher.NewEqualMatcher("userId", request.Params["userId"]),
    matcher.NewEqualMatcher("name", request.Params["name"]))
```

### 设置故障行为
```shell script
if chaosResponse := c.Set(context.TODO(), action.NewDelayAction()); chaosResponse.Success {
	log.Println("trigger delay success")
}
```

## 支持的 API
```text
// 业务请求
// Post: http://localhost:8000/execute
// {"Params":{"userId":"1.3.0.1", "md5":"b493a3d94fbe67d74f39b8dd2c742024", "name":"AHAS_AGENT", "ts":"1535426008477","offset":"1s"}}

// 故障注入
// Post: http://localhost:9526/chaos/inject
// 例如panic: {"target": "main.(*Business).Execute","action":"panic","Flags":{"userId":"1.3.0.1", "message":"mock panic message"}}

// 故障恢复，请求同对应的故障注入参数
// Post: http://localhost:9526/chaos/recover
// 例如恢复panic: {"target": "main.(*Business).Execute","action":"panic","Flags":{"userId":"1.3.0.1", "message":"mock panic message"}}

// 故障命中数查询，请求同对应的故障注入参数
// Post: http://localhost:9526/chaos/metric
// 例如panic: {"target": "main.(*Business).Execute","action":"panic","Flags":{"userId":"1.3.0.1", "message":"mock panic message"}}

// 故障场景熔断
// Post: http://localhost:9526/chaos/circuit

// 取消场景熔断
// Post: http://localhost:9526/chaos/cancelcircuit

// 开启故障场景采集
// Post: http://localhost:9526/chaos/scenario/enable

// 故障场景查询，列出支持的场景，需要业务触发时且打开场景采集开关
// Get: http://localhost:9526/chaos/scenario/list

// 关闭故障场景采集
// Post: http://localhost:9526/chaos/scenario/disable

// 调整日志级别, 支持 error|warn|info|debug
// Get: http://localhost:9526/chaos/log?level=debug
```
