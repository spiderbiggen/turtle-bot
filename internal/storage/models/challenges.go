package models

type ChallengeResult struct {
	ID    ChallengeType
	Value float64
}

type CompareOperator uint8

const (
	Max CompareOperator = iota
	Min
	AboveAverage
	BelowAverage
)

type ChallengeType string

const (
	AssistStreakCount                         ChallengeType = "12AssistStreakCount"
	AbilityUses                               ChallengeType = "abilityUses"
	AcesBefore15Minutes                       ChallengeType = "acesBefore15Minutes"
	AlliedJungleMonsterKills                  ChallengeType = "alliedJungleMonsterKills"
	Assists                                   ChallengeType = "assists"
	BaronTakedowns                            ChallengeType = "baronTakedowns"
	BlastConeOppositeOpponentCount            ChallengeType = "blastConeOppositeOpponentCount"
	BountyGold                                ChallengeType = "bountyGold"
	BountyLevel                               ChallengeType = "bountyLevel"
	BuffsStolen                               ChallengeType = "buffsStolen"
	CompleteSupportQuestInTime                ChallengeType = "completeSupportQuestInTime"
	ConsumablesPurchased                      ChallengeType = "consumablesPurchased"
	ControlWardTimeCoverageInRiverOrEnemyHalf ChallengeType = "controlWardTimeCoverageInRiverOrEnemyHalf"
	ControlWardsPlaced                        ChallengeType = "controlWardsPlaced"
	DamageDealtToBuildings                    ChallengeType = "damageDealtToBuildings"
	DamageDealtToObjectives                   ChallengeType = "damageDealtToObjectives"
	DamageDealtToTurrets                      ChallengeType = "damageDealtToTurrets"
	DamagePerMinute                           ChallengeType = "damagePerMinute"
	DamageTakenOnTeamPercentage               ChallengeType = "damageTakenOnTeamPercentage"
	DamageSelfMitigated                       ChallengeType = "damageSelfMitigated"
	DancedWithRiftHerald                      ChallengeType = "dancedWithRiftHerald"
	Deaths                                    ChallengeType = "deaths"
	DeathsByEnemyChamps                       ChallengeType = "deathsByEnemyChamps"
	DetectorWardsPlaced                       ChallengeType = "detectorWardsPlaced"
	DodgeSkillShotsSmallWindow                ChallengeType = "dodgeSkillShotsSmallWindow"
	DoubleAces                                ChallengeType = "doubleAces"
	DoubleKills                               ChallengeType = "doubleKills"
	DragonTakedowns                           ChallengeType = "dragonTakedowns"
	EarliestBaron                             ChallengeType = "earliestBaron"
	EarliestDragonTakedown                    ChallengeType = "earliestDragonTakedown"
	EarlyLaningPhaseGoldExpAdvantage          ChallengeType = "earlyLaningPhaseGoldExpAdvantage"
	EffectiveHealAndShielding                 ChallengeType = "effectiveHealAndShielding"
	ElderDragonKillsWithOpposingSoul          ChallengeType = "elderDragonKillsWithOpposingSoul"
	ElderDragonMultikills                     ChallengeType = "elderDragonMultikills"
	EnemyChampionImmobilizations              ChallengeType = "enemyChampionImmobilizations"
	EnemyJungleMonsterKills                   ChallengeType = "enemyJungleMonsterKills"
	EpicMonsterKillsNearEnemyJungler          ChallengeType = "epicMonsterKillsNearEnemyJungler"
	EpicMonsterKillsWithin30SecondsOfSpawn    ChallengeType = "epicMonsterKillsWithin30SecondsOfSpawn"
	EpicMonsterSteals                         ChallengeType = "epicMonsterSteals"
	EpicMonsterStolenWithoutSmite             ChallengeType = "epicMonsterStolenWithoutSmite"
	FasterSupportQuestCompletion              ChallengeType = "fasterSupportQuestCompletion"
	FirstTurretKilledTime                     ChallengeType = "firstTurretKilledTime"
	FlawlessAces                              ChallengeType = "flawlessAces"
	FullTeamTakedown                          ChallengeType = "fullTeamTakedown"
	GameLength                                ChallengeType = "gameLength"
	GetTakedownsInAllLanesEarlyJungleAsLaner  ChallengeType = "getTakedownsInAllLanesEarlyJungleAsLaner"
	GoldEarned                                ChallengeType = "goldEarned"
	GoldPerMinute                             ChallengeType = "goldPerMinute"
	GoldSpent                                 ChallengeType = "goldSpent"
	HadAfkTeammate                            ChallengeType = "hadAfkTeammate"
	HadOpenNexus                              ChallengeType = "hadOpenNexus"
	HighestCrowdControlScore                  ChallengeType = "highestCrowdControlScore"
	HighestWardKills                          ChallengeType = "highestWardKills"
	ImmobilizeAndKillWithAlly                 ChallengeType = "immobilizeAndKillWithAlly"
	InitialBuffCount                          ChallengeType = "initialBuffCount"
	InitialCrabCount                          ChallengeType = "initialCrabCount"
	InhibitorKills                            ChallengeType = "inhibitorKills"
	InhibitorTakedowns                        ChallengeType = "inhibitorTakedowns"
	InhibitorsLost                            ChallengeType = "inhibitorsLost"
	ItemsPurchased                            ChallengeType = "itemsPurchased"
	JungleCsBefore10Minutes                   ChallengeType = "jungleCsBefore10Minutes"
	JunglerKillsEarlyJungle                   ChallengeType = "junglerKillsEarlyJungle"
	JunglerTakedownsNearDamagedEpicMonster    ChallengeType = "junglerTakedownsNearDamagedEpicMonster"
	KTurretsDestroyedBeforePlatesFall         ChallengeType = "kTurretsDestroyedBeforePlatesFall"
	Kda                                       ChallengeType = "kda"
	KillAfterHiddenWithAlly                   ChallengeType = "killAfterHiddenWithAlly"
	KillParticipation                         ChallengeType = "killParticipation"
	KilledChampTookFullTeamDamageSurvived     ChallengeType = "killedChampTookFullTeamDamageSurvived"
	KillingSprees                             ChallengeType = "killingSprees"
	Kills                                     ChallengeType = "kills"
	KillsNearEnemyTurret                      ChallengeType = "killsNearEnemyTurret"
	KillsOnOtherLanesEarlyJungleAsLaner       ChallengeType = "killsOnOtherLanesEarlyJungleAsLaner"
	KillsOnRecentlyHealedByAramPack           ChallengeType = "killsOnRecentlyHealedByAramPack"
	KillsUnderOwnTurret                       ChallengeType = "killsUnderOwnTurret"
	KillsWithHelpFromEpicMonster              ChallengeType = "killsWithHelpFromEpicMonster"
	KnockEnemyIntoTeamAndKill                 ChallengeType = "knockEnemyIntoTeamAndKill"
	LandSkillShotsEarlyGame                   ChallengeType = "landSkillShotsEarlyGame"
	LaneMinionsFirst10Minutes                 ChallengeType = "laneMinionsFirst10Minutes"
	LaningPhaseGoldExpAdvantage               ChallengeType = "laningPhaseGoldExpAdvantage"
	LargestCriticalStrike                     ChallengeType = "largestCriticalStrike"
	LargestKillingSpree                       ChallengeType = "largestKillingSpree"
	LargestMultiKill                          ChallengeType = "largestMultiKill"
	LegendaryCount                            ChallengeType = "legendaryCount"
	LongestTimeSpentLiving                    ChallengeType = "longestTimeSpentLiving"
	LostAnInhibitor                           ChallengeType = "lostAnInhibitor"
	MaxCsAdvantageOnLaneOpponent              ChallengeType = "maxCsAdvantageOnLaneOpponent"
	MaxKillDeficit                            ChallengeType = "maxKillDeficit"
	MaxLevelLeadLaneOpponent                  ChallengeType = "maxLevelLeadLaneOpponent"
	MoreEnemyJungleThanOpponent               ChallengeType = "moreEnemyJungleThanOpponent"
	MultiKillOneSpell                         ChallengeType = "multiKillOneSpell"
	MultiTurretRiftHeraldCount                ChallengeType = "multiTurretRiftHeraldCount"
	Multikills                                ChallengeType = "multikills"
	MultikillsAfterAggressiveFlash            ChallengeType = "multikillsAfterAggressiveFlash"
	MythicItemUsed                            ChallengeType = "mythicItemUsed"
	NeutralMinionsKilled                      ChallengeType = "neutralMinionsKilled"
	ObjectivesStolen                          ChallengeType = "objectivesStolen"
	ObjectivesStolenAssists                   ChallengeType = "objectivesStolenAssists"
	OuterTurretExecutesBefore10Minutes        ChallengeType = "outerTurretExecutesBefore10Minutes"
	OutnumberedKills                          ChallengeType = "outnumberedKills"
	OutnumberedNexusKill                      ChallengeType = "outnumberedNexusKill"
	PentaKills                                ChallengeType = "pentaKills"
	PerfectDragonSoulsTaken                   ChallengeType = "perfectDragonSoulsTaken"
	PerfectGame                               ChallengeType = "perfectGame"
	PhysicalDamageDealt                       ChallengeType = "physicalDamageDealt"
	PhysicalDamageDealtToChampions            ChallengeType = "physicalDamageDealtToChampions"
	PhysicalDamageTaken                       ChallengeType = "physicalDamageTaken"
	PickKillWithAlly                          ChallengeType = "pickKillWithAlly"
	PoroExplosions                            ChallengeType = "poroExplosions"
	QuadraKills                               ChallengeType = "quadraKills"
	QuickCleanse                              ChallengeType = "quickCleanse"
	QuickFirstTurret                          ChallengeType = "quickFirstTurret"
	QuickSoloKills                            ChallengeType = "quickSoloKills"
	RiftHeraldTakedowns                       ChallengeType = "riftHeraldTakedowns"
	SaveAllyFromDeath                         ChallengeType = "saveAllyFromDeath"
	ScuttleCrabKills                          ChallengeType = "scuttleCrabKills"
	SkillshotsDodged                          ChallengeType = "skillshotsDodged"
	SkillshotsHit                             ChallengeType = "skillshotsHit"
	SnowballsHit                              ChallengeType = "snowballsHit"
	SoloBaronKills                            ChallengeType = "soloBaronKills"
	SoloKills                                 ChallengeType = "soloKills"
	SoloTurretsLategame                       ChallengeType = "soloTurretsLategame"
	StealthWardsPlaced                        ChallengeType = "stealthWardsPlaced"
	SurvivedSingleDigitHpCount                ChallengeType = "survivedSingleDigitHpCount"
	SurvivedThreeImmobilizesInFight           ChallengeType = "survivedThreeImmobilizesInFight"
	TakedownOnFirstTurret                     ChallengeType = "takedownOnFirstTurret"
	Takedowns                                 ChallengeType = "takedowns"
	TakedownsAfterGainingLevelAdvantage       ChallengeType = "takedownsAfterGainingLevelAdvantage"
	TakedownsBeforeJungleMinionSpawn          ChallengeType = "takedownsBeforeJungleMinionSpawn"
	TakedownsFirstXMinutes                    ChallengeType = "takedownsFirstXMinutes"
	TakedownsInAlcove                         ChallengeType = "takedownsInAlcove"
	TakedownsInEnemyFountain                  ChallengeType = "takedownsInEnemyFountain"
	TeamBaronKills                            ChallengeType = "teamBaronKills"
	TeamDamagePercentage                      ChallengeType = "teamDamagePercentage"
	TeamElderDragonKills                      ChallengeType = "teamElderDragonKills"
	TeamRiftHeraldKills                       ChallengeType = "teamRiftHeraldKills"
	ThreeWardsOneSweeperCount                 ChallengeType = "threeWardsOneSweeperCount"
	TimeCCingOthers                           ChallengeType = "timeCCingOthers"
	TimePlayed                                ChallengeType = "timePlayed"
	TookLargeDamageSurvived                   ChallengeType = "tookLargeDamageSurvived"
	TotalDamageDealt                          ChallengeType = "totalDamageDealt"
	TotalDamageDealtToChampions               ChallengeType = "totalDamageDealtToChampions"
	TotalDamageShieldedOnTeammates            ChallengeType = "totalDamageShieldedOnTeammates"
	TotalDamageTaken                          ChallengeType = "totalDamageTaken"
	TotalHeal                                 ChallengeType = "totalHeal"
	TotalHealsOnTeammates                     ChallengeType = "totalHealsOnTeammates"
	TotalMinionsKilled                        ChallengeType = "totalMinionsKilled"
	TotalTimeCCDealt                          ChallengeType = "totalTimeCCDealt"
	TotalTimeSpentDead                        ChallengeType = "totalTimeSpentDead"
	TotalUnitsHealed                          ChallengeType = "totalUnitsHealed"
	TripleKills                               ChallengeType = "tripleKills"
	TrueDamageDealt                           ChallengeType = "trueDamageDealt"
	TrueDamageDealtToChampions                ChallengeType = "trueDamageDealtToChampions"
	TrueDamageTaken                           ChallengeType = "trueDamageTaken"
	TurretPlatesTaken                         ChallengeType = "turretPlatesTaken"
	TurretTakedowns                           ChallengeType = "turretTakedowns"
	TurretsLost                               ChallengeType = "turretsLost"
	TurretsTakenWithRiftHerald                ChallengeType = "turretsTakenWithRiftHerald"
	TwentyMinionsIn3SecondsCount              ChallengeType = "twentyMinionsIn3SecondsCount"
	UnseenRecalls                             ChallengeType = "unseenRecalls"
	UnrealKills                               ChallengeType = "unrealKills"
	VisionScore                               ChallengeType = "visionScore"
	VisionScoreAdvantageLaneOpponent          ChallengeType = "visionScoreAdvantageLaneOpponent"
	VisionScorePerMinute                      ChallengeType = "visionScorePerMinute"
	VisionWardsBoughtInGame                   ChallengeType = "visionWardsBoughtInGame"
	WardTakedowns                             ChallengeType = "wardTakedowns"
	WardTakedownsBefore20M                    ChallengeType = "wardTakedownsBefore20M"
	WardsPlaced                               ChallengeType = "wardsPlaced"
	WardsGuarded                              ChallengeType = "wardsGuarded"
)

