package repo

import "example.com/app/domain"

type AuthRepo interface {
	Login(username string, password string, ip string, ips []string) (*domain.Admin, string, error)
}

