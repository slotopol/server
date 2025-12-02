local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	1, --  1 wild      15000
	1, --  2 scatter
	3, --  3 mask      750
	3, --  4 fins      750
	3, --  5 compass   400
	3, --  6 clownfish 250
	3, --  7 butterfly 250
	3, --  8 ace       125
	3, --  9 king      125
	3, -- 10 queen     100
	3, -- 11 jack      100
	3, -- 12 ten       100
	3, -- 13 nine      100+
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 2, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  3 mask
	{ 1, 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  4 fins
	{ 1, 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  5 compass
	{ 1, 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  6 clownfish
	{ 1, 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  7 butterfly
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13 nine
}

math.randomseed(os.time())

local function reelgen()
	return makereel(symset, neighbours)
end

if autoscan then
	return reelgen
end

printreel(reelgen())
