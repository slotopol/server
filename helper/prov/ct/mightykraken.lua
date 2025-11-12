local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

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

do
	print "reel 1, 5"
	local n1 = symset[1]
	symset[1] = 0
	printreel(makereelhot(symset, 3, scat, chunklen))
	symset[1] = n1
end

do
	print "reel 2, 4"
	local n3 = symset[3]
	symset[3] = 0
	printreel(makereelhot(symset, 3, scat, chunklen))
	symset[3] = n3
end

do
	print "reel 3"
	printreel(makereelhot(symset, 3, scat, chunklen))
end
