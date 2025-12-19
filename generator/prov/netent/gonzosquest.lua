local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, -- 1 wild     (2, 3, 4 reels only)
	2, -- 2 freefall (1, 2, 3 reels only)
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

local function reelgen(n)
	local n1, n2 = symset[1], symset[2]
	if n == 1 or n == 5 then -- 2, 3, 4 reels only
		symset[1] = 0
	end
	if n == 4 or n == 5 then -- 1, 2, 3 reels only
		symset[2] = 0
	end
	local reel, iter = makereel(symset, neighbours)
	if n == 1 or n == 5 then -- 2, 3, 4 reels only
		symset[1] = n1
	end
	if n == 4 or n == 5 then -- 1, 2, 3 reels only
		symset[2] = n2
	end
	return reel, iter
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen(1))
printreel(reelgen(2))
printreel(reelgen(3))
printreel(reelgen(4))
printreel(reelgen(5))
