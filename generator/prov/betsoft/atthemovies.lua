local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 oscar
	1, --  2 popcorn
	2, --  3 poster
	4, --  4 a
	5, --  5 dummy
	6, --  6 maw
	7, --  7 starship
	7, --  8 heart
	1, --  9 masks (2, 4 reels only)
	1, -- 10 projector
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, }, --  1
	{ 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, }, --  2
	{ 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, }, --  3
	{ 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 1, }, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, }, -- 10
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n ~= 2 and n ~= 4 then
		ss[9] = 0
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
