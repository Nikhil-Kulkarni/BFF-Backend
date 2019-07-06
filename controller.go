package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// Controller controls requests
type Controller struct {
	repository *Repository
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// CorsHandler cors pre-flight options handler
func (c *Controller) CorsHandler(w http.ResponseWriter, r *http.Request) {
	writeCorsHeaders(&w)
	return
}

// Login creates user profile if needed and returns scores - Called by mobile client
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Invalid request"})
		return
	}

	if request.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"No user id found"})
		return
	}

	ctx := context.Background()
	questions, err := c.repository.FetchAllQuestions(ctx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Something went wrong"})
		return
	}

	scores, err := c.repository.FetchScores(ctx, request.UserID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Something went wrong"})
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{scores, questions})
}

// SubmitScore submits a score - Called by web client
func (c *Controller) SubmitScore(w http.ResponseWriter, r *http.Request) {
	writeCorsHeaders(&w)
	var request Score
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Invalid request"})
		return
	}

	if request.UserID == "" {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"No user id found"})
		return
	}

	timestamp := time.Now().Unix()
	ctx := context.Background()
	submitScore := SubmittedScore{request.Name, request.Value, timestamp}
	err = c.repository.SubmitScore(ctx, submitScore, request.UserID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Something went wrong"})
		return
	}

	json.NewEncoder(w).Encode(SuccessResponse{true})
}

// FetchScores fetches top scores for id - called by mobile client
func (c *Controller) FetchScores(w http.ResponseWriter, r *http.Request) {
	var request FetchScoresRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Invalid request"})
		return
	}

	ctx := context.Background()
	scores, err := c.repository.FetchScores(ctx, request.UserID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Something went wrong"})
		return
	}

	json.NewEncoder(w).Encode(FetchScoresResponse{scores})
}

// FetchAllQuestions fetches all question - Called by mobile client
func (c *Controller) FetchAllQuestions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	questions, err := c.repository.FetchAllQuestions(ctx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Something went wrong"})
		return
	}

	json.NewEncoder(w).Encode(FetchAllQuestionsResponse{questions})
}

// CreateTest creates a test - Called by mobile client
func (c *Controller) CreateTest(w http.ResponseWriter, r *http.Request) {
	var test Test
	err := json.NewDecoder(r.Body).Decode(&test)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Invalid request"})
		return
	}

	if test.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"No user id found"})
		return
	}

	if len(test.SubmittedQuestions) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"No questions in request"})
		return
	}

	ctx := context.Background()
	testID := randomID()
	err = c.repository.CreateTest(ctx, testID, test)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Failed to create test"})
		return
	}

	json.NewEncoder(w).Encode(CreateTestResponse{"https://playbff.herokuapp.com/" + testID})
}

// FetchTest fetches a test - Called by web client
func (c *Controller) FetchTest(w http.ResponseWriter, r *http.Request) {
	writeCorsHeaders(&w)
	var request FetchTestRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Invalid request"})
		return
	}

	if request.TestID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"No test id"})
		return
	}

	ctx := context.Background()

	test, err := c.repository.FetchTest(ctx, request)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Something went wrong"})
		return
	}

	questions, err := c.repository.FetchQuestions(ctx, test)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{"Something went wrong"})
		return
	}

	json.NewEncoder(w).Encode(FetchTestResponse{questions, test.UserID})
}

func randomID() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func writeCorsHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://playbff.herokuapp.com, https://playbff.herokuapp.com")
	(*w).Header().Add("Access-Control-Allow-Methods", "POST")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type")
}
