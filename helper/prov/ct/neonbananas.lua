local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

-- symbols set for 1-st reel
local symset1 = {
	2, --  1 seven
	2, --  2 dollar
	8, --  3 banana
	3, --  4 lucky slot
	2, --  5 grapes
	2, --  6 melon
	2, --  7 apple
	2, --  8 pear
	2, --  9 peach
	2, -- 10 orange
	2, -- 11 lemon
	2, -- 12 plum
	2, -- 13 cherry
}

local neighbours1 = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1 seven
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 dollar
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 banana
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4 lucky slot
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  5 grapes
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  6 melon
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  7 apple
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8 pear
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9 peach
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10 orange
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11 lemon
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12 plum
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13 cherry
}

-- symbols set for 2, 3, 4, 5 reels
local symset = {
	2, --  1 seven
	1, --  2 dollar
	7, --  3 banana
	5, --  4 lucky slot
	2, --  5 grapes
	2, --  6 melon
	2, --  7 apple
	2, --  8 pear
	2, --  9 peach
	2, -- 10 orange
	2, -- 11 lemon
	2, -- 12 plum
	2, -- 13 cherry
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1 seven
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 dollar
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 banana
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4 lucky slot
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  5 grapes
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  6 melon
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  7 apple
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8 pear
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9 peach
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10 orange
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11 lemon
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12 plum
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13 cherry
}

math.randomseed(os.time())

local function reelgen(n)
	if n == 1 then
		return makereel(symset1, neighbours1)
	else
		return makereel(symset, neighbours)
	end
end

if autoscan then
	return reelgen
end

print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4, 5"
printreel(reelgen(2))

local lsm = {
	10, 10, 10, 10, 10, 25, 25, 25, 50, 50, 50, 100, 200, 300, 1000,
}
local Els = 0
for _, v in ipairs(lsm) do
	Els = Els + v
end
Els = Els/#lsm
print(string.format("Els = %g, len = %d", Els, #lsm))