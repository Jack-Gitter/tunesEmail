package user

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type UserService struct {
    DB *sql.DB
    UserDAO IUserDAO
}

type IUserService interface {
    GetUserFollowerEmails(spotifyID string) ([]string, error)
}

func(us *UserService) GetUserFollowerEmails(spotifyID string) ([]string, error) {
    return us.UserDAO.GetUserFollowerEmails(us.DB, spotifyID)
}

