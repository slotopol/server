local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	3, -- 1 star   10000
	5, -- 2 bell   500
	5, -- 3 grapes 500
	7, -- 4 plum   200
	7, -- 5 orange 200
	8, -- 6 lemon  200
	5, -- 7 cherry 200
	2, -- 8 crown
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8
	{ 2, 0, 0, 0, 0, 0, 0, 2 }, -- 1 star
	{ 0, 2, 0, 0, 0, 0, 0, 0 }, -- 2 bell
	{ 0, 0, 2, 0, 0, 0, 0, 0 }, -- 3 grapes
	{ 0, 0, 0, 2, 0, 0, 0, 0 }, -- 4 plum
	{ 0, 0, 0, 0, 2, 0, 0, 0 }, -- 5 orange
	{ 0, 0, 0, 0, 0, 2, 0, 0 }, -- 6 lemon
	{ 0, 0, 0, 0, 0, 0, 2, 0 }, -- 7 cherry
	{ 2, 0, 0, 0, 0, 0, 0, 2 }, -- 8 crown
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
