local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, --  1 Cleopatra
	2, --  2 sphinx
	2, --  3 mask
	3, --  4 necklace
	3, --  5 beads
	4, --  6 ace
	4, --  7 king
	5, --  8 queen
	5, --  9 jack
	5, -- 10 ten
	1, -- 11 scatter
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 2,}, --  1 Cleopatra
	{ 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 sphinx
	{ 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  3 mask
	{ 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  4 necklace
	{ 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  5 beads
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  6 ace
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7 king
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  9 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 10 ten
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 11 scatter
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
