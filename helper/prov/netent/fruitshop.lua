local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

-- 0 wilds on 1, 5 reel
-- 2 wilds on 2, 3 reel
-- 1 wild on 4 reel
local symset = {
	1, --  1 wild
	3, --  2 cherry
	3, --  3 plum
	3, --  4 lemon
	3, --  5 orange
	3, --  6 melon
	5, --  7 ace
	5, --  8 king
	6, --  9 queen
	6, -- 10 jack
	6, -- 11 ten
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, }, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, }, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, }, -- 11
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
