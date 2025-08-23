local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	4, --  1 wild       (only on reel 2, 3, 4)
	2, --  2 scatter
	4, --  3 dragon     1000
	5, --  4 jaguar     300
	5, --  5 green mask 200
	5, --  6 blue mask  200
	6, --  7 ace        100
	6, --  8 king       100
	6, --  9 queen      100
	6, -- 10 jack       100
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	3, --  3 dragon
	3, --  4 jaguar
	3, --  5 green mask
	3, --  6 blue mask
	3, --  7 ace
	3, --  8 king
	3, --  9 queen
	3, -- 10 jack
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
