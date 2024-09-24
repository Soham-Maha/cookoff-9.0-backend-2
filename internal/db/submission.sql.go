// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: submission.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createSubmission = `-- name: CreateSubmission :exec
INSERT INTO submissions (id, user_id, question_id, language_id)
VALUES ($1, $2, $3, $4)
`

type CreateSubmissionParams struct {
	ID         uuid.UUID
	UserID     uuid.NullUUID
	QuestionID uuid.UUID
	LanguageID int32
}

func (q *Queries) CreateSubmission(ctx context.Context, arg CreateSubmissionParams) error {
	_, err := q.db.Exec(ctx, createSubmission,
		arg.ID,
		arg.UserID,
		arg.QuestionID,
		arg.LanguageID,
	)
	return err
}

const getSubmission = `-- name: GetSubmission :one
SELECT 
    testcases_passed, 
    testcases_failed 
FROM 
    submissions 
WHERE 
    id = $1
`

type GetSubmissionRow struct {
	TestcasesPassed pgtype.Int4
	TestcasesFailed pgtype.Int4
}

func (q *Queries) GetSubmission(ctx context.Context, id uuid.UUID) (GetSubmissionRow, error) {
	row := q.db.QueryRow(ctx, getSubmission, id)
	var i GetSubmissionRow
	err := row.Scan(&i.TestcasesPassed, &i.TestcasesFailed)
	return i, err
}

const getSubmissionByID = `-- name: GetSubmissionByID :one
SELECT
    id,
    question_id,
    testcases_passed,
    testcases_failed,
    description,
    user_id
FROM submissions
WHERE id = $1
`

type GetSubmissionByIDRow struct {
	ID              uuid.UUID
	QuestionID      uuid.UUID
	TestcasesPassed pgtype.Int4
	TestcasesFailed pgtype.Int4
	Description     *string
	UserID          uuid.NullUUID
}

func (q *Queries) GetSubmissionByID(ctx context.Context, id uuid.UUID) (GetSubmissionByIDRow, error) {
	row := q.db.QueryRow(ctx, getSubmissionByID, id)
	var i GetSubmissionByIDRow
	err := row.Scan(
		&i.ID,
		&i.QuestionID,
		&i.TestcasesPassed,
		&i.TestcasesFailed,
		&i.Description,
		&i.UserID,
	)
	return i, err
}

const getSubmissionStatusByID = `-- name: GetSubmissionStatusByID :one
SELECT
    status
FROM submissions
WHERE id = $1
`

func (q *Queries) GetSubmissionStatusByID(ctx context.Context, id uuid.UUID) (*string, error) {
	row := q.db.QueryRow(ctx, getSubmissionStatusByID, id)
	var status *string
	err := row.Scan(&status)
	return status, err
}

const getSubmissionsWithRoundByUserId = `-- name: GetSubmissionsWithRoundByUserId :many
SELECT q.round, q.title, q.description, s.id, s.question_id, s.testcases_passed, s.testcases_failed, s.runtime, s.submission_time, s.language_id, s.description, s.memory, s.user_id, s.status
FROM submissions s
INNER JOIN questions q ON s.question_id = q.id
WHERE s.user_id = $1
`

type GetSubmissionsWithRoundByUserIdRow struct {
	Round           int32
	Title           string
	Description     string
	ID              uuid.UUID
	QuestionID      uuid.UUID
	TestcasesPassed pgtype.Int4
	TestcasesFailed pgtype.Int4
	Runtime         pgtype.Numeric
	SubmissionTime  pgtype.Timestamp
	LanguageID      int32
	Description_2   *string
	Memory          pgtype.Numeric
	UserID          uuid.NullUUID
	Status          *string
}

func (q *Queries) GetSubmissionsWithRoundByUserId(ctx context.Context, userID uuid.NullUUID) ([]GetSubmissionsWithRoundByUserIdRow, error) {
	rows, err := q.db.Query(ctx, getSubmissionsWithRoundByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubmissionsWithRoundByUserIdRow
	for rows.Next() {
		var i GetSubmissionsWithRoundByUserIdRow
		if err := rows.Scan(
			&i.Round,
			&i.Title,
			&i.Description,
			&i.ID,
			&i.QuestionID,
			&i.TestcasesPassed,
			&i.TestcasesFailed,
			&i.Runtime,
			&i.SubmissionTime,
			&i.LanguageID,
			&i.Description_2,
			&i.Memory,
			&i.UserID,
			&i.Status,
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

const getTestCases = `-- name: GetTestCases :many
SELECT id, expected_output, memory, input, hidden, runtime, question_id 
FROM testcases
WHERE question_id = $1
  AND (CASE WHEN $2 = TRUE THEN hidden = FALSE ELSE TRUE END)
`

type GetTestCasesParams struct {
	QuestionID uuid.UUID
	Column2    interface{}
}

func (q *Queries) GetTestCases(ctx context.Context, arg GetTestCasesParams) ([]Testcase, error) {
	rows, err := q.db.Query(ctx, getTestCases, arg.QuestionID, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Testcase
	for rows.Next() {
		var i Testcase
		if err := rows.Scan(
			&i.ID,
			&i.ExpectedOutput,
			&i.Memory,
			&i.Input,
			&i.Hidden,
			&i.Runtime,
			&i.QuestionID,
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

const updateDescriptionStatus = `-- name: UpdateDescriptionStatus :exec
UPDATE submissions
SET description = $1
WHERE id = $2
`

type UpdateDescriptionStatusParams struct {
	Description *string
	ID          uuid.UUID
}

func (q *Queries) UpdateDescriptionStatus(ctx context.Context, arg UpdateDescriptionStatusParams) error {
	_, err := q.db.Exec(ctx, updateDescriptionStatus, arg.Description, arg.ID)
	return err
}

const updateSubmission = `-- name: UpdateSubmission :exec
UPDATE submissions
SET testcases_passed = $1, testcases_failed = $2, runtime = $3, memory = $4
WHERE id = $5
`

type UpdateSubmissionParams struct {
	TestcasesPassed pgtype.Int4
	TestcasesFailed pgtype.Int4
	Runtime         pgtype.Numeric
	Memory          pgtype.Numeric
	ID              uuid.UUID
}

func (q *Queries) UpdateSubmission(ctx context.Context, arg UpdateSubmissionParams) error {
	_, err := q.db.Exec(ctx, updateSubmission,
		arg.TestcasesPassed,
		arg.TestcasesFailed,
		arg.Runtime,
		arg.Memory,
		arg.ID,
	)
	return err
}

const updateSubmissionStatus = `-- name: UpdateSubmissionStatus :exec
UPDATE submissions
SET status = $1
WHERE id = $2
`

type UpdateSubmissionStatusParams struct {
	Status *string
	ID     uuid.UUID
}

func (q *Queries) UpdateSubmissionStatus(ctx context.Context, arg UpdateSubmissionStatusParams) error {
	_, err := q.db.Exec(ctx, updateSubmissionStatus, arg.Status, arg.ID)
	return err
}