func (r *ChallengeResult) Description() string {
	switch r.ID {
	case AssistStreakCount:
		return "12 Assist streak count"
	case AbilityUses:
		return "Ability uses"
	case AcesBefore15Minutes:
		return "Aces before 15 minutes"
	case AlliedJungleMonsterKills:
		return "Allied jungle monster kills"
	case Assists:
		return "Assists"
	case BaronTakedowns:
		return "Baron takedowns"
	case BlastConeOppositeOpponentCount:
		return "Blast cone opponent count"
	case BountyGold:
		return "Bounty gold"
	case BountyLevel:
		return "Bounty level"
	case BuffsStolen:
		return "Buffs stolen"
	case CompleteSupportQuestInTime:
		return "Complete support quest in time"
	case ConsumablesPurchased:
		return "Consumables purchased"
	case ControlWardTimeCoverageInRiverOrEnemyHalf:
		return "Control ward time coverage in river or enemy half"
	case ControlWardsPlaced:
		return "Control wards placed"
	case DamageDealtToBuildings:
		return "Damage dealt to buildings"
	case DamageDealtToObjectives:
		return "Damage dealt to objectives"
	case DamageDealtToTurrets:
		return "Damage dealt to turrets"
	case DamagePerMinute:
		return "Damage per minute"
	case DamageTakenOnTeamPercentage:
		return "Damage taken on team percentage"
	case DamageSelfMitigated:
		return "Damage self mitigated"
	case DancedWithRiftHerald:
		return "Danced with rift herald"
	case Deaths:
		return "Deaths"
	case DeathsByEnemyChamps:
		return "Deaths by enemy champs"
	case DetectorWardsPlaced:
		return "Detector wards placed"
	case DodgeSkillShotsSmallWindow:
		return "Dodge skill shots small window"
	case DoubleAces:
		return "Double aces"
	case DoubleKills:
		return "Double kills"
	case DragonTakedowns:
		return "Dragon takedowns"
	case EarliestBaron:
		return "Earliest baron"
	case EarliestDragonTakedown:
		return "Earliest dragon takedown"
	case EarlyLaningPhaseGoldExpAdvantage:
		return "Early laning phase gold exp advantage"
	case EffectiveHealAndShielding:
		return "Effective heal and shielding"
	case ElderDragonKillsWithOpposingSoul:
		return "Elder dragon kills with opposing soul"
	case ElderDragonMultikills:
		return "Elder dragon multikills"
	case EnemyChampionImmobilizations:
		return "Enemy champion immobilizations"
	case EnemyJungleMonsterKills:
		return "Enemy jungle monster kills"
	case EpicMonsterKillsNearEnemyJungler:
		return "Epic monster kills near enemy jungler"
	case EpicMonsterKillsWithin30SecondsOfSpawn:
		return "Epic monster kills within 30 seconds of spawn"
	case EpicMonsterSteals:
		return "Epic monster steals"
	case EpicMonsterStolenWithoutSmite:
		return "Epic monster stolen without smite"
	case FasterSupportQuestCompletion:
		return "Faster support quest completion"
	case FirstTurretKilledTime:
		return "First turret killed time"
	case FlawlessAces:
		return "Flawless aces"
	case FullTeamTakedown:
		return "Full team takedown"
	case GameLength:
		return "Game length"
	case GetTakedownsInAllLanesEarlyJungleAsLaner:
		return "Get takedowns in all lanes early jungle as laner"
	case GoldEarned:
		return "Gold earned"
	case GoldPerMinute:
		return "Gold per minute"
	case GoldSpent:
		return "Gold spent"
	case HadAfkTeammate:
		return "Had afk teammate"
	case HadOpenNexus:
		return "Had open nexus"
	case HighestCrowdControlScore:
		return "Highest crowd control score"
	case HighestWardKills:
		return "Highest ward kills"
	case ImmobilizeAndKillWithAlly:
		return "Immobilize and kill with ally"
	case InitialBuffCount:
		return "Initial buff count"
	case InitialCrabCount:
		return "Initial crab count"
	case InhibitorKills:
		return "Inhibitor kills"
	case InhibitorTakedowns:
		return "Inhibitor takedowns"
	case InhibitorsLost:
		return "Inhibitors lost"
	case ItemsPurchased:
		return "Items purchased"
	case JungleCsBefore10Minutes:
		return "Jungle CS before 10 minutes"
	case JunglerKillsEarlyJungle:
		return "Jungler kills early jungle"
	case JunglerTakedownsNearDamagedEpicMonster:
		return "Jungler takedowns near damaged epic monster"
	case KTurretsDestroyedBeforePlatesFall:
		return "Turrets destroyed before plates fall"
	case Kda:
		return "KDA"
	case KillAfterHiddenWithAlly:
		return "Kill after hidden with ally"
	case KillParticipation:
		return "Kill participation"
	case KilledChampTookFullTeamDamageSurvived:
		return "Killed champ took full team damage survived"
	case KillingSprees:
		return "Killing sprees"
	case Kills:
		return "Kills"
	case KillsNearEnemyTurret:
		return "Kills near enemy turret"
	case KillsOnOtherLanesEarlyJungleAsLaner:
		return "Kills on other lanes early jungle as laner"
	case KillsOnRecentlyHealedByAramPack:
		return "Kills on recently healed by aram pack"
	case KillsUnderOwnTurret:
		return "Kills under own turret"
	case KillsWithHelpFromEpicMonster:
		return "Kills with help from epic monster"
	case KnockEnemyIntoTeamAndKill:
		return "Knock enemy into team and kill"
	case LandSkillShotsEarlyGame:
		return "Land skill shots early game"
	case LaneMinionsFirst10Minutes:
		return "Lane minions first 10 minutes"
	case LaningPhaseGoldExpAdvantage:
		return "Laning phase gold exp advantage"
	case LargestCriticalStrike:
		return "Largest critical strike"
	case LargestKillingSpree:
		return "Largest killing spree"
	case LargestMultiKill:
		return "Largest multi kill"
	case LegendaryCount:
		return "Legendary count"
	case LongestTimeSpentLiving:
		return "Longest time spent living"
	case LostAnInhibitor:
		return "Lost an inhibitor"
	case MaxCsAdvantageOnLaneOpponent:
		return "Max CS advantage on lane opponent"
	case MaxKillDeficit:
		return "Max kill deficit"
	case MaxLevelLeadLaneOpponent:
		return "Max level lead lane opponent"
	case MoreEnemyJungleThanOpponent:
		return "More enemy jungle than opponent"
	case MultiKillOneSpell:
		return "Multi kill one spell"
	case MultiTurretRiftHeraldCount:
		return "Multi turret rift herald count"
	case Multikills:
		return "Multikills"
	case MultikillsAfterAggressiveFlash:
		return "Multikills after aggressive flash"
	case MythicItemUsed:
		return "Mythic item used"
	case NeutralMinionsKilled:
		return "Neutral minions killed"
	case ObjectivesStolen:
		return "Objectives stolen"
	case ObjectivesStolenAssists:
		return "Objectives stolen assists"
	case OuterTurretExecutesBefore10Minutes:
		return "Outer turret executes before 10 minutes"
	case OutnumberedKills:
		return "Outnumbered kills"
	case OutnumberedNexusKill:
		return "Outnumbered nexus kills"
	case PentaKills:
		return "Penta kills"
	case PerfectDragonSoulsTaken:
		return "Perfect dragon souls taken"
	case PerfectGame:
		return "Perfect game"
	case PhysicalDamageDealt:
		return "Physical damage dealt"
	case PhysicalDamageDealtToChampions:
		return "Physical damage dealt to champions"
	case PhysicalDamageTaken:
		return "Physical damage taken"
	case PickKillWithAlly:
		return "Pick kill with ally"
	case PoroExplosions:
		return "Poro explosions"
	case QuadraKills:
		return "Quadra kills"
	case QuickCleanse:
		return "Quick cleanse"
	case QuickFirstTurret:
		return "Quick first turret"
	case QuickSoloKills:
		return "Quick solo kills"
	case RiftHeraldTakedowns:
		return "Rift herald takedowns"
	case SaveAllyFromDeath:
		return "Save ally from death"
	case ScuttleCrabKills:
		return "Scuttle crab kills"
	case SkillshotsDodged:
		return "Skillshots dodged"
	case SkillshotsHit:
		return "Skillshots hit"
	case SnowballsHit:
		return "Snowballs hit"
	case SoloBaronKills:
		return "Solo baron kills"
	case SoloKills:
		return "Solo kills"
	case SoloTurretsLategame:
		return "Solo turrets late game"
	case StealthWardsPlaced:
		return "Stealth wards placed"
	case SurvivedSingleDigitHpCount:
		return "Survived single digit HP count"
	case SurvivedThreeImmobilizesInFight:
		return "Survived three immobilizes in fight"
	case TakedownOnFirstTurret:
		return "Takedown on first turret"
	case Takedowns:
		return "Takedowns"
	case TakedownsAfterGainingLevelAdvantage:
		return "Takedowns after gaining level advantage"
	case TakedownsBeforeJungleMinionSpawn:
		return "Takedowns before jungle minion spawn"
	case TakedownsFirstXMinutes:
		return "Takedowns first X minutes"
	case TakedownsInAlcove:
		return "Takedowns in alcove"
	case TakedownsInEnemyFountain:
		return "Takedowns in enemy fountain"
	case TeamBaronKills:
		return "Team baron kills"
	case TeamDamagePercentage:
		return "Team damage percentage"
	case TeamElderDragonKills:
		return "Team elder dragon kills"
	case TeamRiftHeraldKills:
		return "Team rift herald kills"
	case ThreeWardsOneSweeperCount:
		return "Three wards one sweeper count"
	case TimeCCingOthers:
		return "Time ccing others"
	case TimePlayed:
		return "Time played"
	case TookLargeDamageSurvived:
		return "Took large damage survived"
	case TotalDamageDealt:
		return "Total damage dealt"
	case TotalDamageDealtToChampions:
		return "Total damage dealt to champions"
	case TotalDamageShieldedOnTeammates:
		return "Total damage shielded on teammates"
	case TotalDamageTaken:
		return "Total damage taken"
	case TotalHeal:
		return "Total heal"
	case TotalHealsOnTeammates:
		return "Total heals on teammates"
	case TotalMinionsKilled:
		return "Total minions killed"
	case TotalTimeCCDealt:
		return "Total time CC dealt"
	case TotalTimeSpentDead:
		return "Total time spent dead"
	case TotalUnitsHealed:
		return "Total units healed"
	case TripleKills:
		return "Triple kills"
	case TrueDamageDealt:
		return "True damage dealt"
	case TrueDamageDealtToChampions:
		return "True damage dealt to champions"
	case TrueDamageTaken:
		return "True damage taken"
	case TurretPlatesTaken:
		return "Turret plates taken"
	case TurretTakedowns:
		return "Turret takedowns"
	case TurretsLost:
		return "Turrets lost"
	case TurretsTakenWithRiftHerald:
		return "Turrets taken with rift herald"
	case TwentyMinionsIn3SecondsCount:
		return "Twenty minions in 3 seconds count"
	case UnseenRecalls:
		return "Unseen recalls"
	case UnrealKills:
		return "Unreal kills"
	case VisionScore:
		return "Vision score"
	case VisionScoreAdvantageLaneOpponent:
		return "Vision score advantage over lane opponent"
	case VisionScorePerMinute:
		return "Vision score per minute"
	case VisionWardsBoughtInGame:
		return "Vision wards bought in game"
	case WardTakedowns:
		return "Ward takedowns"
	case WardTakedownsBefore20M:
		return "Ward takedowns before 20 minutes"
	case WardsPlaced:
		return "Wards placed"
	case WardsGuarded:
		return "Wards guarded"
	default:
		return ""
	}
}
