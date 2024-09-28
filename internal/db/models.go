// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Question struct {
	ID               uuid.UUID
	Description      string
	Title            string
	InputFormat      []string
	Points           int32
	Round            int32
	Constraints      []string
	OutputFormat     []string
	SampleTestInput  []string
	SampleTestOutput []string
	Explanation      []string
}

type Submission struct {
	ID              uuid.UUID
	QuestionID      uuid.UUID
	TestcasesPassed pgtype.Int4
	TestcasesFailed pgtype.Int4
	Runtime         pgtype.Numeric
	SubmissionTime  pgtype.Timestamp
	SourceCode      string
	LanguageID      int32
	Description     *string
	Memory          pgtype.Numeric
	UserID          uuid.NullUUID
	Status          *string
}

type SubmissionResult struct {
	ID           uuid.UUID
	TestcaseID   uuid.NullUUID
	SubmissionID uuid.UUID
	Runtime      pgtype.Numeric
	Memory       pgtype.Numeric
	Status       string
	Description  *string
}

type Testcase struct {
	ID             uuid.UUID
	ExpectedOutput string
	Memory         pgtype.Numeric
	Input          string
	Hidden         bool
	Runtime        pgtype.Numeric
	QuestionID     uuid.UUID
}

type User struct {
	ID             uuid.UUID
	Email          string
	RegNo          string
	Password       string
	Role           string
	RoundQualified int32
	Score          pgtype.Numeric
	Name           string
	IsBanned       bool
}
