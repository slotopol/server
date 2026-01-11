local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild
	1, --  2 scatter
	2, --  3 catcher 5000
	3, --  4 man     2500
	3, --  5 woman   1000
	5, --  6 guy     250
	6, --  7 bull    150
	6, --  8 hatchet 150
	7, --  9 ace     80
	7, -- 10 king    80
	7, -- 11 queen   70
	8, -- 12 jack    60
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 2, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, }, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, }, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 1, 1, 0, 0, 0, 0, }, --  3 catcher
	{ 1, 1, 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, }, --  4 man
	{ 1, 1, 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, }, --  5 woman
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, }, --  6 guy
	{ 1, 1, 1, 1, 1, 0, 2, 0, 0, 0, 0, 0, }, --  7 bull
	{ 1, 1, 1, 1, 1, 0, 0, 2, 0, 0, 0, 0, }, --  8 hatchet
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  9 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, -- 10 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, }, -- 11 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, }, -- 12 jack
}

local function reelgen()
	return makereel(symset, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
