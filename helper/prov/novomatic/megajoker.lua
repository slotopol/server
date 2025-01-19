local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset135 = {
	5, --  1 wild         2000
	4, --  2 red seven    400
	4, --  3 yellow seven 400
	5, --  4 strawberry   240
	5, --  5 pear         240
	6, --  6 grapes       160
	6, --  7 watermelon   160
	7, --  8 plum         100
	7, --  9 orange       100
	7, -- 10 lemon        100
	7, -- 11 cherry       100
	2, -- 12 star
}

local symset24 = {
	5, --  1 wild         2000
	4, --  2 red seven    400
	4, --  3 yellow seven 400
	4, --  4 strawberry   240
	4, --  5 pear         240
	4, --  6 grapes       160
	4, --  7 watermelon   160
	5, --  8 plum         100
	5, --  9 orange       100
	5, -- 10 lemon        100
	5, -- 11 cherry       100
	1, -- 12 star
}

local chunklen = {
	4, --  1 wild
	6, --  2 red seven
	6, --  3 yellow seven
	6, --  4 strawberry
	6, --  5 pear
	8, --  6 grapes
	8, --  7 watermelon
	8, --  8 plum
	8, --  9 orange
	8, -- 10 lemon
	8, -- 11 cherry
	1, -- 12 star
}

math.randomseed(os.time())
print "reel 1, 3, 5"
printreel(makereelhot(symset135, 4, {[12]=true}, chunklen))
print "reel 2, 4"
printreel(makereelhot(symset24, 4, {[12]=true}, chunklen))
