local path = arg[0]:match("(.*[/\\])")
dofile(path.."lib/reelgen.lua")

local symset = {
	4, --  1 wild (only on reel 2, 3, 4)
	1, --  2 scatter
	5, --  3 strawberry
	5, --  4 bell
	5, --  5 greenstar
	5, --  6 redstar
	6, --  7 plum
	6, --  8 peach
	6, --  9 quince
	6, -- 10 cherry
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	8, --  3 strawberry
	8, --  4 bell
	8, --  5 greenstar
	8, --  6 redstar
	8, --  7 plum
	8, --  8 peach
	8, --  9 quince
	8, -- 10 cherry
}

math.randomseed(os.time())
local reel, iter = makereelhot(symset, 4, {[2]=true}, chunklen, true)
printreel(reel, iter)
