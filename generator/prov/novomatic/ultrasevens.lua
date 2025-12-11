local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

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
addsym(reel, 7, 5)
addsym(reel, 6, 5)
addsym(reel, 5, 5)
addsym(reel, 4, 5)
addsym(reel, 3, 5)
addsym(reel, 2, 5)
addsym(reel, 1, 5)
printreel(reel, iter)
