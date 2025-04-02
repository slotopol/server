local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	0, --  1 wild         (2, 3, 4 reel, insert directly)
	2, --  2 Oliver       5000
	2, --  3 friends      1000
	2, --  4 couple       1000
	3, --  5 sweet-stuffs 500
	3, --  6 cocktails    500
	3, --  7 flower       200
	4, --  8 lime         200
	5, --  9 olives       100
	5, -- 10 strawberries 100
	5, -- 11 oranges      100
	4, -- 12 cherry       100
	1, -- 13 scatter
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 2,}, --  1 wild
	{ 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  2 Oliver
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  3 friends
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  4 couple
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, 0,}, --  5 sweet-stuffs
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0,}, --  6 cocktails
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  7 flower
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8 lime
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9 olives
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10 strawberries
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11 oranges
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12 cherry
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13 scatter
}

math.randomseed(os.time())
local reel, iter = makereel(symset, neighbours)
printreel(reel, iter)
for i = 1, 3 do
	table.insert(reel, i, 1)
end
printreel(reel, iter)
