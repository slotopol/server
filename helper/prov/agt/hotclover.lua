local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	 0, --  1 wild (only on reel 2, 3, 4)
	 2, --  2 scatter
	 4, --  3 seven
	 8, --  4 strawberr
	 9, --  5 grapes
	 8, --  6 bar
	10, --  7 plum
	10, --  8 orange
	10, --  9 lemon
	10, -- 10 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 4, {[1]=true, [2]=true}, {}, true))
