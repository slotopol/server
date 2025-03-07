local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symsetreg = {
	1, --  1 wild       10000
	2, --  2 girl       500
	2, --  3 father     500
	2, --  4 doggy      200
	2, --  5 kitty      200
	3, --  6 watermelon 100
	3, --  7 peach      100
	3, --  8 plum       100
	3, --  9 lemon      100
	3, -- 10 cherry     100
	1, -- 11 scatter
	0, -- 12 diamond    (reel 2, 3, 4)
}

local neighboursreg = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 2, 2,}, --  1 wild
	{ 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  2 girl
	{ 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  3 father
	{ 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, 0,}, --  4 doggy
	{ 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0,}, --  5 kitty
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  6 watermelon
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  7 peach
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  8 plum
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  9 lemon
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 10 cherry
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, -- 11 scatter
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, -- 12 diamond
}

local symsetbon = {
	1, -- 1 diamond
	4, -- 2 emerald
	4, -- 3 sapphire
	4, -- 4 ruby
	4, -- 5 heliodor
}

local neighboursbon = {
	--1, 2, 3, 4, 5,
	{ 2, 0, 0, 0, 0,}, -- 1 diamond
	{ 0, 2, 0, 0, 0,}, -- 2 emerald
	{ 0, 0, 2, 0, 0,}, -- 3 sapphire
	{ 0, 0, 0, 2, 0,}, -- 4 ruby
	{ 0, 0, 0, 0, 2,}, -- 5 heliodor
}

math.randomseed(os.time())
printreel(makereel(symsetreg, neighboursreg))
printreel(makereel(symsetbon, neighboursbon))
