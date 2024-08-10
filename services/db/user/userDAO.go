package user

import (
	"database/sql"

	"github.com/Jack-Gitter/tunesEmailService/services/db/general"
)

type UserDAO struct { }

type IUserDAO interface {
    GetUserFollowerEmails(executor general.QueryExecutor, spotifyID string) ([]string, error)
    GetUsername(executor general.QueryExecutor, spotifyID string) (string, error)
}

func(u *UserDAO) GetUsername(executor general.QueryExecutor, spotifyID string) (string, error) {

    query := ` SELECT users.username FROM users WHERE users.spotifyid = $1`

    row := executor.QueryRow(query, spotifyID)

    username := ""
    err := row.Scan(&username)

    if err != nil {
        return "", err
    }

    return username, nil

}

func(u *UserDAO) GetUserFollowerEmails(executor general.QueryExecutor, spotifyID string) ([]string, error) {

    query := ` SELECT users.email
                FROM followers 
                INNER JOIN  users 
                ON users.spotifyid = followers.follower 
                WHERE followers.userfollowed = $1`

    rows, err := executor.Query(query, spotifyID)

    if err != nil {
        return []string{}, err
    }

    emails := []string{}

    for rows.Next() {
        email := sql.NullString{}
        err := rows.Scan(&email)
        if err != nil {
            return nil, err
        }
        emails = append(emails, email.String)
    }

    return emails, nil

}
