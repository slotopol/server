local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset15 = {
	1, --  1 scatter
	0, --  2 wild    1500 (always 0 here)
	3, --  3 heart   150
	4, --  4 blue    40
	4, --  5 green   40
	4, --  6 yellow  40
	5, --  7 melon   10
	5, --  8 jujube  10
	5, --  9 plum    10
	5, -- 10 cherry  10
}

local symset234 = {
	 3, --  1 scatter
	 0, --  2 wild    1500
	 8, --  3 heart   150
	10, --  4 blue    40
	10, --  5 green   40
	10, --  6 yellow  40
	12, --  7 melon   10
	12, --  8 jujube  10
	12, --  9 plum    10
	12, -- 10 cherry  10
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1 scatter
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 wild
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 heart
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4 blue
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  5 green
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  6 yellow
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  7 melon
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  8 jujube
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  9 plum
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 10 cherry
}

local function reelgen(n)
	if n == 1 or n == 5 then
		return makereel(symset15, neighbours)
	elseif n == 2 or n == 4 then
		local reel, iter = makereel(symset234, neighbours)
		addsym(reel, 2, 3)
		return reel, iter
	else -- n == 3
		local reel, iter = makereel(symset234, neighbours)
		local mt = getmetatable(reel)
		setmetatable(reel, nil)
		table.insert(reel, 5, 2)
		setmetatable(reel, mt)
		addsym(reel, 2, 3)
		return reel, iter
	end
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 4"
printreel(reelgen(2))
print "reel 3"
printreel(reelgen(3))
