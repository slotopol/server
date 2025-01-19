local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 wild
	2, --  2 emerald
	2, --  3 heliodor
	3, --  4 ruby
	3, --  5 amethyst
	3, --  6 sapphire
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
	4, -- 11 ten
	4, -- 12 nine
	1, -- 13 scatter
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 2,}, --  1 wild
	{ 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  2 emerald
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  3 heliodor
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  4 ruby
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, 0,}, --  5 amethyst
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0,}, --  6 sapphire
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12 nine
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13 scatter
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
