package users

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func mysqlDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		log.Fatalf("mysql dsn parse error: %s", err)
	}
	userPass := u.User.Username()
	if pass, ok := u.User.Password(); ok {
		userPass += ":" + pass
	}
	return fmt.Sprintf("%s@tcp(%s)%s?charset=utf8&parseTime=True&loc=Local", userPass, u.Host, u.Path)
}

func TestMysqlRepo_Create(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx  context.Context
		user *User
	}
	db, _ := gorm.Open(mysql.Open(mysqlDSN("tcp://root:faceittestpassword@localhost:3306/faceit")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	user := &User{
		ID:        "test-safak",
		FirstName: "test name",
		LastName:  "test last name",
		Nickname:  "Safak",
		Password:  "1234656",
		Email:     "test@test.com",
		Country:   "TR",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Create", fields{db: db}, args{ctx: context.Background(), user: user}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MysqlRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlRepo_Delete(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	db, _ := gorm.Open(mysql.Open(mysqlDSN("tcp://root:faceittestpassword@localhost:3306/faceit")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Delete", fields{db: db}, args{ctx: context.Background(), userID: "test-safak"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MysqlRepo{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlRepo_List(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx    context.Context
		filter Filter
	}
	db, _ := gorm.Open(mysql.Open(mysqlDSN("tcp://root:faceittestpassword@localhost:3306/faceit")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	var user = User{
		ID:        "847365e8-8f2a-4a99-8567-ae104c8fca0f",
		FirstName: "ok",
		LastName:  "ok",
		Nickname:  "ok",
		Password:  "123456",
		Email:     "safakciplak1990@gmail.com",
		Country:   "UK",
		CreatedAt: time.Date(2022, time.May, 6, 10, 6, 9, 0, time.Now().UTC().Add(time.Hour*3).Local().Location()),
		UpdatedAt: time.Date(2022, time.May, 6, 10, 6, 9, 0, time.Now().UTC().Add(time.Hour*3).Local().Location()),
	}
	var filter = Filter{
		Page:     "1",
		PageSize: "10",
		Country:  "UK",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []User
		wantErr bool
	}{
		{"List", fields{db: db}, args{ctx: context.Background(), filter: filter}, []User{user}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MysqlRepo{
				db: tt.fields.db,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMysqlRepo_Update(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx  context.Context
		user *User
	}
	db, _ := gorm.Open(mysql.Open(mysqlDSN("tcp://root:faceittestpassword@localhost:3306/faceit")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	_ = &User{
		ID:        "test-safak",
		FirstName: "test name",
		LastName:  "test last name",
		Nickname:  "Safak",
		Password:  "1234656",
		Email:     "test@test.com",
		Country:   "TR",
	}
	editedUser := &User{
		ID:        "edited-test-safak",
		FirstName: "test name",
		LastName:  "test last name",
		Nickname:  "Safak",
		Password:  "1234656",
		Email:     "test@test.com",
		Country:   "TR",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Update", fields{db: db}, args{ctx: context.Background(), user: editedUser}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MysqlRepo{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
