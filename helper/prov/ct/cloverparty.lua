local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	 3, -- 1 wild      1000
	 3, -- 2 scatter
	 8, -- 3 seven     400
	10, -- 4 horseshoe 200
	10, -- 5 bell      200
	14, -- 6 peach     100
	14, -- 7 plum      100
	14, -- 8 cherry    100
}

local chunklen = {
	3, -- 1 wild
	1, -- 2 scatter
	4, -- 3 seven
	4, -- 4 horseshoe
	4, -- 5 bell
	6, -- 6 peach
	6, -- 7 plum
	6, -- 8 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
