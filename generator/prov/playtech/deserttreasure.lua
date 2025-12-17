local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild
	1, --  2 scatter
	3, --  3 shield
	3, --  4 swords
	4, --  5 lamp
	3, --  6 ligature1
	4, --  7 ligature2
	4, --  8 ligature3
	4, --  9 ligature4
	4, -- 10 ligature5
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, }, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, }, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, }, -- 10
}

local function reelgen()
	return makereel(symset, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
