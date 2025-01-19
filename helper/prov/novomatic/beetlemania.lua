local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, --  1 bee
	3, --  2 snail
	3, --  3 fly
	3, --  4 worm
	4, --  5 ace
	4, --  6 king
	5, --  7 queen
	5, --  8 jack
	5, --  9 ten
	1, -- 10 note
	0, -- 11 jazzbee -- only on 3 reel
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, --  1 bee
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 snail
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 fly
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  4 worm
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  5 ace
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  6 king
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7 queen
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  9 ten
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, -- 10 note
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, -- 11 jazzbee
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
