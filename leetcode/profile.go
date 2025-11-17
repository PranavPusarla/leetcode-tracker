package leetcode

import (
	"encoding/json"
	"log/slog"
)

const _USER_PROFILE_QUERY = `query getUserProfile($username: String!) {
	allQuestionsCount {
		difficulty
		count
	}
	matchedUser(username: $username) {
		submitStats {
			acSubmissionNum {
				difficulty
				count
				submissions
			}
			totalSubmissionNum {
				difficulty
				count
				submissions
			}
		}
	}
	recentAcSubmissionList(username: $username, limit: 20) {
		id
		title
		titleSlug
		timestamp
		statusDisplay
		lang
	}
}`


type userProfile struct {
	AllQuestionsCount []questionCount `json:"allQuestionsCount"`
	MatchedUser matchedUser `json:"matchedUser"`
	RecentAcSubmissionList []recentAcSubmission `json:"recentAcSubmissionList"`
}

type questionCount struct {
	Difficulty string `json:"difficulty"`
	Count int `json:"count"`
	Submissions *int `json:"submissions"`
}

type matchedUser struct {
	SubmitStats submitStats `json:"submitStats"`
}

type submitStats struct {
	AcSubmissionNum []questionCount `json:"acSubmissionNum"`
	TotalSubmissionNum []questionCount `json:"totalSubmissionNum"`
}

type recentAcSubmission struct {
	Id string `json:"id"`
	Title string `json:"title"`
	TitleSlug string `json:"titleSlug"`
	Timestamp string `json:"timestamp"`
	StatusDisplay string `json:"statusDisplay"`
	Lang string `json:"lang"`
}


func (w *Worker) getUserProfile(username string) (*userProfile, error) {
	var variables = map[string]any{
		"username": username,
	}
	var rawJson, graphqlErr = w.sendGraphqlQuery(_USER_PROFILE_QUERY, variables)
	if graphqlErr != nil {
		return nil, graphqlErr
	}

	var profileData userProfile
	if unmarshalErr := json.Unmarshal(*rawJson, &profileData); unmarshalErr != nil {
		slog.Error("Failed to unmarshal JSON response", "json", string(*rawJson))
		return nil, unmarshalErr
	}
	return &profileData, nil
}
