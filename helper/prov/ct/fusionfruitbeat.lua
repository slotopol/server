local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	 3, -- 1 wild    2000
	 3, -- 2 scatter
	 8, -- 3 apple   400
	10, -- 4 orange  200
	10, -- 5 melon   200
	15, -- 6 lemon   100
	15, -- 7 plum    100
	15, -- 8 cherry  100
}

local chunklen = {
	3, -- 1 wild
	1, -- 2 scatter
	4, -- 3 apple
	4, -- 4 orange
	4, -- 5 melon
	6, -- 6 lemon
	6, -- 7 plum
	6, -- 8 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
