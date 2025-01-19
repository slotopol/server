local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	4, --  1 wild (only on reel 2, 3, 4)
	1, --  2 scatter
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
	4, --  1 wild
	1, --  2 scatter
	8, --  3 strawberry
	8, --  4 pear
	8, --  5 greenstar
	8, --  6 redstar
	8, --  7 plum
	8, --  8 peach
	8, --  9 papaya
	8, -- 10 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 4, {[2]=true}, chunklen, true))
