package users

import (
	"context"
	"errors"
	"strconv"

	"bitbucket.org/faceit/app"
	"gorm.io/gorm"
)

func NewMysqlRepo(db *gorm.DB) *MysqlRepo {
	return &MysqlRepo{db: db}
}

type MysqlRepo struct {
	db *gorm.DB
}

func (r *MysqlRepo) Get(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.First(&user, "email=?", email).Error; err != nil {
		return nil, app.WrapError(err)
	}

	return &user, nil
}

func (r *MysqlRepo) Create(ctx context.Context, user *User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return app.WrapError(err)
	}

	return nil
}

func (r *MysqlRepo) Update(ctx context.Context, user *User) error {
	if err := r.db.Updates(&user).Error; err != nil {
		return app.WrapError(err)
	}
	return nil
}

func (r *MysqlRepo) Delete(ctx context.Context, userID string) error {
	var user *User
	if err := r.db.Delete(&user, "id=?", userID).Error; err != nil {
		return app.WrapError(err)
	}
	return nil
}

type Filter struct {
	Page     string
	PageSize string
	Country  string
}

func (r *MysqlRepo) List(ctx context.Context, filter Filter) ([]User, error) {
	var users []User

	var query = r.db.Scopes(paginate(filter))
	if filter.Country != "" {
		query = query.Where("country=?", filter.Country)
	}
	if err := query.Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorNotfound
		}
		return nil, app.WrapError(err)
	}

	return users, nil
}

func paginate(filter Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		page, _ := strconv.Atoi(filter.Page)
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(filter.PageSize)
		switch {
		case pageSize > 20:
			pageSize = 20
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
