local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	0, --  1 wild    (insert directly)
	0, --  2 scatter (2, 3, 4 reel)
	3, --  3 prince
	3, --  4 princess
	4, --  5 castle
	5, --  6 ruby
	6, --  7 shoes
	6, --  8 carpet
	7, --  9 ace
	7, -- 10 king
	7, -- 11 queen
	7, -- 12 jack
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1
	{ 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  3
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  4
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  5
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  6
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, -- 11
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, -- 12
}

math.randomseed(os.time())
local reel, iter = makereel(symset, neighbours)
for i = 1, 4 do
	table.insert(reel, i, 1)
end
printreel(reel, iter)
