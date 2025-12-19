local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	0, --  1 wild (insert directly)
	0, --  2 scatter
	4, --  3 car
	4, --  4 tv
	4, --  5 recorder
	4, --  6 projector
	4, --  7 boots
	4, --  8 column
	5, --  9 ace
	5, -- 10 king
	5, -- 11 queen
	5, -- 12 jack
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1
	{ 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2
	{ 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  3
	{ 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, }, --  4
	{ 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, }, --  5
	{ 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, }, --  6
	{ 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, }, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, }, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, }, -- 11
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, }, -- 12
}

local function reelgen()
	local reel, iter = makereel(symset, neighbours)
	addsym(reel, 1, 4)
	return reel, iter
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
