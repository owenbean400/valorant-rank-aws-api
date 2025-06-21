package structure

type ValorantClipDynamoDbRecord struct {
	ID         string `dynamodbdav:"id"`
	BaseURL    string `dynamodbdev:"base_url"`
	FileName   string `dynamodbdev:"file_name"`
	Extenstion string `dynamodbdev:"extenstion"`
	FilePath   string `dynamodbdev:"file_path"`
	FullUrl    string `dynamodbdev:"full_url"`
}

type ValorantClipJSON struct {
	ID         string `json:"id"`
	BaseURL    string `json:"base_url"`
	FileName   string `json:"file_name"`
	Extenstion string `json:"extenstion"`
	FilePath   string `json:"file_path"`
	FullUrl    string `json:"full_url"`
}

// id: 23b99bfc-b91e-438a-8a73-1b45bc21f96c
// base url: https://www.beanballer.com/clips
// file name: 23b99bfc-b91e-438a-8a73-1b45bc21f96c.mp4
// extension: mp4
// file path: /clips/23b99bfc-b91e-438a-8a73-1b45bc21f96c.mp4
// full url: https://www.beanballer.com/clips/23b99bfc-b91e-438a-8a73-1b45bc21f96c.mp4
