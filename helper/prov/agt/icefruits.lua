local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	3, -- 1 wild
	3, -- 2 scatter
	8, -- 3 grapes
	11, -- 4 watermelon
	12, -- 5 plum
	14, -- 6 orange
	14, -- 7 pear
	15, -- 8 cherry
}

local chunklen = {
	4, -- 1 wild
	1, -- 2 scatter
	4, -- 3 grapes
	4, -- 4 watermelon
	4, -- 5 plum
	4, -- 6 orange
	4, -- 7 pear
	4, -- 8 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
