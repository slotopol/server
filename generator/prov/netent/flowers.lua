local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 wild
	1, --  2 scatter
	1, --  3 scatter2 (2, 3, 4 reels only)
	5, --  4 red
	2, --  5 red2
	5, --  6 yellow
	2, --  7 yellow2
	5, --  8 green
	2, --  9 green2
	5, -- 10 pink
	2, -- 11 pink2
	5, -- 12 blue
	2, -- 13 blue2
	12, -- 14 ace
	12, -- 15 king
	12, -- 16 queen
	12, -- 17 jack
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,14,15,16,17,
	{ 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1 wild
	{ 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2 scatter
	{ 2, 2, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 0, 0, }, --  3 scatter2
	{ 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  4 red
	{ 0, 0, 2, 2, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 0, 0, }, --  5 red2
	{ 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  6 yellow
	{ 0, 0, 2, 0, 2, 2, 2, 0, 2, 0, 2, 0, 2, 0, 0, 0, 0, }, --  7 yellow2
	{ 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, }, --  8 green
	{ 0, 0, 2, 0, 2, 0, 2, 2, 2, 0, 2, 0, 2, 0, 0, 0, 0, }, --  9 green2
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0, }, -- 10 pink
	{ 0, 0, 2, 0, 2, 0, 2, 0, 2, 2, 2, 0, 2, 0, 0, 0, 0, }, -- 11 pink2
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, }, -- 12 blue
	{ 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 2, 2, 0, 0, 0, 0, }, -- 13 blue2
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, -- 14 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, -- 15 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, }, -- 16 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, }, -- 17 jack
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[3] = 0
	end
	return makereel(ss, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 4, 3"
printreel(reelgen(2))
