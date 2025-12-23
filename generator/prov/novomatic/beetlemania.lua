local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 bee
	3, --  2 snail
	3, --  3 fly
	3, --  4 worm
	4, --  5 ace
	4, --  6 king
	5, --  7 queen
	5, --  8 jack
	5, --  9 ten
	1, -- 10 note
	0, -- 11 jazzbee (3 reel only)
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, --  1 bee
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 snail
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 fly
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  4 worm
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  5 ace
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  6 king
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7 queen
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  9 ten
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, -- 10 note
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, -- 11 jazzbee
}

local function reelgen(n)
	local function make()
		return makereel(symset, neighbours)
	end
	if n ~= 3 then
		local n11 = symset[11]
		symset[11] = 0
		local reel, iter = make()
		symset[11] = n11
		return reel, iter
	else
		return make()
	end
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 3, 4, 5"
printreel(reelgen(1))
print "reel 2"
printreel(reelgen(2))
