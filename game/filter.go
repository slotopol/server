package game

import (
	"regexp"
	"strconv"

	"github.com/slotopol/server/util"
)

type Filter func(*GameInfo) bool

var FiltMap = map[string]Filter{
	"all":  func(gi *GameInfo) bool { return true },
	"slot": func(gi *GameInfo) bool { return gi.GT == GTslot },

	"keno":       func(gi *GameInfo) bool { return gi.GT == GTkeno },
	"agt":        func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "agt" },
	"aristocrat": func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "aristocrat" },
	"ct":         func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "ct" || util.ToID(gi.Prov) == "ctinteractive" },
	"betsoft":    func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "betsoft" },
	"igt":        func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "igt" },
	"megajack":   func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "megajack" },
	"netent":     func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "netent" },
	"novomatic":  func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "novomatic" },
	"playngo":    func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "playngo" },
	"playtech":   func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "playtech" },
	"slotopol":   func(gi *GameInfo) bool { return util.ToID(gi.Prov) == "slotopol" },

	"lines":  func(gi *GameInfo) bool { return gi.LN > 0 },
	"ways":   func(gi *GameInfo) bool { return gi.WN > 0 },
	"casc":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPcasc > 0 },
	"jack":   func(gi *GameInfo) bool { return gi.GP&GPjack > 0 },
	"fg":     func(gi *GameInfo) bool { return gi.GP&(GPfghas+GPretrig) > 0 },
	"bon":    func(gi *GameInfo) bool { return gi.BN > 0 },
	"nodate": func(gi *GameInfo) bool { return gi.Date == 0 },
}
var (
	reReel = regexp.MustCompile(`^(\d)x$`)
	reScrn = regexp.MustCompile(`^(\d)x(\d)$`)
	reYEQ  = regexp.MustCompile(`^y(=|eq)(\d{2}|\d{4})$`)
	reYLT  = regexp.MustCompile(`^y(<|lt)(\d{2}|\d{4})$`)
	reYGT  = regexp.MustCompile(`^y(>|gt)(\d{2}|\d{4})$`)
	reLNEQ = regexp.MustCompile(`^ln(=|eq)(\d{1,3})$`)
	reLNLT = regexp.MustCompile(`^ln(<|lt)(\d{1,3})$`)
	reLNGT = regexp.MustCompile(`^ln(>|gt)(\d{1,3})$`)
	reWNEQ = regexp.MustCompile(`^wn(=|eq)(\d{1,4})$`)
	reWNLT = regexp.MustCompile(`^wn(<|lt)(\d{1,4})$`)
	reWNGT = regexp.MustCompile(`^wn(>|gt)(\d{1,4})$`)
)

func GetFilter(key string) Filter {
	key = util.ToLower(key)
	if f, ok := FiltMap[key]; ok {
		return f
	}
	if s := reReel.FindStringSubmatch(key); len(s) > 0 {
		var x, _ = strconv.Atoi(s[1])
		return func(gi *GameInfo) bool { return gi.SX == x }
	}
	if s := reScrn.FindStringSubmatch(key); len(s) > 0 {
		var x, _ = strconv.Atoi(s[1])
		var y, _ = strconv.Atoi(s[2])
		return func(gi *GameInfo) bool { return gi.SX == x && gi.SY == y }
	}
	if s := reYEQ.FindStringSubmatch(key); len(s) > 0 {
		var year, _ = strconv.Atoi(s[2])
		if year < 100 {
			year += 2000
		}
		return func(gi *GameInfo) bool { return gi.Date.Year() == year }
	}
	if s := reYLT.FindStringSubmatch(key); len(s) > 0 {
		var year, _ = strconv.Atoi(s[2])
		if year < 100 {
			year += 2000
		}
		return func(gi *GameInfo) bool { return gi.Date != 0 && gi.Date.Year() < year }
	}
	if s := reYGT.FindStringSubmatch(key); len(s) > 0 {
		var year, _ = strconv.Atoi(s[2])
		if year < 100 {
			year += 2000
		}
		return func(gi *GameInfo) bool { return gi.Date != 0 && gi.Date.Year() > year }
	}
	if s := reLNEQ.FindStringSubmatch(key); len(s) > 0 {
		var ln, _ = strconv.Atoi(s[2])
		return func(gi *GameInfo) bool { return gi.LN == ln }
	}
	if s := reLNLT.FindStringSubmatch(key); len(s) > 0 {
		var ln, _ = strconv.Atoi(s[2])
		return func(gi *GameInfo) bool { return gi.LN != 0 && gi.LN < ln }
	}
	if s := reLNGT.FindStringSubmatch(key); len(s) > 0 {
		var ln, _ = strconv.Atoi(s[2])
		return func(gi *GameInfo) bool { return gi.LN != 0 && gi.LN > ln }
	}
	if s := reWNEQ.FindStringSubmatch(key); len(s) > 0 {
		var wn, _ = strconv.Atoi(s[2])
		return func(gi *GameInfo) bool { return gi.WN == wn }
	}
	if s := reWNLT.FindStringSubmatch(key); len(s) > 0 {
		var wn, _ = strconv.Atoi(s[2])
		return func(gi *GameInfo) bool { return gi.WN != 0 && gi.WN < wn }
	}
	if s := reWNGT.FindStringSubmatch(key); len(s) > 0 {
		var wn, _ = strconv.Atoi(s[2])
		return func(gi *GameInfo) bool { return gi.WN != 0 && gi.WN > wn }
	}
	return nil
}

func Passes(gi *GameInfo, finclist, fexclist []Filter) bool {
	var is bool
	for _, f := range finclist {
		if f(gi) {
			is = true
			break
		}
	}
	if !is {
		return false
	}
	for _, f := range fexclist {
		if f(gi) {
			return false
		}
	}
	return true
}
