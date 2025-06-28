package structure

type ValorantClipDynamoDbRecord struct {
	ID         string `dynamodbav:"uuid"`
	BaseURL    string `dynamodbav:"base_url"`
	FileName   string `dynamodbav:"file_name"`
	Extenstion string `dynamodbav:"extenstion"`
	FilePath   string `dynamodbav:"file_path"`
	FullUrl    string `dynamodbav:"full_url"`
	MatchId    string `dynamodbav:"match_id"`
}

type ValorantClipJSON struct {
	ID         string `json:"uuid"`
	BaseURL    string `json:"base_url"`
	FileName   string `json:"file_name"`
	Extenstion string `json:"extenstion"`
	FilePath   string `json:"file_path"`
	FullUrl    string `json:"full_url"`
	MatchId    string `json:"match_id"`
}

// id: 23b99bfc-b91e-438a-8a73-1b45bc21f96c
// base url: https://www.beanballer.com/clips
// file name: 23b99bfc-b91e-438a-8a73-1b45bc21f96c.mp4
// extension: mp4
// file path: /clips/23b99bfc-b91e-438a-8a73-1b45bc21f96c.mp4
// full url: https://www.beanballer.com/clips/23b99bfc-b91e-438a-8a73-1b45bc21f96c.mp4
