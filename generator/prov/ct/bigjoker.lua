local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset1 = {
	1, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	2, --  3 seven   750
	3, --  4 grape   200
	3, --  5 melon   200
	3, --  6 apple   200
	4, --  7 orange  100
	4, --  8 lemon   100
	4, --  9 plum    100
	4, -- 10 cherry  100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0,}, --  3 seven
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0,}, --  4 grape
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0,}, --  5 melon
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0,}, --  6 apple
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  7 orange
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  8 lemon
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, --  9 plum
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 10 cherry
}

local symset2 = {
	1, --  1 wild
	1, --  2 scatter
	4, --  3 seven
	4, --  4 grape
	4, --  5 melon
	4, --  6 apple
	4, --  7 orange
	4, --  8 lemon
	4, --  9 plum
	4, -- 10 cherry
}

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	6, --  3 seven
	6, --  4 grape
	6, --  5 melon
	6, --  6 apple
	6, --  7 orange
	6, --  8 lemon
	6, --  9 plum
	6, -- 10 cherry
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		local reel1, iter1 = makereel(symset1, neighbours)
		local reel2, iter2 = makereelhot(symset2, 3, {[2]=true}, chunklen)
		return reelglue(reel1, reel2), iter1, iter2
	end
	if n == 1 or n == 5 then
		local n11, n21, n12 = symset1[1], symset2[1], symset1[2]
		symset1[1], symset2[1], symset1[2] = 0, 0, 0
		local reel, iter = make()
		symset1[1], symset2[1], symset1[2] = n11, n21, n12
		return reel, iter
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
