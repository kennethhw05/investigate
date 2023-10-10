package colossus

type SportCode string

const (
	GenericEsports SportCode = "ESPORTS"
)

type PoolType string

const (
	PoolTypeHeadToHead              PoolType = "H2H"
	PoolTypeOverUnder               PoolType = "OVER_UNDER"
	ColossusToStarCurrencyUnitValue          = 100
	StarMinUnitPerLine                       = 1.0
	StarMaxUnitPerLine                       = 1.0
	StarCurrencyCode                         = "STR"
	MinStarPerTicket                         = 1.0
	MaxStarPerTicket                         = 500000000.0
	PoolPrizeDefaultGuarantee                = 0.0
	PoolPrizeDefaultCarryIn                  = 0.0
	PoolPrizeDefaultAllocation               = 0.0
	PoolPrizeDefaultPrizeTypeCode            = "WIN"
)

type PoolStatus string

const (
	PoolStatusCreated   PoolStatus = "CREATED"
	PoolStatusOpen      PoolStatus = "OPEN"
	PoolStatusInPlay    PoolStatus = "IN_PLAY"
	PoolStatusOfficial  PoolStatus = "OFFICIAL"
	PoolStatusAbandoned PoolStatus = "ABANDONED"
)

type SportEventStatus string

const (
	SportEventStatusNotStarted SportEventStatus = "NOT_STARTED"
	SportEventStatusInPlay     SportEventStatus = "IN_PLAY"
	SportEventStatusCompleted  SportEventStatus = "COMPLETED"
	SportEventStatusOfficial   SportEventStatus = "OFFICIAL"
	SportEventStatusAbandoned  SportEventStatus = "ABANDONED"
)

const (
	ConsolationPrizeN1 string = "CON_N1"
	ConsolationPrizeN2 string = "CON_N2"
	ConsolationPrizeN3 string = "CON_N3"
)
