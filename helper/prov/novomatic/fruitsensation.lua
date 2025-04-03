local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, -- 1 seven  5000
	5, -- 2 bells  500
	5, -- 3 melon  500
	7, -- 4 plum   200
	7, -- 5 orange 200
	7, -- 6 lemon  200
	7, -- 7 cherry 200
}

local chunklen = {
	1, -- 1 seven
	3, -- 2 bells
	3, -- 3 melon
	3, -- 4 plum
	3, -- 5 orange
	3, -- 6 lemon
	3, -- 7 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {}, chunklen, true))
