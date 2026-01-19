local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

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

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	elseif n == 2 or n == 4 then
		ss[3] = 0
	end
	return makereelhot(ss, 3, scat, chunklen)
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
