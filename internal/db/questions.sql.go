// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: questions.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createQuestion = `-- name: CreateQuestion :one
INSERT INTO questions (id, description, title, "input_format", points, round, constraints, output_format)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, description, title, input_format, points, round, constraints, output_format
`

type CreateQuestionParams struct {
	ID           uuid.UUID
	Description  *string
	Title        *string
	InputFormat  *string
	Points       pgtype.Int4
	Round        int32
	Constraints  *string
	OutputFormat *string
}

func (q *Queries) CreateQuestion(ctx context.Context, arg CreateQuestionParams) (Question, error) {
	row := q.db.QueryRow(ctx, createQuestion,
		arg.ID,
		arg.Description,
		arg.Title,
		arg.InputFormat,
		arg.Points,
		arg.Round,
		arg.Constraints,
		arg.OutputFormat,
	)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.Description,
		&i.Title,
		&i.InputFormat,
		&i.Points,
		&i.Round,
		&i.Constraints,
		&i.OutputFormat,
	)
	return i, err
}

const deleteQuestion = `-- name: DeleteQuestion :exec
DELETE FROM questions 
WHERE id = $1
`

func (q *Queries) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteQuestion, id)
	return err
}

const getQuestion = `-- name: GetQuestion :one
SELECT id, description, title, input_format, points, round, constraints, output_format FROM questions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetQuestion(ctx context.Context, id uuid.UUID) (Question, error) {
	row := q.db.QueryRow(ctx, getQuestion, id)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.Description,
		&i.Title,
		&i.InputFormat,
		&i.Points,
		&i.Round,
		&i.Constraints,
		&i.OutputFormat,
	)
	return i, err
}

const getQuestionByRound = `-- name: GetQuestionByRound :many
SELECT id, description, title, input_format, points, round, constraints, output_format FROM questions
WHERE round = $1
`

func (q *Queries) GetQuestionByRound(ctx context.Context, round int32) ([]Question, error) {
	rows, err := q.db.Query(ctx, getQuestionByRound, round)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Question
	for rows.Next() {
		var i Question
		if err := rows.Scan(
			&i.ID,
			&i.Description,
			&i.Title,
			&i.InputFormat,
			&i.Points,
			&i.Round,
			&i.Constraints,
			&i.OutputFormat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQuestions = `-- name: GetQuestions :many
SELECT id, description, title, input_format, points, round, constraints, output_format
FROM questions
`

func (q *Queries) GetQuestions(ctx context.Context) ([]Question, error) {
	rows, err := q.db.Query(ctx, getQuestions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Question
	for rows.Next() {
		var i Question
		if err := rows.Scan(
			&i.ID,
			&i.Description,
			&i.Title,
			&i.InputFormat,
			&i.Points,
			&i.Round,
			&i.Constraints,
			&i.OutputFormat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateQuestion = `-- name: UpdateQuestion :exec
UPDATE questions
SET description = $1, title = $2, input_format = $3, points = $4, round = $5, constraints = $6, output_format = $7
WHERE id = $8
`

type UpdateQuestionParams struct {
	Description  *string
	Title        *string
	InputFormat  *string
	Points       pgtype.Int4
	Round        int32
	Constraints  *string
	OutputFormat *string
	ID           uuid.UUID
}

func (q *Queries) UpdateQuestion(ctx context.Context, arg UpdateQuestionParams) error {
	_, err := q.db.Exec(ctx, updateQuestion,
		arg.Description,
		arg.Title,
		arg.InputFormat,
		arg.Points,
		arg.Round,
		arg.Constraints,
		arg.OutputFormat,
		arg.ID,
	)
	return err
}
