local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset1 = {
	2, --  1 wild    (2, 3, 4, 5 reels only)
	2, --  2 scatter
	2, --  3 samurai 2000
	2, --  4 geisha  400
	2, --  5 bowl    300
	2, --  6 coins   300
	2, --  7 ace     200
	2, --  8 king    200
	2, --  9 queen   200
	2, -- 10 jack    100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  3 samurai
	{ 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  4 geisha
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  5 bowl
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  6 coins
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 10 jack
}

local symset2 = {
	0, --  1 wild
	0, --  2 scatter
	4, --  3 samurai
	4, --  4 geisha
	4, --  5 bowl
	4, --  6 coins
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	6, --  3 samurai
	6, --  4 geisha
	6, --  5 bowl
	6, --  6 coins
	6, --  7 ace
	6, --  8 king
	6, --  9 queen
	6, -- 10 jack
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		local reel1, iter1 = makereel(symset1, neighbours)
		local reel2, iter2 = makereelhot(symset2, 3, {[2]=true}, chunklen, true)
		return reelglue(reel1, reel2), iter1, iter2
	end
	if n == 1 then
		local n1, n2 = symset1[1], symset2[1]
		symset1[1], symset2[1] = 0, 0
		local reel, iter1, iter2 = make()
		symset1[1], symset2[1] = n1, n2
		return reel, iter1, iter2
	else
		return make()
	end
end

if autoscan then
	return reelgen
end

print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4, 5"
printreel(reelgen(2))
