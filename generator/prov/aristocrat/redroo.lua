local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild
	1, --  2 scatter
	3, --  3 redroo 250
	3, --  4 girl   200
	3, --  5 boy    200
	5, --  6 dog    150
	5, --  7 parrot 150
	5, --  8 ace    140
	5, --  9 king   140
	5, -- 10 queen  120
	5, -- 11 jack   120
	5, -- 12 ten    100
	5, -- 13 nine   100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 3, 3, 3, 3, 1, 1, 1, 0, 0, 0, 0, 0, 2,}, --  1
	{ 3, 3, 3, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2
	{ 3, 3, 3, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  3
	{ 3, 3, 3, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  4
	{ 1, 1, 1, 1, 3, 1, 1, 0, 0, 0, 0, 0, 0,}, --  5
	{ 1, 1, 1, 1, 1, 3, 1, 0, 0, 0, 0, 0, 0,}, --  6
	{ 1, 1, 1, 1, 1, 1, 3, 0, 0, 0, 0, 0, 0,}, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0,}, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0,}, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0,}, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0,}, -- 11
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0,}, -- 12
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3,}, -- 13
}

local function reelgen(n)
	local function make()
		return makereel(symset, neighbours)
	end
	if n == 1 or n == 5 then
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
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
