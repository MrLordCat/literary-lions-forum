package db

import (
	"database/sql"
)

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
func UpdateUserKarma(db *sql.DB, userID int) error {
	// Получаем суммарное количество лайков для постов пользователя
	var postLikes int
	postLikesQuery := `SELECT COALESCE(SUM(like_type), 0) FROM post_likes WHERE post_id IN (SELECT id FROM posts WHERE author_id = ?)`
	err := db.QueryRow(postLikesQuery, userID).Scan(&postLikes)
	if err != nil {
		return err
	}

	// Получаем суммарное количество лайков для комментариев пользователя
	var commentLikes int
	commentLikesQuery := `SELECT COALESCE(SUM(like_type), 0) FROM comment_likes WHERE comment_id IN (SELECT id FROM comments WHERE author_id = ?)`
	err = db.QueryRow(commentLikesQuery, userID).Scan(&commentLikes)
	if err != nil {
		return err
	}

	// Рассчитываем новую карму пользователя
	karma := postLikes + (commentLikes / 2)

	// Обновляем карму в таблице users
	updateQuery := `UPDATE users SET karma = ? WHERE id = ?`
	_, err = db.Exec(updateQuery, karma, userID)
	if err != nil {
		return err
	}

	return nil
}
