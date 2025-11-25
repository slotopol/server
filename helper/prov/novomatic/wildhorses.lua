local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	0, --  1 wild        1000 (insert directly)
	3, --  2 scatter     (2, 3, 4 reel)
	4, --  3 white       750
	4, --  4 black       750
	4, --  5 blue amulet 400
	4, --  6 red amulet  400
	7, --  7 ace         300
	7, --  8 king        300
	8, --  9 queen       200
	8, -- 10 jack        200
	8, -- 11 ten         100
	8, -- 12 nine        100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 0, 3, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  1 wild
	{ 3, 3, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  3 white
	{ 1, 1, 1, 3, 1, 1, 0, 0, 0, 0, 0, 0,}, --  4 black
	{ 1, 1, 1, 1, 3, 1, 0, 0, 0, 0, 0, 0,}, --  5 blue amulet
	{ 1, 1, 1, 1, 1, 3, 0, 0, 0, 0, 0, 0,}, --  6 red amulet
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 11 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 12 nine
}

math.randomseed(os.time())
local reel, iter = makereel(symset, neighbours)
for i = 1, 4 do
	table.insert(reel, i, 1)
end
printreel(reel, iter)
