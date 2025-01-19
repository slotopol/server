local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	3, -- 1 seven
	4, -- 2 strawberry
	4, -- 3 bell
	5, -- 4 star
	5, -- 5 lemon
	8, -- 6 blueberry
	8, -- 7 plum
	8, -- 8 orange
	8, -- 9 cherry
}

local chunklen = {
	3, -- 1 seven
	4, -- 2 strawberry
	4, -- 3 bell
	5, -- 4 star
	5, -- 5 lemon
	5, -- 6 blueberry
	5, -- 7 plum
	5, -- 8 orange
	5, -- 9 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {}, chunklen, true))
