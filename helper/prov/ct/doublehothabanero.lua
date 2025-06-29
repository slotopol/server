local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	3, --  1 wild     (only on reel 2, 3, 4)
	2, --  2 scatter
	6, --  3 woman    1000
	6, --  4 man      500
	8, --  5 crayfish 200
	8, --  6 shrimp   150
	8, --  7 ananas   50
	8, --  8 lime     50
	8, --  9 corn     50
	8, -- 10 banana   50
}

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	1, --  3 woman
	1, --  4 man
	3, --  5 crayfish
	3, --  6 shrimp
	4, --  7 ananas
	4, --  8 lime
	4, --  9 corn
	4, -- 10 banana
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
