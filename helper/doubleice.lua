local path = arg[0]:match("(.*[/\\])")
dofile(path.."lib/reelgen.lua")

local symset = {
	3, -- 1 seven
	4, -- 2 strawberry
	4, -- 3 bell
	5, -- 4 star
	5, -- 5 lemon
	8, -- 6 blueberry
	8, -- 7 plum
	8, -- 8 orange
	8, -- 9 cherry
}

math.randomseed(os.time())
local reel, iter = makereelhot(symset, 5, {}, true)
printreel(reel, iter)
