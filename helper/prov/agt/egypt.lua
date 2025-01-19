local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, --  1 wild
	1, --  2 scatter
	2, --  3 anubis
	2, --  4 sphinx
	3, --  5 snake
	3, --  6 mummy
	3, --  7 wall
	4, --  8 cat
	5, --  9 ace
	7, -- 10 king
	7, -- 11 queen
	7, -- 12 jack
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1 wild
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2 scatter
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  3 anubis
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, }, --  4 sphinx
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, }, --  5 snake
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, }, --  6 mummy
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  7 wall
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  8 cat
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  9 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, -- 10 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, }, -- 11 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, }, -- 12 jack
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
