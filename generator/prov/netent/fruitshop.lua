local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 wild (2, 3, 4 reels only)
	2, --  2 cherry 2000
	3, --  3 plum   1000
	3, --  4 lemon  750
	3, --  5 orange 500
	4, --  6 melon  200
	4, --  7 ace    150
	5, --  8 king   100
	6, --  9 queen  75
	6, -- 10 jack   60
	7, -- 11 ten    50
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, }, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, }, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, }, -- 11
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
printreel(reelgen(1))
printreel(reelgen(2))
printreel(reelgen(3))
printreel(reelgen(4))
printreel(reelgen(5))
