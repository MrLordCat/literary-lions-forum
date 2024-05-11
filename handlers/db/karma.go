package db

import "database/sql"

func CalculateUserKarma(db *sql.DB, userID int) (int, error) {
	// Карма от постов
	var postKarma int
	err := db.QueryRow(`
        SELECT COALESCE(SUM(like_type), 0) FROM post_likes WHERE user_id = ?
    `, userID).Scan(&postKarma)
	if err != nil {
		return 0, err
	}

	// Карма от комментариев
	var commentKarma int
	err = db.QueryRow(`
        SELECT COALESCE(SUM(like_type), 0) FROM comment_likes WHERE user_id = ?
    `, userID).Scan(&commentKarma)
	if err != nil {
		return 0, err
	}

	// Конвертируем карму комментариев
	commentKarma /= 2

	// Общая карма
	totalKarma := postKarma + commentKarma
	if totalKarma < 0 {
		totalKarma = 0 // Карма не может быть отрицательной
	}

	return totalKarma, nil
}
