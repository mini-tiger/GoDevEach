package cron

import (
	"godev/works/train/g"
	"time"
	"sync"
	"fmt"
	"tjtools/net/udp"
	"tjtools/conversion"
	"encoding/json"
)

//const imageExt = "_image"  //redis 中存入IMAGES SET时key的名字  vehicle_serial+imageExt
const vehicleInterval = 10 // 扫描vehicle 数据库间隔

type VehicleResult struct {
	sync.RWMutex
	VehicleJson *VehicleJson // 每个要 车厢的json
}

type VehicleJson struct {
	TRAIN_ID       string   `json:"TRAIN_ID"`
	PASS_TIME_CHAR string   `json:"PASS_TIME"`
	INDEX_ID       string   `json:"FOLDER"`
	VehicleOrder   string   `json:"VEHICLE_ORDER"`
	VehicleType    string   `json:"VEHICLE_TYPE"`
	POS_AB_FALG    string   `json:"AB"`
	Images         []string `json:"Images"`//todo omitempty, 字段是默认值时，生成json时不会有此项

	SaveUnixStamp  int      `json:"Save_unix_stamp,omitempty"` // todo 发送时要 清空此项, 存入redis对比 时间是否超过最大    "maxvehicletime":120
	Send           int     `json:"Send,omitempty"`            // todo 是否发送过 ，存入redis对比 时间是否超过最大,  1 没有，0 有
	Vehicle_serial string   `json:"Vehicle_serial,omitempty"`            // todo 程序需要使用，发送前清空
}

type TrainInfo struct {
	TRAIN_ID       string
	PASS_TIME_CHAR string
	INDEX_ID       string
	TRAIN_SERIAL   string
}

func GetPicMission(result map[string]string) {
	g.Logger().Debug("开始获取train_serial:%s ", result["TRAIN_SERIAL"])
	var trainInfo TrainInfo
	conversion.MapToStructStr(result, &trainInfo) // 将map 转换为 结构体
	go trainInfo.BeginMVehicle()
}

func (t *TrainInfo) BeginMVehicle() {
	now := g.GetNow()
	for { //  todo 每次都获取 train_serial下面的所有 veh ，每次 判断存入redis的vehJSON是否 到时，到时则获取sys1发送
		t.CheckMVehicle()
		if int(time.Now().Unix())-now >= g.Config().Train.TrainInterval+30 { //todo 当前时间 － 循环开始时间  超过两个列车的间隔延长60秒
			// todo 把没发送的，所有都获取 sys 发出去,还没有写 这里的代码  send=1 代表没有发送过
			break
		}
		time.Sleep(time.Duration(vehicleInterval) * time.Second)
	}
}

func (t *TrainInfo) CheckMVehicle() {

	sql := fmt.Sprintf("SELECT TRAIN_SERIAL,VEHICLE_SERIAL,VEHICLE_ORDER,VEHICLE_TYPE,POS_AB_FLAG FROM TF_OP_VEHICLE "+
		"WHERE train_serial  ='%s' and vehicle_order is not null  order by vehicle_order ", t.TRAIN_SERIAL)
	if mveh, b := g.GetSqlData(sql); b { // todo 扫描 train_serial 下的所有vehicle
		now := g.GetNow()
		for _, v := range mveh {
			if g.RedisMission.StringExists(v["VEHICLE_SERIAL"]) { // todo 存入过判断是否超时，是否发送过
			//fmt.Println("exist ",v["VEHICLE_SERIAL"])
				go t.SendJson(v["VEHICLE_SERIAL"])
			} else {
				g.Logger().Debug("不存在 加入redis,train_serial:%s ,vehicle_serial:%s", t.TRAIN_SERIAL, v["VEHICLE_SERIAL"])
				var tmpVehJson VehicleResult
				tmpVehJson.SetRedis(t, v, now)
			}
		}
	}
}

