local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	4, --  1 wild (only on reel 2, 3, 4)
	2, --  2 scatter
	5, --  3 strawberry
	5, --  4 pear
	5, --  5 greenstar
	5, --  6 redstar
	6, --  7 plum
	6, --  8 peach
	6, --  9 papaya
	6, -- 10 cherry
}

local chunklen = {
	3, --  1 wild
	1, --  2 scatter
	3, --  3 strawberry
	3, --  4 pear
	3, --  5 greenstar
	3, --  6 redstar
	3, --  7 plum
	3, --  8 peach
	3, --  9 papaya
	3, -- 10 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
