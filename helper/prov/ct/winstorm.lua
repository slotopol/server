local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	4, --  1 wild      (2, 3, 4 reel)
	1, --  2 scatter
	4, --  3 seven     1000
	4, --  4 coin      100+
	4, --  5 horseshoe 100+
	4, --  6 bell      100+
	4, --  7 ace       100
	4, --  8 king      100
	4, --  9 queen     100
	4, -- 10 jack      100
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	3, --  3 seven
	6, --  4 coin
	6, --  5 horseshoe
	6, --  6 bell
	6, --  7 ace
	6, --  8 king
	6, --  9 queen
	6, -- 10 jack
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
