local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild (on 2, 3, 4 reels)
	1, --  2 UFO (on all reels)
	2, --  3 banana (on 1, 3, 5 reels)
	2, --  4 seven  3000
	4, --  5 apple  500
	4, --  6 orange 500
	9, --  7 bell   200
	9, --  8 melon  100
	9, --  9 lemon  100
	9, -- 10 plum   100
	9, -- 11 cherry 100
}

local chunklen = {
	1, --  1 wild (on 2, 3, 4 reels)
	1, --  2 UFO (on all reels)
	1, --  3 banana (on 1, 3, 5 reels)
	1, --  4 seven
	3, --  5 apple
	3, --  6 orange
	3, --  7 bell
	6, --  8 melon
	6, --  9 lemon
	6, -- 10 plum
	6, -- 11 cherry
}

local scat = {[1]=true, [2]=true, [3]=true}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	elseif n == 2 or n == 4 then
		ss[3] = 0
	end
	return makereelhot(ss, 4, scat, chunklen)
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
