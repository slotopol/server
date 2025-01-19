local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	4, -- 1 seven
	5, -- 2 melon
	5, -- 3 grapes
	7, -- 4 plum
	7, -- 5 orange
	7, -- 6 lemon
	7, -- 7 cherry
}

local chunklen = {
	5, -- 1 seven
	6, -- 2 melon
	6, -- 3 grapes
	6, -- 4 plum
	6, -- 5 orange
	6, -- 6 lemon
	6, -- 7 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {}, chunklen, true))
