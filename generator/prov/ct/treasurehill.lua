local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	0, --  1 wild      1000 (insert directly)
	0, --  2 scatter   (2, 3, 4 reel)
	3, --  3 clover    400+
	3, --  4 horseshoe 400+
	3, --  5 treasure  400
	3, --  6 rainbow   400
	5, --  7 beer      200
	5, --  8 smoke     200
	6, --  9 ace       100
	6, -- 10 king      100
	6, -- 11 queen     100
	6, -- 12 jack      100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1 wild
	{ 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2 scatter (2, 3, 4 reel)
	{ 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  3 clover
	{ 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, }, --  4 horseshoe
	{ 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, }, --  5 treasure
	{ 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, }, --  6 rainbow
	{ 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, }, --  7 beer
	{ 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, }, --  8 smoke
	{ 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, }, --  9 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, }, -- 10 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, }, -- 11 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, }, -- 12 jack
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		local reel, iter = makereel(symset, neighbours)
		addsym(reel, 1, 4)
		return reel, iter
	end
	if n == 1 or n == 5 then
		local n2 = symset[2]
		symset[2] = 0
		local reel, iter = make()
		symset[2] = n2
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