func (v *VehicleResult) SetRedis(trainInfo *TrainInfo, veh map[string]string, now int) {
	v.Lock()
	defer v.Unlock()

	v.VehicleJson = &VehicleJson{TRAIN_ID: trainInfo.TRAIN_ID,
		PASS_TIME_CHAR: trainInfo.PASS_TIME_CHAR,
		INDEX_ID: trainInfo.INDEX_ID,
		VehicleOrder: veh["VEHICLE_ORDER"],
		VehicleType: veh["VEHICLE_TYPE"],
		POS_AB_FALG: veh["POS_AB_FLAG"],
		Vehicle_serial: veh["VEHICLE_SERIAL"],
		SaveUnixStamp: now, // 写入当前时间，后面判断是超时
		Send: 1} //没有发送过

	err := g.RedisMission.SetJson(veh["VEHICLE_SERIAL"], v.VehicleJson, time.Duration(1800)*time.Second) // 超时时间30分钟,自动删除redis中的数据\
	if err != nil {
		g.Logger().Debug("存入redis 失败 json:%+v", v.VehicleJson)
	}
	g.Logger().Debug("存入redis 成功 json:%+v", v.VehicleJson)
}

func (t *TrainInfo) SendJson(vehicleKey string) {
	var vehjsonStruct VehicleJson
	if err := g.RedisMission.GetJson(vehicleKey, &vehjsonStruct); err != nil {
		g.Logger().Error("获取json key:%s,err:%s", vehicleKey, err)
		return
	}
	//fmt.Println(vehjsonStruct.Send==0,(g.GetNow()-vehjsonStruct.SaveUnixStamp) < g.Config().Train.MaxVehicleTime)
	if vehjsonStruct.Send ==0 || (g.GetNow()-vehjsonStruct.SaveUnixStamp) < g.Config().Train.MaxVehicleTime {
		return // todo 发送过或者没有到发送的时间
	} else {
		g.Logger().Printf("准备发送vehicleSerial:%s\n", vehicleKey)
		vehImages := GetImages(vehicleKey, &vehjsonStruct)
		vehjsonStruct.Images = vehImages
		vehjsonStruct.Send=0
		vehjsonStruct.SaveUnixStamp=0
		vehjsonStruct.Vehicle_serial=""
		byteVeh, err := json.Marshal(&vehjsonStruct)
		if err != nil {
			g.Logger().Error("vehicle 转换为json 失败err:%s, data:%+v", err, &vehjsonStruct)
			return
		}
		//g.RedisMission.SetJson(vehicleKey, byteVeh, time.Duration(1800)*time.Second) // todo 更新json, 将images,send进去

		err = udp.SendUdp(g.Config().UdpInfo.UdpAddr, byteVeh)
		if err != nil {
			g.Logger().Error("发送 vehicle json 失败 err:%s, data:%+v", err, &vehjsonStruct)
			return
		}
		g.Logger().Debug("发送 vehicle json 成功, data:%+v", &vehjsonStruct)
		g.Redis.SetAdd(t.TRAIN_SERIAL, vehicleKey) // todo 做记录 发达过车厢的集合

		// todo 发送过后 更新发送状态
		//vehjsonStruct.Send=true
		//byteVeh, err = json.Marshal(&vehjsonStruct)
		//if err != nil {
		//	g.Logger().Error("vehicle 转换为json 失败err:%s, data:%+v", err, &vehjsonStruct)
		//	return
		//}
		g.RedisMission.SetJson(vehicleKey, &vehjsonStruct, time.Duration(1800)*time.Second) // todo 更新json, 将images,send进去
	}
}

func GetImages(vehicleKey string, vehjsonStruct *VehicleJson) (tmpSlice []string) {
	//sql := fmt.Sprintf("select train_serial,vehicle_serial,img_index,img_name,vehicle_order, t.rowid  from SYS1 t "+
	//	"here vehicle_serial is not null and vehicle_order is not null and vehicle_order = %s and vehicle_serial='%s'", vehjsonStruct.VehicleOrder, vehicleKey)
	sql := fmt.Sprintf("select img_index,img_name  from SYS1 t " +
		"where vehicle_serial is not null and vehicle_order is not null and vehicle_serial='%s'", vehicleKey) //todo 应该使用 vehicle_order字段筛选
	if result, b := g.GetSqlData(sql); b {
		//tmpSlice := make([]string, 0)
		for _, v := range result {
			tmpSlice = append(tmpSlice, v["IMG_NAME"])

		}
	} else {
		g.Logger().Error("GetImage sql vehicleKey:%s Error\n", vehicleKey)
	}
	return
}
