local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	0, --  1 wild      1000 (insert directly)
	3, --  2 scatter   (2, 3, 4 reels only)
	4, --  3 fox       400
	4, --  4 squirrel  400
	4, --  5 owl       250
	4, --  6 eagle     250
	4, --  7 rockfish  150
	5, --  8 bulltrout 150
	5, --  9 ace       100
	5, -- 10 king      100
	5, -- 11 queen     100
	5, -- 12 jack      100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1 wild
	{ 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2 scatter
	{ 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  3 fox
	{ 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, }, --  4 squirrel
	{ 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, }, --  5 owl
	{ 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, }, --  6 eagle
	{ 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, }, --  7 rockfish
	{ 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, }, --  8 bulltrout
	{ 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, }, --  9 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, }, -- 10 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, }, -- 11 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, }, -- 12 jack
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[2] = 0
	end
	local reel, iter = makereel(ss, neighbours)
	addsym(reel, 1, 4)
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
