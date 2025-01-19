local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	4, --  1 seven
	4, --  2 bell
	4, --  3 melon
	5, --  4 plum
	5, --  5 orange
	5, --  6 lemon
	5, --  7 cherry
	1, --  8 bonus
	1, --  9 joker
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0 }, --  1
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0 }, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0 }, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0 }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0 }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0 }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0 }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 2 }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 2, 2 }, --  9
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
