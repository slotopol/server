local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, -- 1 seven
	3, -- 2 blueberry
	3, -- 3 strawberry
	3, -- 4 plum
	4, -- 5 pear
	5, -- 6 peach
	5, -- 7 cherry
	1, -- 8 bell
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[1]=true, [8]=true}, {}))
