local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	1, -- 1 wild    (2, 4 reel)
	1, -- 2 scatter
	4, -- 3 heart   1000
	7, -- 4 sun     300
	7, -- 5 beer    300
	11, -- 6 pizza   100
	11, -- 7 bomb    100
	11, -- 8 flower  100
}

local chunklen = {
	1, -- 1 wild
	1, -- 2 scatter
	1, -- 3 heart
	3, -- 4 sun
	3, -- 5 beer
	3, -- 6 pizza
	3, -- 7 bomb
	3, -- 8 flower
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[1]=true, [2]=true}, chunklen, true))
