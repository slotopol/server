local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 seven  1000
	3, --  2 bell   500
	4, --  3 melon  200
	4, --  4 grapes 200
	5, --  5 plum   150
	5, --  6 orange 150
	5, --  7 lemon  100
	5, --  8 cherry 100
	0, --  9 star
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0 }, --  1
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0 }, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0 }, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0 }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0 }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0 }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0 }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0 }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2 }, --  9
}

local function reelgen()
	return makereel(symset, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
