local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset1 = {
	3, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	3, --  3 blue
	3, --  4 red
	3, --  5 swords
	3, --  6 axe
	3, --  7 ace
	3, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  3 blue
	{ 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  4 red
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  5 swords
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  6 axe
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 10 jack
}

local symset2 = {
	3, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 blue
	4, --  4 red
	4, --  5 swords
	4, --  6 axe
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local chunklen = {
	3, --  1 wild
	1, --  2 scatter
	6, --  3 blue
	6, --  4 red
	6, --  5 swords
	6, --  6 axe
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
	if n == 1 or n == 5 then
		local n11, n21 = symset1[1], symset2[1]
		symset1[1], symset2[1] = 0, 0
		local reel, iter1, iter2 = make()
		symset1[1], symset2[1] = n11, n21
		return reel, iter1, iter2
	else
		return make()
	end
end

if autoscan then
	return reelgen
end

print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
