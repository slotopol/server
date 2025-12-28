local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symsetreg = {
	1, --  1 wild      15000
	1, --  2 scatter
	3, --  3 mask      750
	3, --  4 fins      750
	3, --  5 compass   400
	3, --  6 clownfish 250
	3, --  7 butterfly 250
	4, --  8 ace       125
	4, --  9 king      125
	4, -- 10 queen     100
	4, -- 11 jack      100
	5, -- 12 ten       100
	5, -- 13 nine      100+
}

local symsetbon = {
	2, --  1 wild      15000
	2, --  2 scatter
	2, --  3 mask      750
	3, --  4 fins      750
	3, --  5 compass   400
	3, --  6 clownfish 250
	3, --  7 butterfly 250
	4, --  8 ace       125
	4, --  9 king      125
	4, -- 10 queen     100
	4, -- 11 jack      100
	4, -- 12 ten       100
	5, -- 13 nine      100+
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

local function reelgen(_, isbon)
	if isbon then
		return makereel(symsetbon, neighbours)
	else
		return makereel(symsetreg, neighbours)
	end
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen(1, false))
