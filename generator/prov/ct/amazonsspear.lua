local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	3, --  3 brunette 2000
	3, --  4 blonde   625
	3, --  5 helmet   625
	3, --  6 harp     625
	3, --  7 ace      150
	4, --  8 king     150
	4, --  9 queen    100
	4, -- 10 jack     100
	4, -- 11 ten      100
	4, -- 12 nine     100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  1 wild (2, 3, 4 reels only)
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  3 brunette
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  4 blonde
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  5 helmet
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  6 harp
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 11 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 12 nine
}

local function reelgen(n)
	local n1 = symset[1]
	if n == 1 or n == 5 then
		symset[1] = 0
	end
	local reel, iter = makereel(symset, neighbours)
	symset[1] = n1
	return reel, iter
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
