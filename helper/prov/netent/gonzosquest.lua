local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, -- 1 wild     (2, 3, 4 reel)
	2, -- 2 freefall (1, 2, 3 reel)
	3, -- 3 mask1    2500
	3, -- 4 mask2    1000
	4, -- 5 mask3    500
	4, -- 6 mask4    200
	4, -- 7 mask5    100
	5, -- 8 mask6    75
	5, -- 9 mask7    50
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0,}, -- 1 wild
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0,}, -- 2 freefall
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0,}, -- 3 mask1
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0,}, -- 4 mask2
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0,}, -- 5 mask3
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 6 mask4
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 7 mask5
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 8 mask6
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 9 mask7
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))
