local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	0, --  1 moon wolf  1000 (insert directly)
	3, --  2 grey wolf  400
	3, --  3 white wolf 400
	3, --  4 idol1      250
	3, --  5 idol2      250
	5, --  6 ace        150
	5, --  7 king       150
	5, --  8 queen      100
	5, --  9 jack       100
	5, -- 10 ten        100
	5, -- 11 nine       100
	3, -- 12 bonus      (2, 3, 4 reels only)
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 3, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 3,}, --  1 moon wolf
	{ 1, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  2 grey wolf
	{ 1, 1, 3, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  3 white wolf
	{ 1, 1, 1, 3, 1, 0, 0, 0, 0, 0, 0, 0,}, --  4 idol1
	{ 1, 1, 1, 1, 3, 0, 0, 0, 0, 0, 0, 0,}, --  5 idol2
	{ 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0,}, --  6 ace
	{ 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0,}, --  7 king
	{ 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0,}, --  8 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0,}, --  9 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0,}, -- 10 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0,}, -- 11 nine
	{ 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3,}, -- 12 bonus
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		local reel, iter = makereel(symset, neighbours)
		addsym(reel, 1, 4)
		return reel, iter
	end
	if n == 1 or n == 5 then
		local n12 = symset[12]
		symset[12] = 0
		local reel, iter = make()
		symset[12] = n12
		return reel, iter
	else
		return make()
	end
end

if autoscan then
	return reelgen
end

print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
