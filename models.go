package main

// Exception message
type Exception struct {
	Message string `json:"message"`
}

// Choice an answer choice
type Choice struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// CreateTestResponse url for created test
type CreateTestResponse struct {
	TestURL string `json:"url"`
}

// FetchAllQuestionsResponse all questions
type FetchAllQuestionsResponse struct {
	Questions []Question `json:"questions"`
}

// FetchScoresRequest request for fetching scores
type FetchScoresRequest struct {
	UserID string `json:"userId"`
}

// FetchScoresResponse scores for user
type FetchScoresResponse struct {
	Scores []SubmittedScore `json:"scores"`
}

// FetchTestRequest fetch test request
type FetchTestRequest struct {
	TestID string `json:"testId"`
}

// FetchTestResponse questions and answers
type FetchTestResponse struct {
	Questions []QuestionWithCorrectChoice `json:"questions"`
	UserID    string                      `json:"userId"`
}

// LoginRequest request to login
type LoginRequest struct {
	UserID string `json:"userId"`
}

// LoginResponse scores and questions
type LoginResponse struct {
	Scores    []SubmittedScore `json:"scores"`
	Questions []Question       `json:"questions"`
}

// Question a question on a test
type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Choices []Choice `json:"choices"`
}

// QuestionWithCorrectChoice question with correct choice selected by creator
type QuestionWithCorrectChoice struct {
	Text            string   `json:"text"`
	Choices         []Choice `json:"choices"`
	CorrectChoiceID string   `json:"correctChoiceId"`
}

// Score a score
type Score struct {
	Name   string `json:"name"`
	Value  int    `json:"value"`
	UserID string `json:"userId"`
}

// SubmittedQuestion question submitted by creator
type SubmittedQuestion struct {
	QuestionID       string `json:"id" firestore:"id,omitempty"`
	SelectedChoiceID string `json:"selectedChoiceId" firestore:"selectedChoiceId,omitempty"`
}

// SubmittedScore score submitted by test taker
type SubmittedScore struct {
	Name      string `json:"name" firestore:"name,omitempty"`
	Value     int    `json:"value" firestore:"value,omitempty"`
	Timestamp int64  `json:"timestamp" firestore:"timestamp,omitempty"`
}

// SuccessResponse success response
type SuccessResponse struct {
	Success bool `json:"success"`
}

// Test a test
type Test struct {
	SubmittedQuestions []SubmittedQuestion `json:"submittedQuestions" firestore:"questions,omitempty"`
	UserID             string              `json:"userId" firestore:"userId,omitempty"`
}
