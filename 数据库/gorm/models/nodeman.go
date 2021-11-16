package models

type Hostinfo struct {
	ID         int       `json:"id" `
	Uuid       string    `json:"uuid"`
	Ip         string    `json:"ip"`
	Passwd     string    `json:"-"`
	Username   string    `json:"username"`
	Port       int       `json:"port"`
	Bizid      int       `json:"bizid" gorm:"column:bizid"`
	Createtime string    `json:"createtime"`
	Run        []Hostrun `gorm:"ForeignKey:Host_id;AssociationForeignKey:ID"` //这边是一,其它表 host_id 关联本表ID列
}

type Hostrun struct {
	ID      int    `json:"id"`
	Step    int    `json:"step"`
	Result  string `json:"result" gorm:"column:result"`
	Runtime string `json:"runtime"`
	Status  int    `json:"status"`
	Host_id int    `json:"host_id"`
}
