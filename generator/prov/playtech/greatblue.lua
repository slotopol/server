local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild     10000
	3, --  2 dolphin  750
	3, --  3 turtle   750
	4, --  4 fish     400
	4, --  5 seahorse 250
	4, --  6 starfish 250
	4, --  7 ace      150
	4, --  8 king     150
	5, --  9 queen    100
	5, -- 10 jack     100
	5, -- 11 ten      100
	4, -- 12 nine     100
	1, -- 13 scatter
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 2,}, --  1
	{ 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  2
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  3
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  4
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, 0,}, --  5
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0,}, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 2 or n == 4 then
		ss[1] = ss[1] + 1
	end
	return makereel(ss, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 3, 5"
printreel(reelgen(1))
print "reel 2, 4"
printreel(reelgen(2))
