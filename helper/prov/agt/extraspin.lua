local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset15 = {
	2, --  1 wild
	0, --  2 scatter (always 0 here)
	3, --  3 strawberry
	4, --  4 papaya
	4, --  5 grapes
	7, --  6 orange
	7, --  7 plum
	8, --  8 cherry
	8, --  9 pear
}

local symset234 = {
	4, --  1 wild
	1, --  2 scatter (always 1 here)
	6, --  3 strawberry
	10, --  4 papaya
	10, --  5 grapes
	13, --  6 orange
	13, --  7 plum
	15, --  8 cherry
	15, --  9 pear
}

local chunklen15 = {
	1, --  1 wild
	1, --  2 scatter
	1, --  3 strawberry
	1, --  4 papaya
	1, --  5 grapes
	3, --  6 orange
	3, --  7 plum
	3, --  8 cherry
	3, --  9 pear
}

local chunklen234 = {
	1, --  1 wild
	1, --  2 scatter
	1, --  3 strawberry
	1, --  4 papaya
	1, --  5 grapes
	7, --  6 orange
	7, --  7 plum
	5, --  8 cherry
	5, --  9 pear
}

math.randomseed(os.time())
printreel(makereelhot(symset15, 3, {[1]=true, [2]=true}, chunklen15, true))
printreel(makereelhot(symset234, 3, {[1]=true, [2]=true}, chunklen234, true))
