local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	0, --  1 moon wolf  1000 (insert directly)
	3, --  2 grey wolf  400
	3, --  3 white wolf 400
	3, --  4 idol1      250
	3, --  5 idol2      250
	5, --  6 ace        150
	5, --  7 king       150
	5, --  8 queen      100
	5, --  9 jack       100
	5, -- 10 ten        100
	5, -- 11 nine       100
	3, -- 12 bonus      (reel 2, 3, 4)
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 3, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 3,}, --  1 moon wolf
	{ 1, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  2 grey wolf
	{ 1, 1, 3, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  3 white wolf
	{ 1, 1, 1, 3, 1, 0, 0, 0, 0, 0, 0, 0,}, --  4 idol1
	{ 1, 1, 1, 1, 3, 0, 0, 0, 0, 0, 0, 0,}, --  5 idol2
	{ 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0,}, --  6 ace
	{ 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0,}, --  7 king
	{ 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0,}, --  8 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0,}, --  9 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0,}, -- 10 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0,}, -- 11 nine
	{ 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3,}, -- 12 bonus
}

math.randomseed(os.time())
local reel, iter = makereel(symset, neighbours)
for i = 1, 4 do
	table.insert(reel, i, 1)
end
printreel(reel, iter)
