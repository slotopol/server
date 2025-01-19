local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, -- 1 seven
	2, -- 2 strawberry
	3, -- 3 grapes
	5, -- 4 plum
	5, -- 5 pear
	5, -- 6 cherry
	1, -- 7 star
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {}, {}))
