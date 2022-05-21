package jsonCheck

import (
	"fmt"
	"github.com/joker1222/tools"
	"sort"
	"reflect"
)

var jsonTypeToGoType = map[string]string{
	"Object":"map[string]interface {}",
	"String":"string",
	"Number":"float64",
	"Boolean":"bool",
	"Array":"[]interface {}",
}
var goTypeToJsonType = map[string]string{
	"map[string]interface {}":"Object",
	"string":"String",
	"float64":"Number",
	"bool":"Boolean",
	"[]interface {}":"Array",
}
/*
| JsonType|       GolangType      |
| ------- | --------------------- |
| Object  | map[string]interface{}|
| Number  | float64               |
| Boolean | bool                  |
| String  | string                |
| Array   | []interface{}         |
*/

func checkType(v interface{},k ,jsonType string) error {
	switch jsonType{
	case  "Object":
		if _,ok:=v.(map[string]interface{});!ok{
			return fmt.Errorf("<JsonKey:\"%v\"> type error  , your valueType is <%v> ,  it should be <%v>",k,goTypeToJsonType[reflect.TypeOf(v).String()],jsonType)
		}
	case  "Number":
		if _,ok:=v.(float64);!ok{
			return fmt.Errorf("<JsonKey:\"%v\"> type error  , your valueType is <%v> ,  it should be <%v>",k,goTypeToJsonType[reflect.TypeOf(v).String()],jsonType)
		}
	case "Boolean":
		if _,ok:=v.(bool);!ok{
			return fmt.Errorf("<JsonKey:\"%v\"> type error  , your valueType is <%v> ,  it should be <%v>",k,goTypeToJsonType[reflect.TypeOf(v).String()],jsonType)
		}
	case "Array":
		if _,ok:=v.([]interface{});!ok{
			return fmt.Errorf("<JsonKey:\"%v\"> type error  , your valueType is <%v> ,  it should be <%v>",k,goTypeToJsonType[reflect.TypeOf(v).String()],jsonType)
		}
	case "String":
		if _,ok:=v.(string);!ok{
			return fmt.Errorf("<JsonKey:\"%v\"> type error  , your valueType is <%v> ,  it should be <%v>",k,goTypeToJsonType[reflect.TypeOf(v).String()],jsonType)
		}
	default:
		return fmt.Errorf("rule jsonType invalid , your jsonType is <%v> , it should be <'Object','Number','Boolean','Array'>",jsonType)
	}
	fmt.Println(fmt.Sprintf("<JsonKey:\"%v\"> <_Type:%v>  checkType successful",k,jsonType))
	return nil
}

type Node struct{
	keys []string
	value interface{}
}

func parseRule(rule map[string]interface{}) []Node{
	q2:=tools.NewQueue()
	queue:=tools.NewQueue()
	queue.Push(Node{
		keys: make([]string,0),
		value:rule,
	})
	//广度优先遍历
	//因为规则配置文件中不会出现[](数组类型)，所以下面省略了数组类型的判定
	for queue.Len()!=0 {
		node:=queue.Pop().(Node)
		if object,ok:=node.value.(map[string]interface{});ok{
			for k,v:=range object {
				if k[0]!='_' || k == "_Element"{
					queue.Push(Node{
						keys: append(append(make([]string, 0),node.keys...),k),
						value:v,
					})
					q2.Push(Node{
						keys: append(append(make([]string, 0),node.keys...),k),
						value:v,
					})
				}
			}
		}
	}
	var nodeList =make([]Node,0)
	for q2.Len() != 0{
		node:=q2.Pop().(Node)
		nodeList=append(nodeList,node)
	}
	sort.Sort(NodeList(nodeList))
	return nodeList
}

