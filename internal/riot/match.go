package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type MatchIdsOptions struct {
	Type      string
	StartTime int64
	EndTime   int64
	Queue     int32
	Start     int32
	Count     int8
}

type Match struct {
	Metadata *Metadata `json:"metadata"`
	Info     *Info     `json:"info"`
}

type Metadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type Challenges struct {
	AssistStreakCount                         float64 `json:"12AssistStreakCount"`
	AbilityUses                               float64 `json:"abilityUses"`
	AcesBefore15Minutes                       float64 `json:"acesBefore15Minutes"`
	AlliedJungleMonsterKills                  float64 `json:"alliedJungleMonsterKills"`
	BaronTakedowns                            float64 `json:"baronTakedowns"`
	BlastConeOppositeOpponentCount            float64 `json:"blastConeOppositeOpponentCount"`
	BountyGold                                float64 `json:"bountyGold"`
	BuffsStolen                               float64 `json:"buffsStolen"`
	CompleteSupportQuestInTime                float64 `json:"completeSupportQuestInTime"`
	ControlWardTimeCoverageInRiverOrEnemyHalf float64 `json:"controlWardTimeCoverageInRiverOrEnemyHalf"`
	ControlWardsPlaced                        float64 `json:"controlWardsPlaced"`
	DamagePerMinute                           float64 `json:"damagePerMinute"`
	DamageTakenOnTeamPercentage               float64 `json:"damageTakenOnTeamPercentage"`
	DancedWithRiftHerald                      float64 `json:"dancedWithRiftHerald"`
	DeathsByEnemyChamps                       float64 `json:"deathsByEnemyChamps"`
	DodgeSkillShotsSmallWindow                float64 `json:"dodgeSkillShotsSmallWindow"`
	DoubleAces                                float64 `json:"doubleAces"`
	DragonTakedowns                           float64 `json:"dragonTakedowns"`
	EarliestBaron                             float64 `json:"earliestBaron"`
	EarliestDragonTakedown                    float64 `json:"earliestDragonTakedown"`
	EarlyLaningPhaseGoldExpAdvantage          float64 `json:"earlyLaningPhaseGoldExpAdvantage"`
	EffectiveHealAndShielding                 float64 `json:"effectiveHealAndShielding"`
	ElderDragonKillsWithOpposingSoul          float64 `json:"elderDragonKillsWithOpposingSoul"`
	ElderDragonMultikills                     float64 `json:"elderDragonMultikills"`
	EnemyChampionImmobilizations              float64 `json:"enemyChampionImmobilizations"`
	EnemyJungleMonsterKills                   float64 `json:"enemyJungleMonsterKills"`
	EpicMonsterKillsNearEnemyJungler          float64 `json:"epicMonsterKillsNearEnemyJungler"`
	EpicMonsterKillsWithin30SecondsOfSpawn    float64 `json:"epicMonsterKillsWithin30SecondsOfSpawn"`
	EpicMonsterSteals                         float64 `json:"epicMonsterSteals"`
	EpicMonsterStolenWithoutSmite             float64 `json:"epicMonsterStolenWithoutSmite"`
	FasterSupportQuestCompletion              float64 `json:"fasterSupportQuestCompletion"`
	FirstTurretKilledTime                     float64 `json:"firstTurretKilledTime"`
	FlawlessAces                              float64 `json:"flawlessAces"`
	FullTeamTakedown                          float64 `json:"fullTeamTakedown"`
	GameLength                                float64 `json:"gameLength"`
	GetTakedownsInAllLanesEarlyJungleAsLaner  float64 `json:"getTakedownsInAllLanesEarlyJungleAsLaner"`
	GoldPerMinute                             float64 `json:"goldPerMinute"`
	HadAfkTeammate                            float64 `json:"hadAfkTeammate"`
	HadOpenNexus                              float64 `json:"hadOpenNexus"`
	HighestCrowdControlScore                  float64 `json:"highestCrowdControlScore"`
	HighestWardKills                          float64 `json:"highestWardKills"`
	ImmobilizeAndKillWithAlly                 float64 `json:"immobilizeAndKillWithAlly"`
	InitialBuffCount                          float64 `json:"initialBuffCount"`
	InitialCrabCount                          float64 `json:"initialCrabCount"`
	JungleCsBefore10Minutes                   float64 `json:"jungleCsBefore10Minutes"`
	JunglerKillsEarlyJungle                   float64 `json:"junglerKillsEarlyJungle"`
	JunglerTakedownsNearDamagedEpicMonster    float64 `json:"junglerTakedownsNearDamagedEpicMonster"`
	KTurretsDestroyedBeforePlatesFall         float64 `json:"kTurretsDestroyedBeforePlatesFall"`
	Kda                                       float64 `json:"kda"`
	KillAfterHiddenWithAlly                   float64 `json:"killAfterHiddenWithAlly"`
	KillParticipation                         float64 `json:"killParticipation"`
	KilledChampTookFullTeamDamageSurvived     float64 `json:"killedChampTookFullTeamDamageSurvived"`
	KillingSprees                             float64 `json:"killingSprees"`
	KillsNearEnemyTurret                      float64 `json:"killsNearEnemyTurret"`
	KillsOnOtherLanesEarlyJungleAsLaner       float64 `json:"killsOnOtherLanesEarlyJungleAsLaner"`
	KillsOnRecentlyHealedByAramPack           float64 `json:"killsOnRecentlyHealedByAramPack"`
	KillsUnderOwnTurret                       float64 `json:"killsUnderOwnTurret"`
	KillsWithHelpFromEpicMonster              float64 `json:"killsWithHelpFromEpicMonster"`
	KnockEnemyIntoTeamAndKill                 float64 `json:"knockEnemyIntoTeamAndKill"`
	LandSkillShotsEarlyGame                   float64 `json:"landSkillShotsEarlyGame"`
	LaneMinionsFirst10Minutes                 float64 `json:"laneMinionsFirst10Minutes"`
	LaningPhaseGoldExpAdvantage               float64 `json:"laningPhaseGoldExpAdvantage"`
	LegendaryCount                            float64 `json:"legendaryCount"`
	LostAnInhibitor                           float64 `json:"lostAnInhibitor"`
	MaxCsAdvantageOnLaneOpponent              float64 `json:"maxCsAdvantageOnLaneOpponent"`
	MaxKillDeficit                            float64 `json:"maxKillDeficit"`
	MaxLevelLeadLaneOpponent                  float64 `json:"maxLevelLeadLaneOpponent"`
	MoreEnemyJungleThanOpponent               float64 `json:"moreEnemyJungleThanOpponent"`
	MultiKillOneSpell                         float64 `json:"multiKillOneSpell"`
	MultiTurretRiftHeraldCount                float64 `json:"multiTurretRiftHeraldCount"`
	Multikills                                float64 `json:"multikills"`
	MultikillsAfterAggressiveFlash            float64 `json:"multikillsAfterAggressiveFlash"`
	MythicItemUsed                            float64 `json:"mythicItemUsed"`
	OuterTurretExecutesBefore10Minutes        float64 `json:"outerTurretExecutesBefore10Minutes"`
	OutnumberedKills                          float64 `json:"outnumberedKills"`
	OutnumberedNexusKill                      float64 `json:"outnumberedNexusKill"`
	PerfectDragonSoulsTaken                   float64 `json:"perfectDragonSoulsTaken"`
	PerfectGame                               float64 `json:"perfectGame"`
	PickKillWithAlly                          float64 `json:"pickKillWithAlly"`
	PoroExplosions                            float64 `json:"poroExplosions"`
	QuickCleanse                              float64 `json:"quickCleanse"`
	QuickFirstTurret                          float64 `json:"quickFirstTurret"`
	QuickSoloKills                            float64 `json:"quickSoloKills"`
	RiftHeraldTakedowns                       float64 `json:"riftHeraldTakedowns"`
	SaveAllyFromDeath                         float64 `json:"saveAllyFromDeath"`
	ScuttleCrabKills                          float64 `json:"scuttleCrabKills"`
	SkillshotsDodged                          float64 `json:"skillshotsDodged"`
	SkillshotsHit                             float64 `json:"skillshotsHit"`
	SnowballsHit                              float64 `json:"snowballsHit"`
	SoloBaronKills                            float64 `json:"soloBaronKills"`
	SoloKills                                 float64 `json:"soloKills"`
	SoloTurretsLategame                       float64 `json:"soloTurretsLategame"`
	StealthWardsPlaced                        float64 `json:"stealthWardsPlaced"`
	SurvivedSingleDigitHpCount                float64 `json:"survivedSingleDigitHpCount"`
	SurvivedThreeImmobilizesInFight           float64 `json:"survivedThreeImmobilizesInFight"`
	TakedownOnFirstTurret                     float64 `json:"takedownOnFirstTurret"`
	Takedowns                                 float64 `json:"takedowns"`
	TakedownsAfterGainingLevelAdvantage       float64 `json:"takedownsAfterGainingLevelAdvantage"`
	TakedownsBeforeJungleMinionSpawn          float64 `json:"takedownsBeforeJungleMinionSpawn"`
	TakedownsFirstXMinutes                    float64 `json:"takedownsFirstXMinutes"`
	TakedownsInAlcove                         float64 `json:"takedownsInAlcove"`
	TakedownsInEnemyFountain                  float64 `json:"takedownsInEnemyFountain"`
	TeamBaronKills                            float64 `json:"teamBaronKills"`
	TeamDamagePercentage                      float64 `json:"teamDamagePercentage"`
	TeamElderDragonKills                      float64 `json:"teamElderDragonKills"`
	TeamRiftHeraldKills                       float64 `json:"teamRiftHeraldKills"`
	ThreeWardsOneSweeperCount                 float64 `json:"threeWardsOneSweeperCount"`
	TookLargeDamageSurvived                   float64 `json:"tookLargeDamageSurvived"`
	TurretPlatesTaken                         float64 `json:"turretPlatesTaken"`
	TurretTakedowns                           float64 `json:"turretTakedowns"`
	TurretsTakenWithRiftHerald                float64 `json:"turretsTakenWithRiftHerald"`
	TwentyMinionsIn3SecondsCount              float64 `json:"twentyMinionsIn3SecondsCount"`
	UnseenRecalls                             float64 `json:"unseenRecalls"`
	VisionScoreAdvantageLaneOpponent          float64 `json:"visionScoreAdvantageLaneOpponent"`
	VisionScorePerMinute                      float64 `json:"visionScorePerMinute"`
	WardTakedowns                             float64 `json:"wardTakedowns"`
	WardTakedownsBefore20M                    float64 `json:"wardTakedownsBefore20M"`
	WardsGuarded                              float64 `json:"wardsGuarded"`
}

