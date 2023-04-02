package repository

import (
	models "URLShortner/Repository/Models"
	"database/sql"
	"fmt"
)

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	var linkRepo LinkRepository = LinkRepository{db: db}
	return &linkRepo
}

func (ins *LinkRepository) List() ([]*models.Link, error) {
	rows, err := ins.db.Query("SELECT * FROM links;")
	if err != nil {
		return nil, err
	}
	var links []*models.Link = []*models.Link{}
	for rows.Next() {
		var link models.Link = models.Link{}
		if err := rows.Scan(&link.LinkId, &link.UserId, &link.Name, &link.Url, &link.ShortLink); err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, nil
}

func (ins *LinkRepository) GetById(id int) (*models.Link, error) {
	var link models.Link = models.Link{}
	var row *sql.Row = ins.db.QueryRow("SELECT * FROM links WHERE LinkId=?;", id)
	if err := row.Scan(&link.LinkId, &link.UserId, &link.Name, &link.Url, &link.ShortLink); err != nil {
		return nil, err
	}

	return &link, nil
}

func (ins *LinkRepository) GetByUserId(userId int) ([]*models.Link, error) {
	rows, err := ins.db.Query("SELECT * FROM links WHERE UserId=?;", userId)
	if err != nil {
		return nil, err
	}
	var links []*models.Link = []*models.Link{}
	for rows.Next() {
		var link models.Link = models.Link{}
		if err := rows.Scan(&link.LinkId, &link.UserId, &link.Name, &link.Url, &link.ShortLink); err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, nil
}

func (ins *LinkRepository) CreateLink(link *models.Link) error {

	var isExist int = 0
	var row *sql.Row = ins.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE UserId=?);", link.UserId)
	if err := row.Scan(&isExist); err != nil {
		return err
	}

	if isExist == 0 {
		err := fmt.Errorf("wrong UserId: %d", link.UserId)
		return err
	}

	trans, err := ins.db.Begin()
	if err != nil {
		return err
	}
	result, err := trans.Exec("INSERT INTO links (UserId, Name, Url, ShortLink) VALUES (?, ?, ?, ?);", link.UserId, link.Name, link.Url, link.ShortLink)
	if err != nil {
		return err
	}

	linkId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	link.LinkId = linkId

	return trans.Commit()
}

func (ins *LinkRepository) Update(link *models.Link) error {
	trans, err := ins.db.Begin()
	if err != nil {
		return err
	}

	if _, err := trans.Exec("UPDATE links SET Name=?, Url=?, ShortLink=? WHERE LinkId=?;", link.Name, link.Url, link.ShortLink, link.LinkId); err != nil {
		return err
	}

	return trans.Commit()
}

func (ins *LinkRepository) Delete(id int) error {
	trans, err := ins.db.Begin()
	if err != nil {
		return err
	}

	if _, err := trans.Exec("DELETE FROM links WHERE LinkId=?;", id); err != nil {
		return err
	}
	return trans.Commit()
}

func (ins *LinkRepository) GetLinkByShortLink(shorted string) (*models.Link, error) {
	var link models.Link = models.Link{}
	var row *sql.Row = ins.db.QueryRow("SELECT * FROM links WHERE ShortLink=?;", shorted)
	if err := row.Scan(&link.LinkId, &link.UserId, &link.Name, &link.Url, &link.ShortLink); err != nil {
		return nil, err
	}

	return &link, nil
}
