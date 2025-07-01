package data

import (
	"math/rand"
	"nancalacc/internal/data/models"
	"time"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	// 确保只执行一次
	if db.Migrator().HasTable(&models.Account{}) && db.Find(&models.Account{}).RowsAffected > 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())

	// 使用事务确保数据完整性
	return db.Transaction(func(tx *gorm.DB) error {
		// // 1. 创建角色种子数据
		// roles := []models.Role{
		// 	{Name: "admin", Description: "Administrator"},
		// 	{Name: "user", Description: "Regular User"},
		// }
		// if err := tx.Create(&roles).Error; err != nil {
		// 	return err
		// }

		// 2. 创建用户种子数据
		users := []models.Account{
			{
				ID:       1,
				Username: "user1",
				Email:    "user1@example.com",
				Phone:    "13800000000",
				Password: hashPassword("admin123"), // 实际项目中应该加密
				//RoleID:   roles[0].ID,
				Status:    1,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			{
				ID:        2,
				Username:  "user2",
				Email:     "user2@example.com",
				Phone:     "13800000001",
				Password:  hashPassword("user456"),
				Status:    1,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
		}
		if err := tx.Create(&users).Error; err != nil {
			return err
		}

		// 3. 可以继续添加其他种子数据...

		return nil
	})
}

func hashPassword(password string) string {
	// 实际项目中应该使用bcrypt等加密算法
	return password // 这里简化处理
}
