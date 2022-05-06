package main

import (
	"testing"
)

func Test_getADDR(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"localhost", ":8080"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getADDR(); got != tt.want {
				t.Errorf("getADDR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mustGetEnv(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustGetEnv(tt.args.key); got != tt.want {
				t.Errorf("mustGetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mysqlDSN(t *testing.T) {
	type args struct {
		dsn string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"localhost:3060", args{"tcp://root:faceittestpassword@localhost:3306/faceit"}, "root:faceittestpassword@tcp(localhost:3306)/faceit?charset=utf8&parseTime=True&loc=Local"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mysqlDSN(tt.args.dsn); got != tt.want {
				t.Errorf("mysqlDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
