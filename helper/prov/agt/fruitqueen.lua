local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	 1, -- 1 scatter
	 4, -- 2 wild       1200
	10, -- 3 grapes     400
	12, -- 4 strawberry 200
	14, -- 5 plum       200
	16, -- 6 pear       120
	18, -- 7 orange     100
	20, -- 8 cherry     80
}

local chunklen = { -- set all to full length
	  1, -- 1 scatter
	100, -- 2 wild
	100, -- 3 grapes
	100, -- 4 strawberry
	100, -- 5 plum
	100, -- 6 pear
	100, -- 7 orange
	100, -- 8 cherry
}

math.randomseed(os.time())
printreel(makereelhot(symset, 6, {[1]=true}, chunklen, true))
