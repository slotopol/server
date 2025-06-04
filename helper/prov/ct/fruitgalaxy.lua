local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 wild (on 2, 3, 4 reels)
	1, --  2 UFO (on all reels)
	2, --  3 banana (on 1, 3, 5 reels)
	2, --  4 seven  3000
	3, --  5 apple  500
	3, --  6 orange 500
	5, --  7 bell   200
	6, --  8 melon  100
	6, --  9 lemon  100
	6, -- 10 plum   100
	6, -- 11 cherry 100
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

math.randomseed(os.time())
printreel(makereelhot(symset, 4, {[1]=true, [2]=true, [3]=true}, chunklen))
