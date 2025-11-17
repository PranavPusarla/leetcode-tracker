package database

import (
	"database/sql"
	"log/slog"
	"time"
)


const scoreboardQuery = `WITH daily_cte AS (
	SELECT
		username,
		timestamp::date,
		count(*) AS problems,
		sum(
			CASE
				WHEN p.difficulty = 'easy' THEN 1
				WHEN p.difficulty = 'medium' THEN 2
				WHEN p.difficulty = 'hard' THEN 3
				ELSE NULL
			END
		) AS daily_score
	FROM problem AS p
	INNER JOIN valid_solution AS s
	ON p.id = s.problem_id
	WHERE timestamp >= $1
	GROUP BY username, timestamp::date
),
total_cte AS (
	SELECT
		username,
		count(*) FILTER (WHERE daily_score >= 2) AS days,
		sum(problems) AS problems,
		sum(daily_score) AS score
	FROM daily_cte
	GROUP BY username
)
SELECT username, days, problems, score, rank() OVER (ORDER BY days DESC, score DESC, problems DESC) AS place
FROM total_cte
ORDER BY place;`

func GetScoreboard(db *sql.DB, startDate time.Time) ([]ScoreboardEntry, error) {
	var rows, queryErr = db.Query(scoreboardQuery, startDate)
	if queryErr != nil {
		slog.Error("Failed to query scoreboard from database")
		return nil, queryErr
	}
	defer rows.Close()

	var scoreboard []ScoreboardEntry
	for rows.Next() {
		var scoreboardEntry ScoreboardEntry
		if err := rows.Scan(&scoreboardEntry.Username, &scoreboardEntry.Days, &scoreboardEntry.Problems, &scoreboardEntry.Score, &scoreboardEntry.Place); err != nil {
			slog.Error("Failed to fetch values from query result")
			return nil, err
		}
		scoreboard = append(scoreboard, scoreboardEntry)
	}
	if err := rows.Err(); err != nil {
		slog.Error("Failed to fetch values from query result")
		return nil, err
	}

	return scoreboard, nil
}
