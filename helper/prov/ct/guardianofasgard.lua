local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset15 = {
	1, --  1 scatter
	0, --  2 wild    1500 (always 0 here)
	3, --  3 bowl    150
	4, --  4 wolf    40
	4, --  5 helm    40
	4, --  6 axe     40
	5, --  7 ace     10
	5, --  8 king    10
	5, --  9 queen   10
	5, -- 10 jack    10
}

local symset234 = {
	3, --  1 scatter
	0, --  2 wild    1500
	8, --  3 bowl    150
	10, --  4 wolf    40
	10, --  5 helm    40
	10, --  6 axe     40
	12, --  7 ace     10
	12, --  8 king    10
	12, --  9 queen   10
	12, -- 10 jack    10
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1 scatter
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 wild
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 bowl
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4 wolf
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  5 helm
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  6 axe
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 10 jack
}

math.randomseed(os.time())
printreel(makereel(symset15, neighbours))
local reel, iter = makereel(symset234, neighbours)
for i = 1, 3 do
	table.insert(reel, i, 2)
end
printreel(reel, iter)
