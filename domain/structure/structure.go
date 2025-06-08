package structure

// Match Info
type MatchData struct {
	Data MatchInfo `json:"data"`
}

type MatchInfo struct {
	MetaData    MatchMetadata `json:"metadata"`
	PlayersInfo AllPlayer     `json:"players"`
	Team        Team          `json:"teams"`
}

type MatchMetadata struct {
	Map    string `json:"map"`
	ModeId string `json:"mode_id"`
}

type AllPlayer struct {
	Players []Player `json:"all_players"`
}

type Player struct {
	PUUID          string     `json:"puuid"`
	Character      string     `json:"character"`
	Team           string     `json:"team"`
	PlayerStats    PlayerStat `json:"stats"`
	DamageMade     int        `json:"damage_made"`
	DamageReceived int        `json:"damage_received"`
}

type PlayerStat struct {
	Score     int `json:"score"`
	Kills     int `json:"kills"`
	Deaths    int `json:"deaths"`
	Assists   int `json:"assists"`
	BodyShots int `json:"bodyshots"`
	HeadShots int `json:"headshots"`
	LegShot   int `json:"legshots"`
}
type Team struct {
	Red  TeamStat `json:"red"`
	Blue TeamStat `json:"blue"`
}

type TeamStat struct {
	RoundsWon  int `json:"rounds_won"`
	RoundsLost int `json:"rounds_lost"`
}

// MMR Structs
type MrrData struct {
	MrrStats []MrrStat `json:"data"`
}

type MrrStat struct {
	CurrentTier      int       `json:"currenttier"`
	CurrentTierPatch string    `json:"currenttierpatched"`
	Image            RankImage `json:"iamges"`
	MatchId          string    `json:"match_id"`
	MmmrChange       int       `json:"mmr_change_to_last_game"`
	Date             string    `json:"date"`
	DateRaw          int       `json:"date_raw"`
}

type RankImage struct {
	Small string `json:"large"`
	Large string `json:"small"`
}

// Stats Save
type RankStatGameSave struct {
	PUUID          string             `json:"puuid"`
	MatchId        string             `json:"match_id"`
	RawDateInt     int                `json:"raw_date_int"`
	DateStr        string             `json:"date"`
	MmrChange      int                `json:"mmr_change_to_last_game"`
	Map            string             `json:"map"`
	Character      string             `json:"character"`
	PlayerMetaStat PlayerMetaStatSave `json:"stats"`
	RoundsWon      int                `json:"rounds_won"`
	RoundsLost     int                `json:"rounds_lost"`
}

type PlayerMetaStatSave struct {
	Score          int `json:"score"`
	Kills          int `json:"kills"`
	Deaths         int `json:"deaths"`
	Assists        int `json:"assists"`
	BodyShots      int `json:"bodyshots"`
	HeadShots      int `json:"headshots"`
	LegShots       int `json:"legshots"`
	DamageMade     int `json:"damage_made"`
	DamageRecieved int `json:"damage_received"`
}

type ValorantRankDynamoDbRecord struct {
	KeyPuuidMatch  string `dynamodbav:"puuid_match"`
	RawDateInt     int    `dynamodbav:"raw_date_int"`
	PUUID          string `dynamodbav:"puuid"`
	MatchId        string `dynamodbav:"match_id"`
	DateStr        string `dynamodbav:"date"`
	MmrChange      int    `dynamodbav:"mmr_change_to_last_game"`
	Map            string `dynamodbav:"map"`
	Character      string `dynamodbav:"character"`
	RoundsWon      int    `dynamodbav:"rounds_won"`
	RoundsLost     int    `dynamodbav:"rounds_lost"`
	Score          int    `dynamodbav:"score"`
	Kills          int    `dynamodbav:"kills"`
	Deaths         int    `dynamodbav:"deaths"`
	Assists        int    `dynamodbav:"assists"`
	BodyShots      int    `dynamodbav:"bodyshots"`
	HeadShots      int    `dynamodbav:"headshots"`
	LegShots       int    `dynamodbav:"legshots"`
	DamageMade     int    `dynamodbav:"damage_made"`
	DamageRecieved int    `dynamodbav:"damagae_received"`
}
