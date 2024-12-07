local path = arg[0]:match("(.*[/\\])")
dofile(path.."lib/reelgen.lua")

local symset = {
	4, --  1 wild (only on reel 2, 3, 4)
	2, --  2 scatter
	2, --  3 seven
	5, --  4 strawberry
	5, --  5 blueberry
	7, --  6 pear
	9, --  7 plum
	9, --  8 peach
	9, --  9 quince
	9, -- 10 cherry
}

local chunklen = {
	3, --  1 wild
	1, --  2 scatter
	1, --  3 seven
	1, --  4 strawberry
	1, --  5 blueberry
	3, --  6 pear
	3, --  7 plum
	3, --  8 peach
	3, --  9 quince
	3, -- 10 cherry
}

math.randomseed(os.time())
local reel, iter = makereelhot(symset, 3, {[2]=true}, chunklen, true)
printreel(reel, iter)
