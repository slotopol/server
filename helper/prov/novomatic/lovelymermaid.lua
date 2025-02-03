local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	6-4, --  1 mermaid    2000
	5-5, --  2 lobster    400
	5-5, --  3 turtle     400
	5-4, --  4 blowfish   300
	5-4, --  5 seahorse   200
	5-4, --  6 parrotfish 200
	5, --  7 ace        100
	5, --  8 king       100
	5, --  9 queen      80
	5, -- 10 jack       80
	5, -- 11 ten        80
	5, -- 12 nine       80
	1, -- 13 scatter    80
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 3, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 3,}, --  1 mermaid
	{ 1, 3, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  2 lobster
	{ 1, 1, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  3 turtle
	{ 1, 1, 1, 3, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  4 blowfish
	{ 1, 1, 1, 1, 3, 1, 0, 0, 0, 0, 0, 0, 0,}, --  5 seahorse
	{ 1, 1, 1, 1, 1, 3, 0, 0, 0, 0, 0, 0, 0,}, --  6 parrotfish
	{ 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0,}, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0,}, -- 11 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0,}, -- 12 nine
	{ 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3,}, -- 13 scatter
}

math.randomseed(os.time())
local reel, iter = makereel(symset, neighbours)
for i = 1, 4 do
	table.insert(reel, i, 1)
end
for i = 1, 5 do
	table.insert(reel, i, 2)
end
for i = 1, 5 do
	table.insert(reel, i, 3)
end
for i = 1, 4 do
	table.insert(reel, i, 4)
end
for i = 1, 4 do
	table.insert(reel, i, 5)
end
for i = 1, 4 do
	table.insert(reel, i, 6)
end
printreel(reel, iter)
