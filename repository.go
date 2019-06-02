package main

import (
	"context"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

// Repository used to db actions
type Repository struct {
	client *firestore.Client
}

// SubmitScore submits a score
func (r *Repository) SubmitScore(ctx context.Context, score SubmittedScore, userID string) error {
	_, _, err := r.client.Collection("users").Doc(userID).Collection("scores").Add(ctx, score)
	if err != nil {
		return err
	}

	return nil
}

// FetchScores fetches scores for a user
func (r *Repository) FetchScores(ctx context.Context, userID string) ([]SubmittedScore, error) {
	scores := make([]SubmittedScore, 0)
	scoresIter := r.client.Collection("users").Doc(userID).Collection("scores").Documents(ctx)
	for {
		scoreDoc, err := scoresIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var score SubmittedScore
		scoreDoc.DataTo(&score)
		scores = append(scores, score)
	}

	return scores, nil
}

// CreateTest creates a test
func (r *Repository) CreateTest(ctx context.Context, id string, request Test) error {
	_, err := r.client.Collection("tests").Doc(id).Set(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

// FetchTest fetches a test
func (r *Repository) FetchTest(ctx context.Context, request FetchTestRequest) (*Test, error) {
	testDoc, err := r.client.Collection("tests").Doc(request.TestID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var test Test
	testDoc.DataTo(&test)
	return &test, nil
}

// FetchAllQuestions fetches all questions
func (r *Repository) FetchAllQuestions(ctx context.Context) ([]Question, error) {
	questions := make([]Question, 0)

	questionsDocIter := r.client.Collection("questions").Documents(ctx)
	for {
		questionDoc, err := questionsDocIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		questionMap := questionDoc.Data()
		questionText := questionMap["text"].(string)

		choices := make([]Choice, 0)
		choicesDocIter := r.client.Collection("questions").Doc(questionDoc.Ref.ID).Collection("choices").Documents(ctx)
		for {
			choiceDoc, err := choicesDocIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}

			choiceMap := choiceDoc.Data()
			choice := Choice{choiceDoc.Ref.ID, choiceMap["text"].(string)}
			choices = append(choices, choice)
		}

		question := Question{questionDoc.Ref.ID, questionText, choices}
		questions = append(questions, question)
	}

	return questions, nil
}

// FetchQuestions fetches questions and choices for a test
func (r *Repository) FetchQuestions(ctx context.Context, test *Test) ([]QuestionWithCorrectChoice, error) {
	questions := make([]QuestionWithCorrectChoice, 0)

	for _, question := range test.SubmittedQuestions {
		questionDoc, err := r.client.Collection("questions").Doc(question.QuestionID).Get(ctx)
		if err != nil {
			return nil, err
		}

		questionMap := questionDoc.Data()
		questionText := questionMap["text"].(string)

		// Fetch choices available for question
		choices := make([]Choice, 0)
		choicesIter := r.client.Collection("questions").Doc(question.QuestionID).Collection("choices").Documents(ctx)
		for {
			doc, err := choicesIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}

			choiceMap := doc.Data()
			choice := Choice{doc.Ref.ID, choiceMap["text"].(string)}
			choices = append(choices, choice)
		}

		questionWithCorrectChoice := QuestionWithCorrectChoice{questionText, choices, question.SelectedChoiceID}
		questions = append(questions, questionWithCorrectChoice)
	}

	return questions, nil
}
