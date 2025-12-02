local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	1, --  1 wild     (on 2, 3, 4 reels)
	1, --  2 star     (on all reels)
	2, --  3 dollar   (on 1, 3, 5 reels)
	2, --  4 captain  3000
	4, --  5 compass  500
	4, --  6 anchor   500
	5, --  7 bell     200
	6, --  8 spades   100
	6, --  9 hearts   100
	6, -- 10 diamonds 100
	6, -- 11 clubs    100
}

local chunklen = {
	1, --  1 wild (on 2, 3, 4 reels)
	1, --  2 star (on all reels)
	1, --  3 dollar (on 1, 3, 5 reels)
	1, --  4 captain
	1, --  5 compass
	1, --  6 anchor
	1, --  7 bell
	3, --  8 spades
	3, --  9 hearts
	3, -- 10 diamonds
	3, -- 11 clubs
}

local scat = {[1]=true, [2]=true, [3]=true}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		return makereelhot(symset, 3, scat, chunklen)
	end
	if n == 1 or n == 5 then
		local n1 = symset[1]
		symset[1] = 0
		local reel, iter = make()
		symset[1] = n1
		return reel, iter
	elseif n == 2 or n == 4 then
		local n3 = symset[3]
		symset[3] = 0
		local reel, iter = make()
		symset[3] = n3
		return reel, iter
	else -- n == 3
		return make()
	end
end

if autoscan then
	return reelgen
end

print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 4"
printreel(reelgen(2))
print "reel 3"
printreel(reelgen(3))
