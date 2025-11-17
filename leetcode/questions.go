package leetcode

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)


const _QUESTIONS_QUERY_FRAGMENT = `fragment questionFields on QuestionNode {
	questionId
	questionFrontendId
	title
	difficulty
}`

type question struct {
	QuestionId string `json:"questionId"`
	QuestionFrontendId string `json:"questionFrontendId"`
	Title string `json:"title"`
	Difficulty string `json:"difficulty"`
}

func (w *Worker) getQuestions(titleSlugs []string) (*map[string]question, error) {
	if len(titleSlugs) == 0 {
		slog.Error("getQuestions() needs at least one question")
		return nil, fmt.Errorf("No questions to fetch")
	}

	var variables = make(map[string]any, len(titleSlugs))
	var queryBuilder strings.Builder
	queryBuilder.WriteString(_QUESTIONS_QUERY_FRAGMENT)
	queryBuilder.WriteString("\nquery getQuestions(")
	for i, titleSlug := range titleSlugs {
		variables[fmt.Sprintf("titleSlug%d", i + 1)] = titleSlug

		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(fmt.Sprintf("$titleSlug%d: String!", i + 1))
	}
	queryBuilder.WriteString(") {")
	for i := range titleSlugs {
		queryBuilder.WriteString(fmt.Sprintf("\nquestion%d: question(titleSlug: $titleSlug%d) {...questionFields}", i + 1, i + 1))
	}
	queryBuilder.WriteString("\n}")

	var rawJson, graphqlErr = w.sendGraphqlQuery(queryBuilder.String(), variables)
	if graphqlErr != nil {
		return nil, graphqlErr
	}

	var questions = make(map[string]question, len(titleSlugs))
	if unmarshalErr := json.Unmarshal(*rawJson, &questions); unmarshalErr != nil {
		slog.Error("Failed to unmarshal JSON response", "json", string(*rawJson))
		return nil, unmarshalErr
	}
	for i, titleSlug := range titleSlugs {
		var key = fmt.Sprintf("question%d", i + 1)
		if questionData, success := questions[key]; success {
			delete(questions, key)
			questions[titleSlug] = questionData
		}
	}
	if len(questions) < len(titleSlugs) {
		slog.Error("Missing Leetcode questions in GraphQL response", slog.Group("request", slog.Int("questions", len(titleSlugs))), slog.Group("response", slog.Int("questions", len(questions))))
		return &questions, fmt.Errorf("Missing Leetcode questions in GraphQL response")
	}
	return &questions, nil
}
