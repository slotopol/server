local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	1, --  1 wild   (on 2, 3, 4 reels)
	1, --  2 star   (on all reels)
	2, --  3 dollar (on 1, 3, 5 reels)
	2, --  4 seven  3000
	3, --  5 apple  500
	3, --  6 orange 500
	6, --  7 bell   200
	7, --  8 melon  100
	7, --  9 lemon  100
	7, -- 10 plum   100
	7, -- 11 cherry 100
}

local chunklen = {
	1, --  1 wild (on 2, 3, 4 reels)
	1, --  2 star (on all reels)
	1, --  3 dollar (on 1, 3, 5 reels)
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
