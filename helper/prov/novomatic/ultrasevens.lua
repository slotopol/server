local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	0, -- 1 seven
	3, -- 2 melon
	3, -- 3 grapes
	8, -- 4 plum
	8, -- 5 orange
	8, -- 6 lemon
	8, -- 7 cherry
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7,
	{ 3, 0, 0, 0, 0, 0, 0, }, -- 1 seven
	{ 0, 3, 0, 0, 0, 0, 0, }, -- 2 melon
	{ 0, 0, 3, 0, 0, 0, 0, }, -- 3 grapes
	{ 0, 0, 0, 2, 0, 0, 0, }, -- 4 plum
	{ 0, 0, 0, 0, 2, 0, 0, }, -- 5 orange
	{ 0, 0, 0, 0, 0, 2, 0, }, -- 6 lemon
	{ 0, 0, 0, 0, 0, 0, 2, }, -- 7 cherry
}

math.randomseed(os.time())
local reel, iter = makereel(symset, neighbours)
for i = 1, 5 do
	table.insert(reel, i, 1)
end
for i = 1, 5 do
	table.insert(reel, i, 2)
end
for i = 1, 5 do
	table.insert(reel, i, 3)
end
for i = 1, 5 do
	table.insert(reel, i, 4)
end
for i = 1, 5 do
	table.insert(reel, i, 5)
end
for i = 1, 5 do
	table.insert(reel, i, 6)
end
for i = 1, 5 do
	table.insert(reel, i, 7)
end
printreel(reel, iter)
