package models

import (
	"fmt"
	"intellectual_property/pkg/utils"
	"os"
	"time"
)

// 设定专利状态码
const (
	AwaitingApproval = iota + 1 //待提交
	UnderReview                 //审核中
	Authorized                  //已授权
	Rejected                    //已驳回
)

// Patent 专利信息表
type Patent struct {
	Id          int       `json:"id" gorm:"id;primaryKey;autoIncrement;type:bigint"` //专利id
	ApplyData   time.Time `json:"apply_data" gorm:"apply_data;type:date"`            //申请日期
	ApplyNo     string    `json:"apply_no" gorm:"apply_no;type:varchar(20)"`         //申请号
	PatentName  string    `json:"patent_name" gorm:"patent_name;type:varchar(50)"`   //专利名称
	WarrantDate time.Time `json:"warrant_date" gorm:"warrant_date;type:date"`        //授权日期
	PatentType  string    `json:"patent_type" gorm:"patent_type;type:varchar(20)"`   //专利类型
	UserID      int       `json:"user_id" gorm:"user_id;type:bigint"`                //发明者id ,也是user_id
	User        User      //发明人信息
	Status      int       `json:"status" gorm:"status;type:int"` //专利状态
}

// CreatePatentTable 迁移user表
func CreatePatentTable() {
	err := utils.DB.AutoMigrate(&Patent{}, &User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

// CreatePatent 新建专利申请
// 单个申请
func CreatePatent(patent *Patent) error { return utils.DB.Create(patent).Error }

// CreatePatentBatch  批量新建专利申请
// 多个申请
func CreatePatentBatch(patent []*Patent) error { return utils.DB.Create(patent).Error }

// GetPatentByID 专利查询
func GetPatentByID(patentID int) (*Patent, error) {
	var patent Patent
	if err := utils.DB.Find(&patent, patentID).Error; err != nil {
		return nil, err
	}
	return &patent, nil
}

// GetPatentByInID 专利查询根据user_id
func GetPatentByInID(InId int) ([]Patent, error) {
	var patent []Patent
	if err := utils.DB.Where("user_id=?", InId).Find(&patent).Error; err != nil {
		return nil, err
	}
	for _, v := range patent {
		id, err := GetUserByID(v.UserID)
		if err != nil {
			return nil, err
		}
		v.User = id
	}
	return patent, nil
}

// GetPatentInformation 查询有所专利信息
func GetPatentInformation() ([]Patent, error) {
	var patents []Patent
	if err := utils.DB.Find(&patents).Error; err != nil {
		return nil, err
	}

	//在根据id查询用户信息

	return patents, nil
}

// UpdatePatent 更新专利信息
func UpdatePatent(patent *Patent) error {
	return utils.DB.Model(&Patent{}).Where("id = ? ", patent.Id).Select("patent_name", "warrant_date", "status").Error
}

// DeletePatent 删除专利信息
func DeletePatent(applyNo string) error {
	return utils.DB.Where("apply_no=?", applyNo).Delete(&Patent{}).Error
}

// FindPatentFuzzy 模糊查询
func FindPatentFuzzy(keyword string, status int) ([]Patent, error) {
	db := utils.DB.Model(&Patent{}).
		Preload("User").
		Joins("LEFT JOIN users ON users.id = patents.user_id")
	if keyword != "" {
		db = db.Where("patents.apply_no like ? OR patents.patent_name like ? OR users.user_name like ? ", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status > 0 {
		db = db.Where("patents.status =?", status)
	}
	var patents []Patent
	if err := db.Find(&patents).Error; err != nil {
		return nil, err
	}
	return patents, nil
}

// GetPatentFile 根据专利号拿到所有文件--返回所有nginx地址前端再次请求
func GetPatentFile(applyNo string) ([]string, error) {
	var ans []string
	baseurl := utils.NgX.LocationDocs + "/" + applyNo //文件夹
	dir, err := os.ReadDir(baseurl)
	if err != nil {
		return nil, err
	}
	for _, file := range dir {
		ans = append(ans, "/docs"+"/"+applyNo+"/"+file.Name())
	}
	return ans, nil
}
