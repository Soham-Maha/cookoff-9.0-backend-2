package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/google/uuid"
)

type resp struct {
	Result []submission.Judgeresp `json:"result"`
}

func RunCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req subreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	question_id, _ := uuid.Parse(req.QuestionID)
	userID, _ := auth.GetUserID(w, r)

	qualified, err := auth.VerifyRound(ctx, userID, question_id)
	if err != nil {
		httphelpers.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if !qualified {
		httphelpers.WriteError(w, http.StatusForbidden, "User is not qualified for this round")
		return
	}

	testcases, err := database.Queries.GetTestCases(ctx, db.GetTestCasesParams{QuestionID: question_id, Column2: true})
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, fmt.Sprintf("error getting test cases for question_id %d: %v", question_id, err))
		return
	}

	judge0URL, err := url.Parse(JUDGE0_URI + "/submissions/")
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing Judge0 URL: %v", err))
		return
	}
	params := url.Values{}
	params.Add("base64_encoded", "true")
	params.Add("wait", "true")
	judge0URL.RawQuery = params.Encode()

	var payload submission.Submission
	response := resp{
		Result: make([]submission.Judgeresp, len(testcases)),
	}

	runtime_mut, err := submission.RuntimeMut(req.LanguageID)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	for i, testcase := range testcases {
		runtime, _ := testcase.Runtime.Float64Value()
		payload = submission.Submission{
			LanguageID: req.LanguageID,
			SourceCode: submission.B64(req.SourceCode),
			Input:      submission.B64(*testcase.Input),
			Output:     submission.B64(*testcase.ExpectedOutput),
			Runtime:    runtime.Float64 * float64(runtime_mut),
		}

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			httphelpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("error marshaling payload: %v", err))
			return
		}

		result, err := http.Post(judge0URL.String(), "application/json", bytes.NewBuffer(payloadJSON))
		if err != nil {
			logger.Errof("Error sending request to Judge0: %v", err)
			httphelpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Error sending request to Judge0: %v", err))
			return
		}
		defer result.Body.Close()

		var data submission.Judgeresp
		data.TestCaseID = testcase.ID.String()
		if err = json.NewDecoder(result.Body).Decode(&data); err != nil {
			logger.Errof("Error decoding response from Judge0: %v", err)
			httphelpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding response from Judge0: %v", err))
			return
		}

		data.CompilerOutput, _ = submission.DecodeB64(data.CompilerOutput)
		response.Result[i] = data
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errof("Error encoding response: %v", err)
		httphelpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Error encoding response: %v", err))
	}
}
