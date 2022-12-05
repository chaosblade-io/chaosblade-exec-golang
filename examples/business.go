package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/action"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/matcher"
)

type BusinessResult struct {
	Code    string
	Success bool
	Error   string
}

func (r *BusinessResult) Print() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "encode failed"
	}
	return string(bytes)
}

type Request struct {
	Id      string
	Headers map[string]string
	Params  map[string]string
}

func main() {
	chaos.Init()
	go func() {
		log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
	}()
	log.Printf("start")
	http.HandleFunc("/execute", func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Fprintf(writer, err.Error())
			return
		}
		r := &Request{}
		err = json.Unmarshal(body, r)
		if err != nil {
			fmt.Fprintf(writer, err.Error())
			return
		}
		business := &Business{}
		response := business.Execute(*r)
		fmt.Fprintf(writer, response.Print())
	})
	hold := make(chan struct{})
	<-hold
}

type Business struct {
}

// 业务请求
// Post: http://localhost:8000/execute
// {"Params":{"userId":"1.3.0.1", "md5":"b493a3d94fbe67d74f39b8dd2c742024", "name":"AHAS_AGENT", "ts":"1535426008477","offset":"1s"}}

// 故障注入，参数见故障点
// Post: http://localhost:9526/chaos/inject

// 故障恢复，同对应的故障注入参数
// Post: http://localhost:9526/chaos/recover

// 故障命中数查询，同对应的故障注入参数
// Post: http://localhost:9526/chaos/metric

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

func (b *Business) Execute(request Request) *BusinessResult {
	// 定义故障注入目标，取值为包名+函数名，例如此处为 main.(*Business).Execute
	c := chaos.New()
	// 添加故障注入目标
	c.AddMatchers(
		matcher.NewEqualMatcher("userId", request.Params["userId"]),
		matcher.NewEqualMatcher("name", request.Params["name"]))
	// TODO 设置故障点：修改变量的值
	// {"target": "main.(*Business).Execute","action":"modify","Flags":{"userId":"1.3.0.1","value":"Hanmeimei","effect-count":"5"}}
	var name = "bob"
	if chaosResponse := c.Set(context.TODO(), action.NewModifyAction(name)); chaosResponse.Success {
		// 修改后的值会保存在 result.result 中
		name = chaosResponse.Result.(string)
		log.Println("bob -> " + name)
	}
	// TODO 设置故障点：注入延迟
	// {"target": "main.(*Business).Execute","action":"delay","Flags":{"userId":"1.3.0.1","time":"3s","effect-count":"5"}}
	if chaosResponse := c.Set(context.TODO(), action.NewDelayAction()); chaosResponse.Success {
		log.Println("trigger delay success")
	}
	// TODO 设置故障点：修改参数的值
	// {"target": "main.(*Business).Execute","action":"modify","Flags":{"userId":"1.3.0.1","value":"{\"id\":\"mockid\"}","effect-count":"5"}}
	if chaosResponse := c.SetWithIndex(context.TODO(), 1, action.NewModifyAction(request)); chaosResponse.Success {
		log.Printf("%+v -> %+v", request, chaosResponse.Result.(Request))
	}
	// TODO 设置故障点：抛异常
	// 同一个场景支持多个故障点设置，可以指定故障点位置，使用 SetWithIndex 函数，不同场景故障点位置可以相同
	// {"target": "main.(*Business).Execute","action":"error","Flags":{"userId":"1.3.0.1","message":"mock error for chaos","index":"1"}}
	if chaosResponse := c.SetWithIndex(context.TODO(), 1, action.NewErrorAction()); chaosResponse.Success {
		log.Printf("expected error at 1 index: %v", chaosResponse.Result.(error))
	}
	// {"target": "main.(*Business).Execute","action":"error","Flags":{"userId":"1.3.0.1","message":"mock error for chaos","index":"2"}}
	if chaosResponse := c.SetWithIndex(context.TODO(), 2, action.NewErrorAction()); chaosResponse.Success {
		log.Printf("expected error at 2 index: %v", chaosResponse.Result.(error))
	}
	// do business....
	result := &BusinessResult{
		Code:    "200",
		Success: true,
		Error:   "",
	}
	// TODO 设置故障点：修改返回对象的值
	// {"target": "main.(*Business).Execute","action":"modify","Flags":{"userId":"1.3.0.1","value":"{\"code\":\"500\",\"success\":false,\"error\":\"mock response error\"}","index":"1"}}
	if chaosResponse := c.SetWithIndex(context.TODO(), 1, action.NewModifyAction(result)); chaosResponse.Success {
		// 指针类型，不用重新赋值也会修改原有值
		log.Printf("new result value: %+v", result)
	}

	// TODO 设置故障点：内存溢出
	// {"target": "main.(*Business).Execute","action":"oom","Flags":{"userId":"1.3.0.1"}}
	if chaosResponse := c.Set(context.TODO(), action.NewOomAction()); chaosResponse.Success {
		log.Println("trigger oom success")
	}

	// TODO 设置故障点：panic
	// {"target": "main.(*Business).Execute","action":"panic","Flags":{"userId":"1.3.0.1", "message":"mock panic message"}}
	c.Set(context.TODO(), action.NewPanicAction())
	return result
}
