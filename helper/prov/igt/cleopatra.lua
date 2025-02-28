local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 cleopatra 10000
	3, --  2 scarab    750
	3, --  3 broom     750
	2, --  4 shield    400
	3, --  5 stick     250
	3, --  6 eye       250
	3, --  7 ace       125
	3, --  8 king      100
	3, --  9 queen     100
	3, -- 10 jack      100
	3, -- 11 ten       100
	3, -- 12 nine      100
	1, -- 13 sphinx
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 2,}, --  1 cleopatra
	{ 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  2 scarab
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  3 broom
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  4 shield
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, 0,}, --  5 stick
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0,}, --  6 eye
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12 nine
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13 sphinx
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
