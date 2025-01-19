local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 scatter
	2, --  2 man
	2, --  3 mind
	2, --  4 internet
	2, --  5 eye
	4, --  6 ace
	4, --  7 king
	4, --  8 queen
	4, --  9 jack
	4, -- 10 ten
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1 scatter
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2 man
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, }, --  3 mind
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, }, --  4 internet
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  5 eye
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  6 ace
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  7 king
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, --  8 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, }, --  9 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, }, -- 10 ten
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
