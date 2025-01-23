local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset135 = {
	1, -- 1 diamond   5000
	2, -- 2 topaz     500
	9, -- 3 sapphire  200
	2, -- 4 heliodor  200
	9, -- 5 ruby      200
	2, -- 6 tanzanite 200
	9, -- 7 emerald   200
	1, -- 8 star
}

local symset24 = {
	1, -- 1 diamond   5000
	9, -- 2 topaz     500
	2, -- 3 sapphire  200
	9, -- 4 heliodor  200
	2, -- 5 ruby      200
	9, -- 6 tanzanite 200
	2, -- 7 emerald   200
	1, -- 8 star
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8
	{ 2, 0, 0, 0, 0, 0, 0, 0 }, -- 1 diamond
	{ 0, 2, 0, 0, 0, 0, 0, 0 }, -- 2 topaz
	{ 0, 0, 2, 0, 0, 0, 0, 0 }, -- 3 sapphire
	{ 0, 0, 0, 2, 0, 0, 0, 0 }, -- 4 heliodor
	{ 0, 0, 0, 0, 2, 0, 0, 0 }, -- 5 ruby
	{ 0, 0, 0, 0, 0, 2, 0, 0 }, -- 6 tanzanite
	{ 0, 0, 0, 0, 0, 0, 2, 0 }, -- 7 emerald
	{ 0, 0, 0, 0, 0, 0, 0, 2 }, -- 8 star
}

math.randomseed(os.time())
printreel(makereel(symset135, neighbours))
printreel(makereel(symset24, neighbours))
