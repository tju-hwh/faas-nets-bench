package handlers

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB
var Nodes []Node
var FcLocates []FcLocate
var Fcs []Fc
var MyPods []MyPod

type MyPod struct {
	ID int `db:"id"`
	FcName string `db:"fcname"`
	PodName string `db:"podname"`
	Locate string `db:"locate"`
}

func (MyPod) TableName() string{
	return "pod"
}

type FcLocate struct {
	ID int `db:"id"`
	FcName string `db:"fc_name"`
	Locate string `db:"locate"`
	Number int `db:"number"`
}
func (FcLocate) TableName() string{
	return "fc_locate"
}

type Fc struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Cpu float64 `db:"cpu"`
	Nodename string `db:"nodename"`
	Image string `db:"image"`
	Delay float64 `db:"delay"`
	Replicas int `db:"replicas"`
}
func (Fc) TableName() string{
	return "fc"
}

type Node struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Cpu float64 `db:"cpu"`
	Mem float64 `db:"mem"`
	Fc string `db:"fc"`
	Rcpu string `db:"rcpu"`
	Rmem string `db:"rmem"`

	CpuFloat []float64
	MemFloat []float64

}

func (Node) TableName() string{
	return "node"
}
