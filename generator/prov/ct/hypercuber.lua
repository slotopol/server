local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset15 = {
	1, --  1 infinity 300
	0, --  2 wild     100 (always 0 here)
	4, --  3 atom     65
	5, --  4 red      20
	5, --  5 yellow   20
	5, --  6 gold     20
	5, --  7 violet   10
	5, --  8 lilac    10
	5, --  9 green    10
	5, -- 10 blue     10
	0, -- 11 cuber
}

local symset234 = {
	1, --  1 infinity 300
	2, --  2 wild     100 (2, 3, 4 reels only)
	4, --  3 atom     65
	5, --  4 red      20
	5, --  5 yellow   20
	5, --  6 gold     20
	5, --  7 violet   10
	5, --  8 lilac    10
	5, --  9 green    10
	5, -- 10 blue     10
	0, -- 11 cuber
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, --  1 infinity
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, --  2 wild (2, 3, 4 reels only)
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 atom
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4 red
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  5 yellow
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  6 gold
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  7 violet
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  8 lilac
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  9 green
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 10 blue
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 11 cuber
}

local function reelgen(n)
	if n == 1 or n == 5 then
		return makereel(symset15, neighbours)
	else
		local reel, iter = makereel(symset234, neighbours)
		--addsym(reel, 2, 3)
		return reel, iter
	end
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))

local cubermult = {2, 2, 2, 2, 3, 3, 3, 5, 5, 5, 10, 10, 10, 15, 15, 100}
local sum = 0
for _, v in ipairs(cubermult) do
	sum = sum + v
end
print("Mavr = "..sum/#cubermult, "len = "..#cubermult)
