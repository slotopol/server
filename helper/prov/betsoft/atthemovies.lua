local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 oscar
	1, --  2 popcorn
	2, --  3 poster
	4, --  4 a
	5, --  5 dummy
	6, --  6 maw
	7, --  7 starship
	7, --  8 heart
	0, --  9 masks
	1, -- 10 projector
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, }, --  1
	{ 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, }, --  2
	{ 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, }, --  3
	{ 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 1, }, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, }, -- 10
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
