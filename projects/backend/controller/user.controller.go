package controller

import (
	"github.com/fes111/rmm/libs/go/models"
	"github.com/fes111/rmm/projects/backend/config"
)

func AddDeviceToUser(id string, token string) error {
	deviceToken, err := GetDeviceToken(token)
	if err != nil {
		return err
	}
	device, err := GetDeviceById(id)
	if err != nil {
		return err
	}
	result := map[string]interface{}{}
	config.Database.Table("user_devices").Where("user_id = ?", deviceToken.UserID).Where("device_id = ?", device.ID).Find(&result)
	if len(result) == 0 {
		config.Database.Table("user_devices").Create(map[string]interface{}{"user_id": deviceToken.UserID, "device_id": device.ID})
	}
	return nil
}

func UpdateUser(user models.User) error {
	err := config.Database.Save(&user).Error
	return err
}

func GetUserById(id uint) (models.User, error) {
	user := models.User{}
	err := config.Database.First(&user, id).Error
	return user, err
}
