package article

import (
	"gorm.io/gorm"
)

type Article struct {
	ID          uint `gorm:"primary_key,auto_increment"`
	Name        string
	Description string
	Price       float64
}

type NotExistsError struct {
}

func (*NotExistsError) Error() string {
	return "User does not exist"
}

func GetArticle(db *gorm.DB, id uint) (*Article, error) {
	var a Article
	res := db.Find(&a, id)
	if res.Error != nil {
		return nil, &NotExistsError{}
	}
	return &a, res.Error
}

func GetAllArticles(db *gorm.DB) ([]Article, error) {
	var a []Article
	res := db.Find(&a)
	if res.Error != nil {
		return nil, res.Error
	}
	return a, nil
}

func Create(db *gorm.DB, a *Article) (uint, error) {
	err := db.Create(a).Error
	if err != nil {
		return 0, err
	}
	return a.ID, nil
}
func Update(db *gorm.DB, a *Article) (uint, error) {
	err := db.Model(&a).Updates(Article{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Price:       a.Price,
	}).Error
	if err != nil {
		return 0, err
	}
	return a.ID, nil
}
func Delete(db *gorm.DB, a *Article) (bool, error) {
	err := db.Delete(&a).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
