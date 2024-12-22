local path = arg[0]:match("(.*[/\\])")
dofile(path.."lib/reelgen.lua")

local symset = {
	1, -- 1 wild
	3, -- 2 scatter
	1, -- 3 snowman
	3, -- 4 ice
	6, -- 5 sled
	10, -- 6 house
	14, -- 7 bell
	16, -- 8 deer
}

local chunklen = {
	6, -- 1 wild
	1, -- 2 scatter
	6, -- 3 snowman
	6, -- 4 ice
	6, -- 5 sled
	6, -- 6 house
	6, -- 7 bell
	6, -- 8 deer
}

math.randomseed(os.time())
local reel, iter = makereelhot(symset, 3, {[2]=true}, chunklen, true)
printreel(reel, iter)
