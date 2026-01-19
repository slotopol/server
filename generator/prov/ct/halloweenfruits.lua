local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset1 = {
	5, --  1 wild    (2, 3, 4, 5 reels only)
	5, --  2 scatter (2, 3, 4 reels only)
	2, --  3 witch   300
	2, --  4 cat     100
	2, --  5 banana  100
	2, --  6 grape   100
	2, --  7 apple   50
	2, --  8 melon   50
	2, --  9 orange  30
	2, -- 10 lemon   30
	2, -- 11 plum    30
	2, -- 12 cherry  30
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  1 wild
	{ 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  3 witch
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  4 cat
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  5 banana
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  6 grape
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  7 apple
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  8 melon
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  9 orange
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 10 lemon
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 11 plum
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 12 cherry
}

local symset2 = {
	4, --  1 wild
	0, --  2 scatter
	6, --  3 witch
	6, --  4 cat
	4, --  5 banana
	4, --  6 grape
	4, --  7 apple
	4, --  8 melon
	4, --  9 orange
	4, -- 10 lemon
	4, -- 11 plum
	4, -- 12 cherry
}

local chunklen = {
	6, --  1 wild
	1, --  2 scatter
	6, --  3 witch
	6, --  4 cat
	6, --  5 banana
	6, --  6 grape
	6, --  7 apple
	6, --  8 melon
	6, --  9 orange
	6, -- 10 lemon
	6, -- 11 plum
	6, -- 12 cherry
}

local function reelgen(n)
	local ss1, ss2 = tcopy(symset1), tcopy(symset2)
	if n == 1 then
		ss1[1], ss1[2], ss2[1], ss2[2] = 0, 0, 0, 0
	elseif n == 5 then
		ss1[2], ss2[2] = 0, 0
	end
	local reel1, iter1 = makereel(ss1, neighbours)
	local reel2, iter2 = makereelhot(ss2, 3, {[2]=true}, chunklen)
	return reelglue(reel1, reel2), iter1, iter2
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
print "reel 5"
printreel(reelgen(5))
