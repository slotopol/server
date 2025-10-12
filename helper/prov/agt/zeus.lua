local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, --  1 scatter
	1, --  2 wild     1000
	1, --  3 nymph    500
	3, --  4 armor    100
	4, --  5 boots    80
	4, --  6 crown    80
	4, --  7 staff    60
	5, --  8 pantheon 40
	5, --  9 shield   40
	6, -- 10 glove    20
}

local chunklen = {
	1, --  1 scatter
	1, --  2 wild
	1, --  3 nymph
	2, --  4 armor
	4, --  5 boots
	4, --  6 crown
	5, --  7 staff
	5, --  8 pantheon
	5, --  9 shield
	5, -- 10 glove
}

math.randomseed(os.time())
printreel(makereelhot(symset, 4, {[1]=true}, chunklen, true))
