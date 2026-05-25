local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	2, --  3 dj       10000
	4, --  4 cocktail 1000
	4, --  5 cup      500
	4, --  6 bull     500
	5, --  7 ace      250
	5, --  8 king     250
	5, --  9 queen    100
	5, -- 10 jack     100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 2, 0, 0, 0, 0, 0, 0, 0,}, --  1 wild (2, 3, 4 reels only)
	{ 2, 2, 2, 0, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 2, 2, 2, 1, 1, 1, 0, 0, 0, 0,}, --  3 dj
	{ 0, 0, 1, 2, 0, 0, 0, 0, 0, 0,}, --  4 cocktail
	{ 0, 0, 1, 0, 2, 0, 0, 0, 0, 0,}, --  5 cup
	{ 0, 0, 1, 0, 0, 2, 0, 0, 0, 0,}, --  6 bull
	{ 0, 0, 0, 0, 0, 0, 1, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 1, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,}, -- 10 jack
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	return makereel(ss, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
