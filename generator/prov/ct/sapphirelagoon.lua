local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset1 = {
	2, --  1 wild    (2, 3, 4, 5 reels only)
	1, --  2 scatter
	2, --  3 man
	3, --  4 woman
	3, --  5 flask
	3, --  6 hook
	3, --  7 ace
	3, --  8 king
	3, --  9 queen
	3, -- 10 jack
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  3 man
	{ 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  4 woman
	{ 1, 1, 0, 0, 2, 1, 0, 0, 0, 0,}, --  5 flask
	{ 1, 1, 0, 0, 1, 2, 0, 0, 0, 0,}, --  6 hook
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 10 jack
}

local symset2 = {
	4, --  1 wild    (2, 3, 4, 5 reels only)
	3, --  2 scatter
	4, --  3 man
	4, --  4 woman
	4, --  5 flask
	4, --  6 hook
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local chunklen = {
	6, --  1 wild
	3, --  2 scatter
	6, --  3 singer
	6, --  4 dancer man
	6, --  5 dancer girl 1
	6, --  6 dancer girl 2
	6, --  7 ace
	6, --  8 king
	6, --  9 queen
	6, -- 10 jack
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		local reel1, iter1 = makereel(symset1, neighbours)
		local reel2, iter2 = makereelhot(symset2, 3, {[2]=true}, chunklen)
		return reelglue(reel1, reel2), iter1, iter2
	end
	if n == 1 then
		local n11, n21 = symset1[1], symset2[1]
		symset1[1], symset2[1] = 0, 0
		local reel, iter = make()
		symset1[1], symset2[1] = n11, n21
		return reel, iter
	else
		return make()
	end
end

if autoscan then
	return reelgen
end

print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4, 5"
printreel(reelgen(2))
