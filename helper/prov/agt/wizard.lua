local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	6, --  1 wild
	0, --  2 scatter (2, 3, 4 reels only)
	5, --  3 owl
	5, --  4 cat
	5, --  5 cauldron
	5, --  6 emerald
	5, --  7 ruby
	5, --  8 ace
	5, --  9 king
	6, -- 10 queen
	6, -- 11 jack
}

math.randomseed(os.time())
printreel(makereelhot(symset, 4, {}, {}))
