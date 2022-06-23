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
	StartTime uint64
	EndTime   uint64
	Queue     uint32
	Start     uint32
	Count     uint8
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
	One2AssistStreakCount                     int     `json:"12AssistStreakCount"`
	AbilityUses                               int     `json:"abilityUses"`
	AcesBefore15Minutes                       int     `json:"acesBefore15Minutes"`
	AlliedJungleMonsterKills                  float64 `json:"alliedJungleMonsterKills"`
	BaronTakedowns                            int     `json:"baronTakedowns"`
	BlastConeOppositeOpponentCount            int     `json:"blastConeOppositeOpponentCount"`
	BountyGold                                int     `json:"bountyGold"`
	BuffsStolen                               int     `json:"buffsStolen"`
	CompleteSupportQuestInTime                int     `json:"completeSupportQuestInTime"`
	ControlWardTimeCoverageInRiverOrEnemyHalf float64 `json:"controlWardTimeCoverageInRiverOrEnemyHalf"`
	ControlWardsPlaced                        int     `json:"controlWardsPlaced"`
	DamagePerMinute                           float64 `json:"damagePerMinute"`
	DamageTakenOnTeamPercentage               float64 `json:"damageTakenOnTeamPercentage"`
	DancedWithRiftHerald                      int     `json:"dancedWithRiftHerald"`
	DeathsByEnemyChamps                       int     `json:"deathsByEnemyChamps"`
	DodgeSkillShotsSmallWindow                int     `json:"dodgeSkillShotsSmallWindow"`
	DoubleAces                                int     `json:"doubleAces"`
	DragonTakedowns                           int     `json:"dragonTakedowns"`
	EarliestBaron                             float64 `json:"earliestBaron"`
	EarliestDragonTakedown                    float64 `json:"earliestDragonTakedown"`
	EarlyLaningPhaseGoldExpAdvantage          int     `json:"earlyLaningPhaseGoldExpAdvantage"`
	EffectiveHealAndShielding                 int     `json:"effectiveHealAndShielding"`
	ElderDragonKillsWithOpposingSoul          int     `json:"elderDragonKillsWithOpposingSoul"`
	ElderDragonMultikills                     int     `json:"elderDragonMultikills"`
	EnemyChampionImmobilizations              int     `json:"enemyChampionImmobilizations"`
	EnemyJungleMonsterKills                   float64 `json:"enemyJungleMonsterKills"`
	EpicMonsterKillsNearEnemyJungler          int     `json:"epicMonsterKillsNearEnemyJungler"`
	EpicMonsterKillsWithin30SecondsOfSpawn    int     `json:"epicMonsterKillsWithin30SecondsOfSpawn"`
	EpicMonsterSteals                         int     `json:"epicMonsterSteals"`
	EpicMonsterStolenWithoutSmite             int     `json:"epicMonsterStolenWithoutSmite"`
	FasterSupportQuestCompletion              int     `json:"fasterSupportQuestCompletion"`
	FirstTurretKilledTime                     float64 `json:"firstTurretKilledTime"`
	FlawlessAces                              int     `json:"flawlessAces"`
	FullTeamTakedown                          int     `json:"fullTeamTakedown"`
	GameLength                                float64 `json:"gameLength"`
	GetTakedownsInAllLanesEarlyJungleAsLaner  int     `json:"getTakedownsInAllLanesEarlyJungleAsLaner"`
	GoldPerMinute                             float64 `json:"goldPerMinute"`
	HadAfkTeammate                            int     `json:"hadAfkTeammate"`
	HadOpenNexus                              int     `json:"hadOpenNexus"`
	HighestCrowdControlScore                  int     `json:"highestCrowdControlScore"`
	HighestWardKills                          int     `json:"highestWardKills"`
	ImmobilizeAndKillWithAlly                 int     `json:"immobilizeAndKillWithAlly"`
	InitialBuffCount                          int     `json:"initialBuffCount"`
	InitialCrabCount                          int     `json:"initialCrabCount"`
	JungleCsBefore10Minutes                   float64 `json:"jungleCsBefore10Minutes"`
	JunglerKillsEarlyJungle                   int     `json:"junglerKillsEarlyJungle"`
	JunglerTakedownsNearDamagedEpicMonster    int     `json:"junglerTakedownsNearDamagedEpicMonster"`
	KTurretsDestroyedBeforePlatesFall         int     `json:"kTurretsDestroyedBeforePlatesFall"`
	Kda                                       float64 `json:"kda"`
	KillAfterHiddenWithAlly                   int     `json:"killAfterHiddenWithAlly"`
	KillParticipation                         float64 `json:"killParticipation"`
	KilledChampTookFullTeamDamageSurvived     int     `json:"killedChampTookFullTeamDamageSurvived"`
	KillingSprees                             int     `json:"killingSprees"`
	KillsNearEnemyTurret                      int     `json:"killsNearEnemyTurret"`
	KillsOnOtherLanesEarlyJungleAsLaner       int     `json:"killsOnOtherLanesEarlyJungleAsLaner"`
	KillsOnRecentlyHealedByAramPack           int     `json:"killsOnRecentlyHealedByAramPack"`
	KillsUnderOwnTurret                       int     `json:"killsUnderOwnTurret"`
	KillsWithHelpFromEpicMonster              int     `json:"killsWithHelpFromEpicMonster"`
	KnockEnemyIntoTeamAndKill                 int     `json:"knockEnemyIntoTeamAndKill"`
	LandSkillShotsEarlyGame                   int     `json:"landSkillShotsEarlyGame"`
	LaneMinionsFirst10Minutes                 int     `json:"laneMinionsFirst10Minutes"`
	LaningPhaseGoldExpAdvantage               int     `json:"laningPhaseGoldExpAdvantage"`
	LegendaryCount                            int     `json:"legendaryCount"`
	LostAnInhibitor                           int     `json:"lostAnInhibitor"`
	MaxCsAdvantageOnLaneOpponent              float64 `json:"maxCsAdvantageOnLaneOpponent"`
	MaxKillDeficit                            int     `json:"maxKillDeficit"`
	MaxLevelLeadLaneOpponent                  int     `json:"maxLevelLeadLaneOpponent"`
	MoreEnemyJungleThanOpponent               float64 `json:"moreEnemyJungleThanOpponent"`
	MultiKillOneSpell                         int     `json:"multiKillOneSpell"`
	MultiTurretRiftHeraldCount                int     `json:"multiTurretRiftHeraldCount"`
	Multikills                                int     `json:"multikills"`
	MultikillsAfterAggressiveFlash            int     `json:"multikillsAfterAggressiveFlash"`
	MythicItemUsed                            int     `json:"mythicItemUsed"`
	OuterTurretExecutesBefore10Minutes        int     `json:"outerTurretExecutesBefore10Minutes"`
	OutnumberedKills                          int     `json:"outnumberedKills"`
	OutnumberedNexusKill                      int     `json:"outnumberedNexusKill"`
	PerfectDragonSoulsTaken                   int     `json:"perfectDragonSoulsTaken"`
	PerfectGame                               int     `json:"perfectGame"`
	PickKillWithAlly                          int     `json:"pickKillWithAlly"`
	PoroExplosions                            int     `json:"poroExplosions"`
	QuickCleanse                              int     `json:"quickCleanse"`
	QuickFirstTurret                          int     `json:"quickFirstTurret"`
	QuickSoloKills                            int     `json:"quickSoloKills"`
	RiftHeraldTakedowns                       int     `json:"riftHeraldTakedowns"`
	SaveAllyFromDeath                         int     `json:"saveAllyFromDeath"`
	ScuttleCrabKills                          int     `json:"scuttleCrabKills"`
	SkillshotsDodged                          int     `json:"skillshotsDodged"`
	SkillshotsHit                             int     `json:"skillshotsHit"`
	SnowballsHit                              int     `json:"snowballsHit"`
	SoloBaronKills                            int     `json:"soloBaronKills"`
	SoloKills                                 int     `json:"soloKills"`
	SoloTurretsLategame                       int     `json:"soloTurretsLategame"`
	StealthWardsPlaced                        int     `json:"stealthWardsPlaced"`
	SurvivedSingleDigitHpCount                int     `json:"survivedSingleDigitHpCount"`
	SurvivedThreeImmobilizesInFight           int     `json:"survivedThreeImmobilizesInFight"`
	TakedownOnFirstTurret                     int     `json:"takedownOnFirstTurret"`
	Takedowns                                 int     `json:"takedowns"`
	TakedownsAfterGainingLevelAdvantage       int     `json:"takedownsAfterGainingLevelAdvantage"`
	TakedownsBeforeJungleMinionSpawn          int     `json:"takedownsBeforeJungleMinionSpawn"`
	TakedownsFirstXMinutes                    int     `json:"takedownsFirstXMinutes"`
	TakedownsInAlcove                         int     `json:"takedownsInAlcove"`
	TakedownsInEnemyFountain                  int     `json:"takedownsInEnemyFountain"`
	TeamBaronKills                            int     `json:"teamBaronKills"`
	TeamDamagePercentage                      float64 `json:"teamDamagePercentage"`
	TeamElderDragonKills                      int     `json:"teamElderDragonKills"`
	TeamRiftHeraldKills                       int     `json:"teamRiftHeraldKills"`
	ThreeWardsOneSweeperCount                 int     `json:"threeWardsOneSweeperCount"`
	TookLargeDamageSurvived                   int     `json:"tookLargeDamageSurvived"`
	TurretPlatesTaken                         int     `json:"turretPlatesTaken"`
	TurretTakedowns                           int     `json:"turretTakedowns"`
	TurretsTakenWithRiftHerald                int     `json:"turretsTakenWithRiftHerald"`
	TwentyMinionsIn3SecondsCount              int     `json:"twentyMinionsIn3SecondsCount"`
	UnseenRecalls                             int     `json:"unseenRecalls"`
	VisionScoreAdvantageLaneOpponent          float64 `json:"visionScoreAdvantageLaneOpponent"`
	VisionScorePerMinute                      float64 `json:"visionScorePerMinute"`
	WardTakedowns                             int     `json:"wardTakedowns"`
	WardTakedownsBefore20M                    int     `json:"wardTakedownsBefore20M"`
	WardsGuarded                              int     `json:"wardsGuarded"`
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
			q.Set("startTime", strconv.FormatUint(options.StartTime, 10))
		}
		if options.EndTime != 0 {
			q.Set("endTime", strconv.FormatUint(options.EndTime, 10))
		}
		if options.Queue != 0 {
			q.Set("queue", strconv.FormatUint(uint64(options.Queue), 10))
		}
		if options.Start != 0 {
			q.Set("start", strconv.FormatUint(uint64(options.Start), 10))
		}
		if options.Count != 0 {
			q.Set("count", strconv.FormatUint(uint64(options.Count), 10))
		}
		u.RawQuery = q.Encode()
	}
	resp, err := c.request(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	defer resp.Body.Close()

	var m Match
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
