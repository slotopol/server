local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 wild (2, 4 reels only)
	1, --  2 scatter
	3, --  3 shoe   1000
	3, --  4 woman  500
	3, --  5 spurs  400
	3, --  6 belt   400
	3, --  7 saddle 120
	3, --  8 hat    120
	3, --  9 boots  120
	4, -- 10 ace    100
	4, -- 11 king   100
	4, -- 12 queen  100
	4, -- 13 jack   100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1 wild (2, 4 reels only)
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2 scatter
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  3 shoe
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  4 woman
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, }, --  5 spurs
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, }, --  6 belt
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, }, --  7 saddle
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  8 hat
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  9 boots
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, -- 10 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, -- 11 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, }, -- 12 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, }, -- 13 jack
}

local function reelgen(n)
	local function make()
		return makereel(symset, neighbours)
	end
	if n == 1 or n == 3 or n == 5 then
		local n1 = symset[1]
		symset[1] = 0
		local reel, iter = make()
		symset[1] = n1
		return reel, iter
	else
		return make()
	end
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 3, 5"
printreel(reelgen(1))
print "reel 2, 4"
printreel(reelgen(2))
