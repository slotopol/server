local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, -- 1 wild
	3, -- 2 bonsai
	5, -- 3 fish
	5, -- 4 fan
	9, -- 5 lamp
	10, -- 6 pot
	10, -- 7 flower
	13, -- 8 button
	1, -- 9 scatter
}

local chunklen = {
	1, -- 1 wild
	3, -- 2 bonsai
	3, -- 3 fish
	3, -- 4 fan
	5, -- 5 lamp
	5, -- 6 pot
	5, -- 7 flower
	5, -- 8 button
	1, -- 9 scatter
}

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {}, chunklen, true))
