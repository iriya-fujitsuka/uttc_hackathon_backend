package dao

import (
	"log"
	"uttc_hackathon_backend/models"
)

func GetCommunities() ([]models.Community, error) {
	rows, err := db.Query("SELECT id, name, created_at FROM communities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []models.Community
	for rows.Next() {
		var community models.Community
		if err := rows.Scan(&community.ID, &community.Name, &community.CreatedAt); err != nil {
			return nil, err
		}
		communities = append(communities, community)
	}
	return communities, nil
}

func AddCommunity(name string) (int, error) {
	result, err := db.Exec("INSERT INTO communities (name) VALUES (?)", name)
	if err != nil {
		log.Printf("Error adding community: %v", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return 0, err
	}
	return int(id), nil
}

func DeleteCommunity(id int) error {
	_, err := db.Exec("DELETE FROM communities WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting community: %v", err)
		return err
	}
	return nil
}