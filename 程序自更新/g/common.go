package g

/**
 * @Author: Tao Jun
 * @Since: 2023/3/20
 * @Desc: common.go
**/

var ExitAppChan chan struct{} = make(chan struct{}, 0)
var NewAppFinishChan chan struct{} = make(chan struct{}, 0)
