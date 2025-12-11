local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild (2, 4 reels only)
	1, --  2 scatter
	2, --  3 parrot   5000
	2, --  4 turtle   1000
	3, --  5 seastar  500
	3, --  6 seahorse 500
	3, --  7 ace      250
	3, --  8 king     250
	3, --  9 queen    150
	3, -- 10 jack     150
	3, -- 11 ten      100
	3, -- 12 nine     100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  1 wild (2, 4 reels only)
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  3 parrot
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  4 turtle
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  5 seastar
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  6 seahorse
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 11 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 12 nine
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		return makereel(symset, neighbours)
	end
	if n == 1 or n == 3 or n == 5 then
		local n1 = symset[1]
		symset[1] = 0
		local reel, iter = make()
		symset[1] = n1
		return reel, iter
	else
		return make()
	end
end

if autoscan then
	return reelgen
end

print "reel 1, 3, 5"
printreel(reelgen(1))
print "reel 2, 4"
printreel(reelgen(2))
