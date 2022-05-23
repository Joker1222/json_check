# Overview

# 使用说明

1.前提
- 所有配置不允许以'_'开头命名（因为规则配置中需要判断带有'_'的规则参数）
- 对于数组Array类型，其中的元素类型必须一致，Object元素中的Key也需要一致
- 不允许出现自定义的Object类型

2.规则用法
|规则参数|说明|示例|
|-|-|-|
|_Type|字段类型(Object/Array/Number/String/Boolean)|Object|
|_Required|是否必填(true/false)|true|
|_Default|_Required为false时，可能需要默认值(如果是非必填参数且被检测配置没配，会自动加入该值)|
|_Element|数组参数的元素类型，是一个Object，配置方式与其他规则|{"_Type":"Number","_Required":true}|
|_Range|用于规定该参数的值在一个特定范围内|["HELLO","WORLD"]|

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
$ go run main.go
<JsonKey:"root"> <_Type:Object>  checkType successful
<JsonKey:"root.leaf"> <_Type:Number>  checkType successful
---------------------------
<JsonKey:"root"> <_Type:Object>  checkType successful
<JsonKey:"root.leaf"> type error  , your valueType is <Number> ,  it should be <String>
```