type NodeList []Node
func (m NodeList) Len() int {
	return len(m)
}
func (m NodeList) Less(i, j int) bool {
	return len(m[i].keys) < len(m[j].keys)
}
func (m NodeList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
type CheckNode struct{
	keyStr string
	jsonType string
	value interface{}
}
func FoundStringArr(k string,l []string) bool{
	for _,v:=range l{
		if k == v{
			return true
		}
	}
	return false
}

func Check(ruleJson ,conf map[string]interface{}) []error {
	rule:=parseRule(ruleJson)
	noRequireds:=map[string]struct{}{}
	requiredErrList:=make([]string,0)
	errList:=make([]error,0)
	for _,node:=range rule{
		var confValue interface{}
		confValue=conf
		ruleKeys:=node.keys
		ruleValue:=node.value.(map[string]interface{})
		ruleType:=ruleValue["_Type"].(string)
		ruleRequired:=ruleValue["_Required"].(bool)
		var ruleRange []interface{}
		var ruleDefault interface{}
		if _Range,ok:=ruleValue["_Range"];ok{
			ruleRange=_Range.([]interface{})
		}
		if _Default,ok:=ruleValue["_Default"];ok{
			ruleDefault=_Default
		}
		skip:=false
		for _,rel:=range requiredErrList{
			if FoundStringArr(rel,ruleKeys){
				skip=true
				break
			}
		}
		if skip{
			continue
		}
		if !ruleRequired{
			noRequireds[ruleKeys[len(ruleKeys)-1]]= struct{}{}
		}
		checkList,err1,err2:=Recursion(ruleDefault,ruleRange,noRequireds,ruleType,"",ruleKeys,confValue)
		if err1!=nil || err2!=nil{
			if err1!=nil{
				errList=append(errList,err1)
				requiredErrList=append(requiredErrList,ruleKeys[len(ruleKeys)-1])
			} else {
				errList=append(errList,err2)
			}
		}else{
			for _,checkNode:=range checkList{
				if err:=checkType(checkNode.value,checkNode.keyStr,checkNode.jsonType);err!=nil{
					errList=append(errList,err)
				}
			}
		}
	}
	return errList
}

func Recursion(ruleDefault interface{},ruleRange []interface{},noRequired map[string]struct{},jsonType , kstr string ,keys []string,value interface{}) ([]CheckNode,error,error){
	for i,k:=range keys{
		if k == "_Element"{
			checkList:=make([]CheckNode,0)
			for j,index:=range value.([]interface{}){
				if i == len(keys)-1{
					checkList=append(checkList,CheckNode{
						keyStr:   fmt.Sprintf("%v[%v]",kstr[:len(kstr)-1],j),
						jsonType: jsonType,
						value:    index,
					})
				} else {
					cl,err1,err2:=Recursion(ruleDefault,ruleRange,noRequired,jsonType,fmt.Sprintf("%v[%v].",kstr[:len(kstr)-1],j),keys[i+1:],index)
					if err1 != nil || err2 != nil{
						return nil,err1,err2
					} else{
						checkList=append(checkList,cl...)
					}
				}
			}
			return checkList,nil,nil
		} else {
			kstr+=k+"."
			v,vok:=value.(map[string]interface{})[k]
			if _,ok:=noRequired[k];ok && !vok {
				//如果不是必填字段，并且用户没配，则不继续向下校验，直接返回空
				fmt.Printf("<JsonKey:%v> this key is no required, and user not config , skip check ... \n",kstr[:len(kstr)-1])
				if ruleDefault != nil{ //此处判断下规则中是否给选填参数提供了默认值，如果没有，后续会删除这个key，否则将默认值加上
					value.(map[string]interface{})[k]= ruleDefault
					fmt.Printf("<JsonKey:%v> this key is no required, and user not config , but have default <%v> ... \n",kstr[:len(kstr)-1],ruleDefault)
				}
				return nil,nil,nil
			}
			if _,ok:=noRequired[k];!ok && !vok{ //如果是必填字段但用户没有配置，则直接报错并返回空
				return nil,fmt.Errorf("<JsonKey:%v> this key is required, but user not config ... ",kstr[:len(kstr)-1]),nil
			}
			//如果是非必填字段，但用户却配了，仍然要向下进行配置检查
			value=v
		}
	}
	if ruleRange!=nil{
		if !FoundRuleRangeArr(value,ruleRange){
			return nil,nil,fmt.Errorf("<JsonKey:%v> this key is existed, and type correct,but its value is not in the specified range, which is <%v> , your value is <%v> ... ",kstr[:len(kstr)-1],ruleRange,value)
		}
	}
	return []CheckNode{
		{
			keyStr:   kstr[:len(kstr)-1],
			jsonType: jsonType,
			value:    value,
		},
	},nil,nil
}

func FoundRuleRangeArr(k interface{},ruleRange []interface{}) bool{
	for _,v:=range ruleRange{
		if k == v {
			return true
		}
	}
	return false
}
