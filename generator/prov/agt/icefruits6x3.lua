local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, -- 1 wild    2000
	1, -- 2 scatter
	4, -- 3 grapes  500
	5, -- 4 melon  300
	5, -- 5 plum   300
	5, -- 6 orange 150
	5, -- 7 pear   150
	5, -- 8 cherry 150
}

local chunklen = {
	1, -- 1 wild
	1, -- 2 scatter
	1, -- 3 grapes
	1, -- 4 melon
	4, -- 5 plum
	4, -- 6 orange
	4, -- 7 pear
	4, -- 8 cherry
}

local function reelgen()
	return makereelhot(symset, 3, {[2]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
