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
	"keno": func(gi *GameInfo) bool { return gi.GT == GTkeno },

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
	"bon":    func(gi *GameInfo) bool { return gi.BN > 0 },
	"lpay":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPlpay == GPlpay },
	"rpay":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPrpay == GPrpay },
	"apay":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPapay == GPapay },
	"cpay":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPcpay == GPcpay },
	"jack":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPjack > 0 },
	"fill":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPfill > 0 },
	"bm":     func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPbmode > 0 },
	"casc":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPcasc > 0 },
	"cm":     func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPcmult > 0 },
	"fg":     func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&(GPfghas+GPretrig) > 0 },
	"fgr":    func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPfgreel > 0 },
	"fgm":    func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPfgmult > 0 },
	"scat":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPscat > 0 },
	"wild":   func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPwild > 0 },
	"rw":     func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPrwild > 0 },
	"bw":     func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPbwild > 0 },
	"wt":     func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPwturn > 0 },
	"wm":     func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPwmult > 0 },
	"big":    func(gi *GameInfo) bool { return gi.GT == GTslot && gi.GP&GPbsym > 0 },
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
	if len(key) > 0 && key[0] == '~' {
		if f := GetFilter(key[1:]); f != nil {
			return func(gi *GameInfo) bool {
				return !f(gi)
			}
		}
		return nil
	}

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

func Passes(gi *GameInfo, finclist, fexclist [][]Filter) bool {
	var is bool
	for _, sum := range finclist {
		if len(sum) == 0 {
			continue
		}
		is = true
		for _, f := range sum {
			is = is && f(gi)
		}
		if is {
			break
		}
	}
	if !is {
		return false
	}
	for _, sum := range fexclist {
		if len(sum) == 0 {
			continue
		}
		is = true
		for _, f := range sum {
			is = is && f(gi)
		}
		if is {
			return false
		}
	}
	return true
}
