package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)


func InsertProblemsAndSolutions(db *sql.DB, problems []Problem, solutions []Solution) error {
	var tx, txErr = db.Begin()
	if txErr != nil {
		slog.Error("Database transaction failed")
		return txErr
	}
	defer tx.Rollback()

	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO problem (id, name, slug, difficulty)\nVALUES ")
	var problemParams = make([]any, 0, len(problems) * 4)
	for i, problem := range problems {
		var suffix string
		if i < len(problems) - 1 {
			suffix = ","
		}
		queryBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)%s\n", i * 4 + 1, i * 4 + 2, i * 4 + 3, i * 4 + 4, suffix))
		problemParams = append(problemParams, problem.Id, problem.Name, problem.Slug, problem.Difficulty)
	}
	queryBuilder.WriteString("ON CONFLICT DO NOTHING;")
	if _, err := tx.Exec(queryBuilder.String(), problemParams...); err != nil {
		slog.Error("Failed to insert problems into database", "query", queryBuilder.String(), "params", problemParams)
		return err
	}

	queryBuilder.Reset()
	queryBuilder.WriteString("INSERT INTO solution (id, problem_id, username, timestamp, language)\nVALUES ")
	var solutionParams = make([]any, 0, len(solutions) * 5)
	for i, solution := range solutions {
		var suffix string
		if i < len(solutions) - 1 {
			suffix = ","
		}
		queryBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)%s\n", i * 5 + 1, i * 5 + 2, i * 5 + 3, i * 5 + 4, i * 5 + 5, suffix))
		solutionParams = append(solutionParams, solution.Id, solution.ProblemId, solution.Username, solution.Timestamp, solution.Language)
	}
	queryBuilder.WriteString("ON CONFLICT DO NOTHING;")
	if _, err := tx.Exec(queryBuilder.String(), solutionParams...); err != nil {
		slog.Error("Failed to insert solutions into database", "query", queryBuilder.String(), "params", solutionParams)
		return err
	}

	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit insert transaction")
		return err
	}
	return nil
}
