local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	4, -- 1 wild    (2, 3, 4 reels only)
	1, -- 2 scatter
	4, -- 3 seven   1000
	5, -- 4 dead    150
	5, -- 5 cat     80
	5, -- 6 vampire 80
	6, -- 7 pot     80
	6, -- 8 hat     80
	6, -- 9 scull   80
}

local chunklen = {
	4, -- 1 wild
	1, -- 2 scatter
	4, -- 3 seven
	4, -- 4 dead
	4, -- 5 cat
	4, -- 6 vampire
	4, -- 7 pot
	4, -- 8 hat
	4, -- 9 scull
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen))
