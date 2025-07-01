local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, --  1 wild    10000
	1, --  2 scatter
	2, --  3 red     625
	2, --  4 blonde  625
	3, --  5 beer    625
	3, --  6 ham     625
	5, --  7 ace     200
	5, --  8 king    200
	6, --  9 queen   100
	7, -- 10 jack    100
	7, -- 11 ten     100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0,}, --  1 wild
	{ 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0,}, --  3 red
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0,}, --  4 blonde
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0,}, --  5 beer
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0,}, --  6 ham
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 11 ten
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
