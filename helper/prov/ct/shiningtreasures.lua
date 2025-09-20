local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	3, --  1 wild    (only on reel 2, 3, 4)
	2, --  2 scatter
	4, --  3 seven   2500
	5, --  4 grape   400
	5, --  5 melon   400
	6, --  6 apple   300
	7, --  7 orange  100
	7, --  8 lemon   100
	7, --  9 plum    100
	7, -- 10 cherry  100
}

local chunklen = {
	3, --  1 wild
	1, --  2 scatter
	4, --  3 seven
	5, --  4 grape
	5, --  5 melon
	3, --  6 apple
	4, --  7 orange
	4, --  8 lemon
	4, --  9 plum
	4, -- 10 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
