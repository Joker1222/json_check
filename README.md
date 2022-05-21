# Overview

# Install
```bash
$ go get github.com/Joker1222/json_check
```

# Example
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/joker1222/json_check"
)

var ruleStr1 =
`
{
	"root":{
		"_Type":"Object",
		"_Required":true,
		"leaf":{
			"_Type":"Number",
			"_Required":true
        	}
	}
}
`
var ruleStr2 =
`
{
	"root":{
		"_Type":"Object",
		"_Required":true,
		"leaf":{
			"_Type":"String",
			"_Required":true
        	}
	}
}
`
var confStr =
`
{
	"root":{
		"leaf":1
	}
}
`
func main()  {
	conf:=map[string]interface{}{}
	rule1:=map[string]interface{}{}
	rule2:=map[string]interface{}{}
	_=json.Unmarshal([]byte(confStr),&conf)
	_=json.Unmarshal([]byte(ruleStr1),&rule1)
	_=json.Unmarshal([]byte(ruleStr2),&rule2)
	errList:=json_check.Check(rule1,conf)
	if len(errList)!=0{
		for _,v:=range errList{
			fmt.Println(v)
		}
	}
	fmt.Println("---------------------------")
	errList=json_check.Check(rule2,conf)
	if len(errList)!=0{
		for _,v:=range errList{
			fmt.Println(v)
		}
	}
}
```
```
$ go run
<JsonKey:"root"> <_Type:Object>  checkType successful
<JsonKey:"root.leaf"> <_Type:Number>  checkType successful
---------------------------
<JsonKey:"root"> <_Type:Object>  checkType successful
<JsonKey:"root.leaf"> type error  , your valueType is <Number> ,  it should be <String>
```
