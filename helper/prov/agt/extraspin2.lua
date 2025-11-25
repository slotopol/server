local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset15 = {
	2, -- 1 wild       1000
	0, -- 2 scatter    (always 0 here)
	6, -- 3 strawberry 500
	6, -- 4 papaya     200
	6, -- 5 grapes     200
	7, -- 6 orange     100
	7, -- 7 plum       100
	9, -- 8 cherry     50
	9, -- 9 pear       50
}

local symset234 = {
	 4, -- 1 wild       1000
	 1, -- 2 scatter    (always 1 here)
	12, -- 3 strawberry 500
	12, -- 4 papaya     200
	12, -- 5 grapes     200
	14, -- 6 orange     100
	14, -- 7 plum       100
	18, -- 8 cherry     50
	18, -- 9 pear       50
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
	5, --  6 orange
	5, --  7 plum
	5, --  8 cherry
	5, --  9 pear
}

math.randomseed(os.time())
printreel(makereelhot(symset15, 3, {[1]=true, [2]=true}, chunklen15, true))
printreel(makereelhot(symset234, 3, {[1]=true, [2]=true}, chunklen234, true))
