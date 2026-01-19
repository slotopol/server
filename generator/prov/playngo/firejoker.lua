local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	4, --  1 seven  100
	4, --  2 bell   50
	5, --  3 melon  50
	4, --  4 plum   20
	4, --  5 orange 20
	5, --  6 lemon  20
	5, --  7 cherry 20
	2, --  8 bonus (1, 3, 5 reels only)
	1, --  9 joker
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0 }, --  1
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0 }, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0 }, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0 }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0 }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0 }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0 }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 2 }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 2, 2 }, --  9
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 2 or n == 4 then
		ss[8] = 0
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
