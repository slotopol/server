local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, -- 1 scatter
	2, -- 2 dollar  5000
	3, -- 3 melon   1000
	3, -- 4 apple   1000
	5, -- 5 orange  200
	5, -- 6 lemon   200
	6, -- 7 plum    200
	6, -- 8 cherry  200
}

local chunklen = {
	1, -- 1 scatter
	1, -- 2 dollar
	3, -- 3 melon
	3, -- 4 apple
	6, -- 5 orange
	6, -- 6 lemon
	6, -- 7 plum
	6, -- 8 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[1]=true}, chunklen, true))
