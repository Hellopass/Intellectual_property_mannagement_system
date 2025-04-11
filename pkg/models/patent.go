package models

import "intellectual_property/pkg/utils"

// Patent 专利信息表
type Patent struct {
	ID                uint     `json:"id" gorm:"primaryKey"`                         //id --主键
	ApplicationNo     string   `json:"application_no" gorm:"application_no"`         //申请号
	Name              string   `json:"name" gorm:"name"`                             //专利名称
	Inventor          string   `json:"inventor" gorm:"inventor"`                     //发明人
	Type              string   `json:"type" gorm:"type" `                            //专利类型
	TechnicalField    string   `json:"technical_field" gorm:"technical_field"`       //技术领域
	TechnicalSolution string   `json:"technical_solution" gorm:"technical_solution"` //技术方案
	RelatedFiles      []string `json:"related_files" gorm:"related_files" `          //相关文件
}

// CreatePatent 创建新专利
func CreatePatent(patent *Patent) error {
	return utils.DB.Create(patent).Error
}

// GetPatent 获取专利
func GetPatent(id uint) (*Patent, error) {
	var patent Patent
	if err := utils.DB.First(&patent, id).Error; err != nil {
		return nil, err
	}
	return &patent, nil
}

// UpdatePatent 更新专利
func UpdatePatent(id uint, patent *Patent) error {
	return utils.DB.Model(&Patent{}).Where("id = ?", id).Updates(patent).Error
}

// DeletePatent 删除专利
func DeletePatent(id uint) error {
	return utils.DB.Delete(&Patent{}, id).Error
}

// ListPatents 列出所有专利
func ListPatents() ([]Patent, error) {
	var patents []Patent
	if err := utils.DB.Find(&patents).Error; err != nil {
		return nil, err
	}
	return patents, nil
}
