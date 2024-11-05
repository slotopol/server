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

math.randomseed(os.time())
local reel, iter = makereelhot(symset, 8, {[2]=true}, true)
printreel(reel, iter)
