# Golang SDK for Chaos Spec

## Supported scenarios
* delay
* error
* modify parameters
* modify variables
* modify return value
* oom
* panic

## How to use
Refer to [business.go](./examples/business.go)

### Go get
```shell script
go get github.com/chaosblade-io/chaosblade-exec-golang
```

### Initialization
Use the following command to initialize, only need to initialize once.
```shell script
chaos.Init()
```

### Define the failure target in the function
```shell script
c:=chaos.New()
```

### Add experiment matchers
```shell script
c.AddMatchers(
    matcher.NewEqualMatcher("userId", request.Params["userId"]),
    matcher.NewEqualMatcher("name", request.Params["name"]))
```

### Set the fault behavior which will not take effect at this time until the experiment rules issued.
```shell script
if chaosResponse := c.Set(context.TODO(), action.NewDelayAction()); chaosResponse.Success {
	log.Println("trigger delay success")
}
```

## Supported APIs
```text
// Demo request url
// Post: http://localhost:8000/execute
// {"Params":{"userId":"1.3.0.1", "md5":"b493a3d94fbe67d74f39b8dd2c742024", "name":"AHAS_AGENT", "ts":"1535426008477","offset":"1s"}}

// Fault injection
// Post: http://localhost:9526/chaos/inject
// Such as panic: {"target": "main.(*Business).Execute","action":"panic","Flags":{"userId":"1.3.0.1", "message":"mock panic message"}}

// Fault recovery. The request is the same as the corresponding fault injection parameter.
// Post: http://localhost:9526/chaos/recover
// Such as panic: {"target": "main.(*Business).Execute","action":"panic","Flags":{"userId":"1.3.0.1", "message":"mock panic message"}}

// Query the number of fault hits. The request is the same as the corresponding fault injection parameter.
// Post: http://localhost:9526/chaos/metric
// Such as panic: {"target": "main.(*Business).Execute","action":"panic","Flags":{"userId":"1.3.0.1", "message":"mock panic message"}}

// The failure scenario circuit and will stop all experiments.
// Post: http://localhost:9526/chaos/circuit

// Cancel circuit and will recover all experiments.
// Post: http://localhost:9526/chaos/cancelcircuit

// Enable failure scenarios collection, for scenarios list.
// Post: http://localhost:9526/chaos/scenario/enable

// Query fault scenarios, and enable the scenarios collection first.
// Get: http://localhost:9526/chaos/scenario/list

// Disable failure scenarios collection
// Post: http://localhost:9526/chaos/scenario/disable

// Modify log level, support error|warn|info|debug
// Get: http://localhost:9526/chaos/log?level=debug
```
