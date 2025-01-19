local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 wild (only on reel 2, 3, 4)
	1, --  2 scatter
	2, --  3 seven
	3, --  4 strawberry
	4, --  5 blueberry
	6, --  6 pear
	7, --  7 plum
	7, --  8 peach
	7, --  9 papaya
	7, -- 10 cherry
}

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	1, --  3 seven
	1, --  4 strawberry
	1, --  5 blueberry
	3, --  6 pear
	3, --  7 plum
	3, --  8 peach
	3, --  9 papaya
	3, -- 10 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[1]=true, [2]=true}, chunklen, true))
