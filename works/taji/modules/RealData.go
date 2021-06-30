package modules

import (
	"sync"
)

/**
 * @Author: Tao Jun
 * @Description: modules
 * @File:  RealData
 * @Version: 1.0.0
 * @Date: 2021/6/29 下午5:00
 */

type RealDataSub struct {
	Data       int16  `json:"data"`
	Gaojing    int16  `json:"gaojing"`
	Yujing     int16  `json:"yujing"`
	Danwei     string `json:"danwei"`
	Fengji     string `json:"fengji"`
	Type       string `json:"type"`
	Isalarm    string `json:"isalarm"`
	Action     string `json:"action"`
	Actiondata string `json:"actiondata"`
}

type RealDataEntity struct {
	Userid     string        `bson:"userid" json:"userid"`
	Username   string        `bson:"username" json:"username"`
	TowerId    string        `bson:"towerId" json:"towerId"`
	MID        string        `bson:"MID" json:"MID"`
	ProjectID  string        `bson:"projectID" json:"projectID"`
	PanbanId   string        `bson:"paibanId" json:"paibanId"`
	InsertTime int64         `bson:"insertTime"`
	Data       []interface{} `bson:"data" json:"data"`
}

var RealDataEntityFree = sync.Pool{
	New: func() interface{} {
		return &RealDataEntity{}
	},
}
