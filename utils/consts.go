package utils

import "github.com/astaxie/beego/logs"

const (
	PathRoom   = "rooms"
	PathPlayer = "players"

	PathUsed   = "used"
	PathPool   = "pool"
	PathConfig = "config"
)

const (
	CharWolf     = "wolf"
	CharVillager = "villager"
	CharProphet  = "prophet"
	CharWitch    = "witch"
	CharHunter   = "hunter"
	CharKingWolf = "kingwolf"
	CharGuard    = "guard"
)

const (
	OPEnter     = "enter"
	OPEnterSucc = "enterSucc"
)

const (
	StatusLive = "live"
	StatusDead = "dead"
)

const (
	AlarmSucc = "succ"
	AlarmFail = "failed"
)

const (
	PoolSize   = 128
	TimesRetry = 3
)

var LogLevel map[string]int = map[string]int{"emergency": logs.LevelEmergency,
	"alert":    logs.LevelAlert,
	"critical": logs.LevelCritical,
	"error":    logs.LevelError,
	"warning":  logs.LevelWarning,
	"notice":   logs.LevelNotice,
	"info":     logs.LevelInformational,
	"debug":    logs.LevelDebug}
