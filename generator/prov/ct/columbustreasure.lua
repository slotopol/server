local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild1    10000
	1, --  2 wild2 (2, 4 reels only)
	1, --  3 scatter
	2, --  4 cardinal 1000
	3, --  5 wizard   250
	3, --  6 sailor   250
	4, --  7 lady     100
	4, --  8 knight   100
	5, --  9 ace      25
	5, -- 10 king     25
	5, -- 11 queen    25
	5, -- 12 jack     25
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1 wild1
	{ 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 wild2 (2, 4 reels only)
	{ 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 scatter
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4 cardinal
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  5 wizard
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  6 sailor
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  7 lady
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  8 knight
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  9 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 10 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 11 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 12 jack
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n ~= 2 and n ~= 4 then
		ss[1] = 0
	end
	return makereel(ss, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen(1))
printreel(reelgen(2))
printreel(reelgen(3))
printreel(reelgen(4))
printreel(reelgen(5))
