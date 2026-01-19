local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 wild     10000
	1, --  2 scatter
	3, --  3 pirate   1000
	3, --  4 lady     1000
	3, --  5 spyglass 1000
	3, --  6 island   500
	3, --  7 parrot   400
	3, --  8 monkey   400
	3, --  9 ace      200
	3, -- 10 king     200
	4, -- 11 queen    150
	4, -- 12 jack     150
	4, -- 13 ten      150
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 2, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0,}, --  3 pirate
	{ 1, 1, 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0,}, --  4 lady
	{ 1, 1, 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0,}, --  5 spyglass
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0,}, --  6 island
	{ 1, 1, 1, 1, 1, 0, 2, 0, 0, 0, 0, 0, 0,}, --  7 parrot
	{ 1, 1, 1, 1, 1, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8 monkey
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13 ten
}

local function reelgen()
	return makereel(symset, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