type Perks struct {
	StatPerks *struct {
		Defense int `json:"defense"`
		Flex    int `json:"flex"`
		Offense int `json:"offense"`
	} `json:"statPerks"`
	Styles []struct {
		Description string `json:"description"`
		Style       int    `json:"style"`
		Selections  []*struct {
			Perk int `json:"perk"`
			Var1 int `json:"var1"`
			Var2 int `json:"var2"`
			Var3 int `json:"var3"`
		} `json:"selections"`
	} `json:"styles"`
}

type Participant struct {
	Assists                        int         `json:"assists"`
	BaronKills                     int         `json:"baronKills"`
	BountyLevel                    int         `json:"bountyLevel"`
	Challenges                     *Challenges `json:"challenges,omitempty"`
	ChampExperience                int         `json:"champExperience"`
	ChampLevel                     int         `json:"champLevel"`
	ChampionID                     int         `json:"championId"`
	ChampionName                   string      `json:"championName"`
	ChampionTransform              int         `json:"championTransform"`
	ConsumablesPurchased           int         `json:"consumablesPurchased"`
	DamageDealtToBuildings         int         `json:"damageDealtToBuildings"`
	DamageDealtToObjectives        int         `json:"damageDealtToObjectives"`
	DamageDealtToTurrets           int         `json:"damageDealtToTurrets"`
	DamageSelfMitigated            int         `json:"damageSelfMitigated"`
	Deaths                         int         `json:"deaths"`
	DetectorWardsPlaced            int         `json:"detectorWardsPlaced"`
	DoubleKills                    int         `json:"doubleKills"`
	DragonKills                    int         `json:"dragonKills"`
	EligibleForProgression         bool        `json:"eligibleForProgression"`
	FirstBloodAssist               bool        `json:"firstBloodAssist"`
	FirstBloodKill                 bool        `json:"firstBloodKill"`
	FirstTowerAssist               bool        `json:"firstTowerAssist"`
	FirstTowerKill                 bool        `json:"firstTowerKill"`
	GameEndedInEarlySurrender      bool        `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender           bool        `json:"gameEndedInSurrender"`
	GoldEarned                     int         `json:"goldEarned"`
	GoldSpent                      int         `json:"goldSpent"`
	IndividualPosition             string      `json:"individualPosition"`
	InhibitorKills                 int         `json:"inhibitorKills"`
	InhibitorTakedowns             int         `json:"inhibitorTakedowns"`
	InhibitorsLost                 int         `json:"inhibitorsLost"`
	Item0                          int         `json:"item0"`
	Item1                          int         `json:"item1"`
	Item2                          int         `json:"item2"`
	Item3                          int         `json:"item3"`
	Item4                          int         `json:"item4"`
	Item5                          int         `json:"item5"`
	Item6                          int         `json:"item6"`
	ItemsPurchased                 int         `json:"itemsPurchased"`
	KillingSprees                  int         `json:"killingSprees"`
	Kills                          int         `json:"kills"`
	Lane                           string      `json:"lane"`
	LargestCriticalStrike          int         `json:"largestCriticalStrike"`
	LargestKillingSpree            int         `json:"largestKillingSpree"`
	LargestMultiKill               int         `json:"largestMultiKill"`
	LongestTimeSpentLiving         int         `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int         `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int         `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int         `json:"magicDamageTaken"`
	NeutralMinionsKilled           int         `json:"neutralMinionsKilled"`
	NexusKills                     int         `json:"nexusKills"`
	NexusLost                      int         `json:"nexusLost"`
	NexusTakedowns                 int         `json:"nexusTakedowns"`
	ObjectivesStolen               int         `json:"objectivesStolen"`
	ObjectivesStolenAssists        int         `json:"objectivesStolenAssists"`
	ParticipantID                  int         `json:"participantId"`
	PentaKills                     int         `json:"pentaKills"`
	Perks                          *Perks      `json:"perks"`
	PhysicalDamageDealt            int         `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int         `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken            int         `json:"physicalDamageTaken"`
	ProfileIcon                    int         `json:"profileIcon"`
	Puuid                          string      `json:"puuid"`
	QuadraKills                    int         `json:"quadraKills"`
	RiotIDName                     string      `json:"riotIdName"`
	RiotIDTagline                  string      `json:"riotIdTagline"`
	Role                           string      `json:"role"`
	SightWardsBoughtInGame         int         `json:"sightWardsBoughtInGame"`
	Spell1Casts                    int         `json:"spell1Casts"`
	Spell2Casts                    int         `json:"spell2Casts"`
	Spell3Casts                    int         `json:"spell3Casts"`
	Spell4Casts                    int         `json:"spell4Casts"`
	Summoner1Casts                 int         `json:"summoner1Casts"`
	Summoner1ID                    int         `json:"summoner1Id"`
	Summoner2Casts                 int         `json:"summoner2Casts"`
	Summoner2ID                    int         `json:"summoner2Id"`
	SummonerID                     string      `json:"summonerId"`
	SummonerLevel                  int         `json:"summonerLevel"`
	SummonerName                   string      `json:"summonerName"`
	TeamEarlySurrendered           bool        `json:"teamEarlySurrendered"`
	TeamID                         int         `json:"teamId"`
	TeamPosition                   string      `json:"teamPosition"`
	TimeCCingOthers                int         `json:"timeCCingOthers"`
	TimePlayed                     int         `json:"timePlayed"`
	TotalDamageDealt               int         `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int         `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int         `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int         `json:"totalDamageTaken"`
	TotalHeal                      int         `json:"totalHeal"`
	TotalHealsOnTeammates          int         `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int         `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int         `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int         `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int         `json:"totalUnitsHealed"`
	TripleKills                    int         `json:"tripleKills"`
	TrueDamageDealt                int         `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int         `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int         `json:"trueDamageTaken"`
	TurretKills                    int         `json:"turretKills"`
	TurretTakedowns                int         `json:"turretTakedowns"`
	TurretsLost                    int         `json:"turretsLost"`
	UnrealKills                    int         `json:"unrealKills"`
	VisionScore                    int         `json:"visionScore"`
	VisionWardsBoughtInGame        int         `json:"visionWardsBoughtInGame"`
	WardsKilled                    int         `json:"wardsKilled"`
	WardsPlaced                    int         `json:"wardsPlaced"`
	Win                            bool        `json:"win"`
}

type Objective struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}
type Objectives struct {
	Baron      *Objective `json:"baron"`
	Champion   *Objective `json:"champion"`
	Dragon     *Objective `json:"dragon"`
	Inhibitor  *Objective `json:"inhibitor"`
	RiftHerald *Objective `json:"riftHerald"`
	Tower      *Objective `json:"tower"`
}
type Teams struct {
	Bans []*struct {
		ChampionID int `json:"championId"`
		PickTurn   int `json:"pickTurn"`
	} `json:"bans"`
	Objectives *Objectives `json:"objectives"`
	TeamID     int         `json:"teamId"`
	Win        bool        `json:"win"`
}
type Info struct {
	GameCreation       int64          `json:"gameCreation"`
	GameDuration       int            `json:"gameDuration"`
	GameEndTimestamp   int64          `json:"gameEndTimestamp"`
	GameID             int64          `json:"gameId"`
	GameMode           string         `json:"gameMode"`
	GameName           string         `json:"gameName"`
	GameStartTimestamp int64          `json:"gameStartTimestamp"`
	GameType           string         `json:"gameType"`
	GameVersion        string         `json:"gameVersion"`
	MapID              int            `json:"mapId"`
	Participants       []*Participant `json:"participants"`
	PlatformID         string         `json:"platformId"`
	QueueID            int            `json:"queueId"`
	Teams              []*Teams       `json:"teams"`
	TournamentCode     string         `json:"tournamentCode"`
}

func (c *Client) MatchIds(ctx context.Context, region Region, puuid string, options *MatchIdsOptions) ([]string, error) {
	r, v := region.Continent()
	if !v {
		return nil, ErrRegionUnknown
	}
	u, err := url.Parse(fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids", r, puuid))
	if err != nil {
		return nil, err
	}
	if options != nil {
		q := u.Query()
		if options.Type != "" {
			q.Set("type", options.Type)
		}
		if options.StartTime != 0 {
			q.Set("startTime", strconv.FormatInt(options.StartTime, 10))
		}
		if options.EndTime != 0 {
			q.Set("endTime", strconv.FormatInt(options.EndTime, 10))
		}
		if options.Queue != 0 {
			q.Set("queue", strconv.FormatInt(int64(options.Queue), 10))
		}
		if options.Start != 0 {
			q.Set("start", strconv.FormatInt(int64(options.Start), 10))
		}
		if options.Count != 0 {
			q.Set("count", strconv.FormatInt(int64(options.Count), 10))
		}
		u.RawQuery = q.Encode()
	}
	resp, err := c.request(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var s []string
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}
	return s, nil
}

func (c *Client) Match(ctx context.Context, region Region, id string) (*Match, error) {
	r, v := region.Continent()
	if !v {
		return nil, ErrRegionUnknown
	}
	u := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s", r, id)
	resp, err := c.request(ctx, u)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var m Match
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
