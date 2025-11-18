local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset135 = {
	1, -- 1 diamond   5000
	1, -- 2 star
	2, -- 3 topaz     500
	9, -- 4 sapphire  200
	2, -- 5 heliodor  200
	9, -- 6 ruby      200
	2, -- 7 tanzanite 200
	9, -- 8 emerald   200
}

local symset24 = {
	1, -- 1 diamond   5000
	1, -- 2 star
	9, -- 3 topaz     500
	2, -- 4 sapphire  200
	9, -- 5 heliodor  200
	2, -- 6 ruby      200
	9, -- 7 tanzanite 200
	2, -- 8 emerald   200
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8
	{ 2, 2, 0, 0, 0, 0, 0, 0 }, -- 1 diamond
	{ 2, 2, 0, 0, 0, 0, 0, 0 }, -- 2 star
	{ 0, 0, 2, 0, 0, 0, 0, 0 }, -- 3 topaz
	{ 0, 0, 0, 2, 0, 0, 0, 0 }, -- 4 sapphire
	{ 0, 0, 0, 0, 2, 0, 0, 0 }, -- 5 heliodor
	{ 0, 0, 0, 0, 0, 2, 0, 0 }, -- 6 ruby
	{ 0, 0, 0, 0, 0, 0, 2, 0 }, -- 7 tanzanite
	{ 0, 0, 0, 0, 0, 0, 0, 2 }, -- 8 emerald
}

math.randomseed(os.time())
printreel(makereel(symset135, neighbours))
printreel(makereel(symset24, neighbours))
