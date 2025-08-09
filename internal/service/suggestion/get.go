package suggestion

import "smart-interview/internal/database"

type Data struct {
	Suggestion string `json:"suggestion"`
	CreatedAt  string `json:"created_at"`
}

func GetSuggestion(uid int64) (*Data, error) {
	db := database.GetDB()

	var data Data
	var exists bool

	row := db.QueryRow("SELECT suggestion, created_at FROM suggestion WHERE user = ?", uid)

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM suggestion WHERE user = ?)", uid).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	err = row.Scan(&data.Suggestion, &data.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
