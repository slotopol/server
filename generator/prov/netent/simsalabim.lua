local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symsetreg = {
	3, --  1 hat   600
	3, --  2 chest 240
	3, --  3 cell  240
	3, --  4 cards 120
	4, --  5 ace   120
	4, --  6 king  90
	4, --  7 queen 60
	5, --  8 jack  60
	3, --  9 bonus
	1, -- 10 wild
	1, -- 11 scatter
}

local symsetbon = {
	3, --  1 hat   600
	3, --  2 chest 240
	3, --  3 cell  240
	3, --  4 cards 120
	3, --  5 ace   120
	4, --  6 king  90
	4, --  7 queen 60
	4, --  8 jack  60
	0, --  9 bonus (absent on free games)
	1, -- 10 wild
	1, -- 11 scatter
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 11
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

local dd = {10, 20, 20, 40, 40, 80, 120, 120, 450}
local sum = 0
for _, v in ipairs(dd) do
	sum = sum + v
end
print("bonus M = "..sum/#dd..", E = 3*M = "..3*sum/#dd)
