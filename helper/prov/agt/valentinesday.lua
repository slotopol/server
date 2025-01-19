local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, -- 1 scatter
	2, -- 2 angel     5000
	4, -- 3 nymph     1000
	5, -- 4 soul      200
	5, -- 5 toy       200
	5, -- 6 balloon   200
	5, -- 7 hearts    160
	6, -- 8 medallion 100
}

local chunklen = {
	1, -- 1 scatter
	1, -- 2 angel
	1, -- 3 nymph
	6, -- 4 soul
	6, -- 5 toy
	6, -- 6 balloon
	6, -- 7 hearts
	6, -- 8 medallion
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[1]=true}, chunklen))
