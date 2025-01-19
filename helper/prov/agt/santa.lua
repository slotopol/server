local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset134 = {
	0, --  1 scatter (always 0 here)
	1, --  2 wild
	2, --  3 gnomes
	4, --  4 snowman
	4, --  5 christmas tree
	5, --  6 socks
	6, --  7 balls
	6, --  8 sweets
	6, --  9 present
	6, -- 10 bells
}

local symset2 = {
	1*1, --  1 scatter (always 1 here)
	1*4, --  2 wild
	2*4, --  3 gnomes
	2*4, --  4 snowman
	3*4-1, --  5 christmas tree
	4*4, --  6 socks
	4*4, --  7 balls
	5*4, --  8 sweets
	5*4, --  9 present
	6*4, -- 10 bells
}

local chunklen = {
	1, --  1 scatter
	1, --  2 wild
	1, --  3 gnomes
	1, --  4 snowman
	3, --  5 christmas tree
	3, --  6 socks
	3, --  7 balls
	4, --  8 sweets
	4, --  9 present
	4, -- 10 bells
}

math.randomseed(os.time())
printreel(makereelhot(symset134, 4, {}, chunklen, true))
printreel(makereelhot(symset2, 4, {}, chunklen, true))
