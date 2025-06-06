local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	4, --  1 wild      (only on reel 2, 3, 4)
	2, --  2 scatter
	4, --  3 clover    1500
	5, --  4 horseshoe 1250
	5, --  5 bell      500
	6, --  6 apple     100+
	6, --  7 orange    100
	6, --  8 lemon     100
	6, --  9 plum      100
	6, -- 10 cherry    100
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	5, --  3 clover
	5, --  4 horseshoe
	3, --  5 bell
	3, --  6 apple
	3, --  7 orange
	3, --  8 lemon
	3, --  9 plum
	3, -- 10 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
