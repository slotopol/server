local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset15reg = {
	4, --  1 wild
	1, --  2 scatter (always 1)
	3, --  3 yacht
	3, --  4 bike
	3, --  5 lady
	5, --  6 salesgirl
	5, --  7 courier
	6, --  8 ace
	6, --  9 king
	6, -- 10 queen
	6, -- 11 jack
	6, -- 12 ten
	6, -- 13 nine
}

--[[local symset15bon = {
	2, --  1 wild
	2, --  2 scatter (always has)
	2, --  3 yacht
	2, --  4 bike
	2, --  5 lady
	2, --  6 salesgirl
	2, --  7 courier
	2, --  8 ace
	2, --  9 king
	2, -- 10 queen
	2, -- 11 jack
	4, -- 12 ten
	4, -- 13 nine
}]]

local symset234reg = {
	2, --  1 wild
	0, --  2 scatter (always 0)
	2, --  3 yacht
	2, --  4 bike
	2, --  5 lady
	3, --  6 salesgirl
	3, --  7 courier
	3, --  8 ace
	3, --  9 king
	3, -- 10 queen
	3, -- 11 jack
	3, -- 12 ten
	3, -- 13 nine
}

--[[local symset234bon = {
	2, --  1 wild
	0, --  2 scatter (always 0)
	2, --  3 yacht
	2, --  4 bike
	2, --  5 lady
	2, --  6 salesgirl
	2, --  7 courier
	2, --  8 ace
	2, --  9 king
	2, -- 10 queen
	2, -- 11 jack
	3, -- 12 ten
	3, -- 13 nine
}]]

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 yacht
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4 bike
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  5 lady
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  6 salesgirl
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  7 courier
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13 nine
}

math.randomseed(os.time())
printreel(makereel(symset15reg, neighbours))
printreel(makereel(symset234reg, neighbours))
