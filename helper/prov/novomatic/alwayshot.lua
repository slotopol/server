local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	3, -- 1 seven
	3, -- 2 star
	4, -- 3 melon
	4, -- 4 grapes
	4, -- 5 bell
	6, -- 6 orange
	6, -- 7 plum
	6, -- 8 lemon
	6, -- 9 cherry
}

local chunklen = {
	6, -- 1 seven
	6, -- 2 star
	6, -- 3 melon
	6, -- 4 grapes
	6, -- 5 bell
	6, -- 6 orange
	6, -- 7 plum
	6, -- 8 lemon
	6, -- 9 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {}, chunklen))
