local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 wild    (reel 2, 4)
	2, --  2 volcano 2000
	2, --  3 drums   750
	2, --  4 guitar  400
	2, --  5 coconut 400
	4, --  6 ace     150
	5, --  7 king    125
	5, --  8 queen   125
	5, --  9 jack    100
	5, -- 10 ten     100
	1, -- 11 scatter
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 2,}, --  1 wild
	{ 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 volcano
	{ 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  3 drums
	{ 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  4 guitar
	{ 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  5 coconut
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  6 ace
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7 king
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  9 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 10 ten
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 11 scatter
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
