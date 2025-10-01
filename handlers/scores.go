package handlers

import (
	"database/sql"
	"reports-api/constants"
	"reports-api/db"
	"reports-api/models"
	"reports-api/utils"

	"github.com/gofiber/fiber/v2"
)

func ListScoresHandler(c *fiber.Ctx) error {
	query := `SELECT department_id, year, month, score FROM scores`
	rows, err := db.DB.Query(query)
	if err != nil {
		return utils.ErrorResponse(c, 500, constants.ErrorFailedToQuery+" scores")
	}
	defer rows.Close()

	var scores []models.Score
	for rows.Next() {
		var score models.Score
		err := rows.Scan(&score.DepartmentID, &score.Year, &score.Month, &score.Score)
		if err != nil {
			return utils.ErrorResponse(c, 500, "Failed to scan score")
		}
		scores = append(scores, score)
	}
	if err := rows.Err(); err != nil {
		return utils.ErrorResponse(c, 500, "Row error")
	}

	return c.JSON(fiber.Map{"scores": scores})
}

func GetScoreDetailHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	query := `SELECT department_id, year, month, score FROM scores WHERE department_id = ?`
	row := db.DB.QueryRow(query, id)

	var score models.Score
	err := row.Scan(&score.DepartmentID, &score.Year, &score.Month, &score.Score)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrorResponse(c, 404, constants.ErrorNotFound)
		}
		return utils.ErrorResponse(c, 500, constants.ErrorFailedToQuery+" score")
	}

	return utils.SuccessResponse(c, score)
}

func UpdateScoreHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var score models.ScoreUpdateRequest
	err := c.BodyParser(&score)
	if err != nil {
		return utils.ErrorResponse(c, 400, constants.MessageInvalidRequest)
	}

	var query string
	var args []interface{}

	if score.Year > 0 && score.Month > 0 {
		query = `UPDATE scores SET score = ? WHERE department_id = ? AND year = ? AND month = ?`
		args = []interface{}{score.Score, id, score.Year, score.Month}
	} else {
		query = `UPDATE scores SET score = ? WHERE department_id = ?`
		args = []interface{}{score.Score, id}
	}

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		return utils.ErrorResponse(c, 500, constants.ErrorFailedToUpdate+" score")
	}

	return utils.UpdatedResponse(c)
}

func DeleteScoreHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var score models.Score
	err := c.BodyParser(&score)

	var query string
	var args []interface{}

	if err == nil && score.Year > 0 && score.Month > 0 {
		query = `DELETE FROM scores WHERE department_id = ? AND year = ? AND month = ?`
		args = []interface{}{id, score.Year, score.Month}
	} else {
		query = `DELETE FROM scores WHERE department_id = ?`
		args = []interface{}{id}
	}

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		return utils.ErrorResponse(c, 500, constants.ErrorFailedToDelete+" score")
	}

	return utils.DeletedResponse(c)
}
