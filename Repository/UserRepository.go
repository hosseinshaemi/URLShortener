package repository

import (
	models "URLShortner/Repository/Models"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	var userRepo UserRepository = UserRepository{db: db}
	return &userRepo
}

func (ins *UserRepository) GetUserLinksByUserId(id int64) ([]models.Link, error) {
	rows, err := ins.db.Query("SELECT * FROM links WHERE UserId=?;", id)
	if err != nil {
		return nil, err
	}

	var links []models.Link = []models.Link{}
	for rows.Next() {
		var link models.Link = models.Link{}
		if err := rows.Scan(&link.LinkId, &link.UserId, &link.Name, &link.Url, &link.ShortLink); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func (ins *UserRepository) List() ([]*models.User, error) {
	rows, err := ins.db.Query("SELECT * FROM users;")
	if err != nil {
		return nil, err
	}
	var users []*models.User = []*models.User{}
	for rows.Next() {
		var user models.User = models.User{}
		if err := rows.Scan(&user.UserId, &user.Firstname, &user.Lastname, &user.Email); err != nil {
			return nil, err
		}
		links, err := ins.GetUserLinksByUserId(user.UserId)
		if err != nil {
			return nil, err
		}
		user.Links = links
		users = append(users, &user)
	}
	return users, nil
}

func (ins *UserRepository) GetById(id int64) (*models.User, error) {
	var user models.User = models.User{}
	var row *sql.Row = ins.db.QueryRow("SELECT * FROM users WHERE UserId=?;", id)
	if err := row.Scan(&user.UserId, &user.Firstname, &user.Lastname, &user.Email); err != nil {
		return nil, err
	}

	links, err := ins.GetUserLinksByUserId(id)
	if err != nil {
		return nil, err
	}

	user.Links = links

	return &user, nil
}

func (ins *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User = models.User{}
	var row *sql.Row = ins.db.QueryRow("SELECT * FROM users WHERE Email=?;", email)
	if err := row.Scan(&user.UserId, &user.Firstname, &user.Lastname, &user.Email); err != nil {
		return nil, err
	}

	links, err := ins.GetUserLinksByUserId(user.UserId)
	if err != nil {
		return nil, err
	}

	user.Links = links

	return &user, nil
}

func (ins *UserRepository) CreateUser(user *models.User) error {
	trans, err := ins.db.Begin()
	if err != nil {
		return err
	}
	result, err := trans.Exec("INSERT INTO users (Firstname, Lastname, Email) VALUES (?, ?, ?);", user.Firstname, user.Lastname, user.Email)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.UserId = userId

	if len(user.Links) != 0 {
		for _, link := range user.Links {
			_, err := trans.Exec("INSERT INTO links (Name, Url, ShortLink, UserId) VALUES (?, ?, ?, ?);", link.Name, link.Url, link.ShortLink, link.UserId)
			if err != nil {
				return err
			}
		}
	}
	return trans.Commit()
}

func (ins *UserRepository) Update(user *models.User) error {
	trans, err := ins.db.Begin()
	if err != nil {
		return err
	}

	if _, err := trans.Exec("UPDATE users SET Firstname=?, Lastname=?, Email=? WHERE UserId=?;", user.Firstname, user.Lastname, user.Email, user.UserId); err != nil {
		return err
	}

	return trans.Commit()
}

func (ins *UserRepository) Delete(id int) error {
	trans, err := ins.db.Begin()
	if err != nil {
		return err
	}

	if _, err := trans.Exec("DELETE FROM users WHERE UserId=?;", id); err != nil {
		return err
	}

	if _, err := trans.Exec("DELETE FROM links WHERE UserId=?;", id); err != nil {
		return err
	}
	return trans.Commit()
}
