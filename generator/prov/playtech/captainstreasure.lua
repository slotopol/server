local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	1, --  3 sabers 5000
	2, --  4 map    2500
	3, --  5 anchor 1000
	5, --  6 ace    500
	5, --  7 king   300
	6, --  8 queen  200
	6, --  9 jack   200
	6, -- 10 ten    100
	6, -- 11 nine   100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 11
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
