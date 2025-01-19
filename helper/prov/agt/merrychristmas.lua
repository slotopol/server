local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	3, -- 1 snowman 500
	2, -- 2 scatter
	4, -- 3 ice     250
	6, -- 4 sled    100
	10, -- 5 house  20
	10, -- 6 bell   10
	12, -- 7 deer   5
}

local chunklen = {
	6, -- 1 snowman
	1, -- 2 scatter
	6, -- 3 ice
	6, -- 4 sled
	7, -- 5 house
	7, -- 6 bell
	7, -- 7 deer
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
