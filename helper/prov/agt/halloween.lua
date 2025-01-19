local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	3, -- 1 pumpkin   1000
	3, -- 2 witch     500
	4, -- 3 castle    200
	5, -- 4 scarecrow 100
	7, -- 5 ghost     30
	8, -- 6 spider    20
	9, -- 7 skeleton  10
	10, -- 8 candles   5
}

local chunklen = {
	1, -- 1 pumpkin
	1, -- 2 witch
	9, -- 3 castle
	9, -- 4 scarecrow
	9, -- 5 ghost
	9, -- 6 spider
	9, -- 7 skeleton
	9, -- 8 candles
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {}, chunklen, true))
