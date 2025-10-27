local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	4, --  1 wild    (2, 3, 4 reels only)
	2, --  2 scatter
	4, --  3 seven   1000
	5, --  4 bell    300
	5, --  5 shoe    200
	5, --  6 coin    200
	6, --  7 peach   100
	6, --  8 apple   100
	6, --  9 plum    100
	6, -- 10 cherry  100
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	3, --  3 seven
	3, --  4 bell
	3, --  5 shoe
	3, --  6 coin
	3, --  7 peach
	3, --  8 apple
	3, --  9 plum
	3, -- 10 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
