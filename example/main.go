package main

import (
	"encoding/json"
	"fmt"
	"github.com/joker1222/json_check"
	"os"
)

func main() {
	var conf, rule map[string]interface{}
	confBytes, _ := os.ReadFile(os.Args[1])
	ruleBytes, _ := os.ReadFile(os.Args[2])
	_ = json.Unmarshal(confBytes, &conf)
	_ = json.Unmarshal(ruleBytes, &rule)
	errList := json_check.Check(rule, conf)
	if len(errList) != 0 {
		for _, v := range errList {
			fmt.Println(v)
		}
	}
}
