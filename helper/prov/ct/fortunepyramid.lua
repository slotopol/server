local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	 3, -- 1 wild    1000
	 2, -- 2 scatter
	 6, -- 3 scarab  500
	 9, -- 4 cross   200
	 9, -- 5 eye     200
	11, -- 6 ring    100
	11, -- 7 cup     100
	11, -- 8 bowl    100
}

local chunklen = {
	3, -- 1 wild
	1, -- 2 scatter
	3, -- 3 scarab
	4, -- 4 cross
	4, -- 5 eye
	5, -- 6 ring
	5, -- 7 cup
	5, -- 8 bowl
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
