local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset1 = {
	1, --  1 wild    (2, 3, 4, 5 reels only)
	1, --  2 scatter
	3, --  3 aryan   1000
	3, --  4 nymph   200
	3, --  5 pot     150
	3, --  6 vase    150
	4, --  7 ace     100
	4, --  8 king    100
	4, --  9 queen   100
	4, -- 10 jack    100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0,}, --  3 aryan
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0,}, --  4 nymph
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0,}, --  5 pot
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0,}, --  6 vase
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 10 jack
}

local symset2 = {
	4, --  1 wild    (2, 3, 4, 5 reels only)
	1, --  2 scatter
	4, --  3 aryan
	4, --  4 nymph
	4, --  5 pot
	4, --  6 vase
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	4, --  3 aryan
	4, --  4 nymph
	4, --  5 pot
	4, --  6 vase
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		local reel1, iter1 = makereel(symset1, neighbours)
		local reel2, iter2 = makereelhot(symset2, 3, {[2]=true}, chunklen, true)
		return reelglue(reel1, reel2), iter1, iter2
	end
	if n == 1 then
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

print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4, 5"
printreel(reelgen(2))
