local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset1 = {
	4, --  1 tomat   100
	4, --  2 corn    100
	4, --  3 lama    100
	4, --  4 frog    100
	4, --  5 jaguar  100
	3, --  6 condor  500
	3, --  7 queen   750
	3, --  8 king    1000
	1, --  9 dragon  10000
	2, -- 10 scatter
	0, -- 11 idol    (2, 3, 4 reel)
	0, -- 12 pyramid (3, 4, 5 reel)
}

local symset2 = {
	4, --  1 tomat   100
	4, --  2 corn    100
	4, --  3 lama    100
	4, --  4 frog    100
	4, --  5 jaguar  100
	3, --  6 condor  500
	3, --  7 queen   750
	3, --  8 king    1000
	1, --  9 dragon  10000
	1, -- 10 scatter
	1, -- 11 idol    (2, 3, 4 reel)
	0, -- 12 pyramid (3, 4, 5 reel)
}

local symset3 = {
	4, --  1 tomat   100
	4, --  2 corn    100
	4, --  3 lama    100
	4, --  4 frog    100
	4, --  5 jaguar  100
	3, --  6 condor  500
	3, --  7 queen   750
	3, --  8 king    1000
	1, --  9 dragon  10000
	2, -- 10 scatter
	1, -- 11 idol    (2, 3, 4 reel)
	2, -- 12 pyramid (3, 4, 5 reel)
}

local symset4 = {
	4, --  1 tomat   100
	4, --  2 corn    100
	4, --  3 lama    100
	4, --  4 frog    100
	4, --  5 jaguar  100
	4, --  6 condor  500
	4, --  7 queen   750
	2, --  8 king    1000
	1, --  9 dragon  10000
	1, -- 10 scatter
	1, -- 11 idol    (2, 3, 4 reel)
	2, -- 12 pyramid (3, 4, 5 reel)
}

local symset5 = {
	4, --  1 tomat   100
	4, --  2 corn    100
	4, --  3 lama    100
	4, --  4 frog    100
	4, --  5 jaguar  100
	4, --  6 condor  500
	4, --  7 queen   750
	2, --  8 king    1000
	1, --  9 dragon  10000
	2, -- 10 scatter
	0, -- 11 idol    (2, 3, 4 reel)
	2, -- 12 pyramid (3, 4, 5 reel)
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  1 tomat
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  2 corn
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, }, --  3 lama
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, }, --  4 frog
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, }, --  5 jaguar
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, }, --  6 condor
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  7 queen
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 2, 0, }, --  9 dragon
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, }, -- 10 scatter
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 0, }, -- 11 idol
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 2, }, -- 12 pyramid
}

math.randomseed(os.time())
printreel(makereel(symset1, neighbours))
printreel(makereel(symset2, neighbours))
printreel(makereel(symset3, neighbours))
printreel(makereel(symset4, neighbours))
printreel(makereel(symset5, neighbours))
