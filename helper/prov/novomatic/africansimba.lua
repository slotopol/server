local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, --  1 wild     (2, 3, 4 reel)
	0, --  2 scatter  (1, 3, 5 reel)
	3, --  3 giraffe  2500
	3, --  4 buffalo  750
	3, --  5 lemur    250
	3, --  6 flamingo 250
	4, --  7 ace      125
	4, --  8 king     125
	4, --  9 queen    125
	4, -- 10 jack     100
	4, -- 11 ten      100
	4, -- 12 nine     100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, }, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, }, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, }, --  3 giraffe
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, }, --  4 buffalo
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, }, --  5 lemur
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, }, --  6 flamingo
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, }, -- 11 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, }, -- 12 nine
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
