local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	3, --  1 wild    1000
	1, --  2 scatter
	3, --  3 crown   400
	3, --  4 cup     400
	5, --  5 melon   200
	5, --  6 apple   200
	6, --  7 orange  100
	6, --  8 lemon   100
	6, --  9 plum    100
	7, -- 10 cherry  100
}

local chunklen = {
	3, --  1 wild
	1, --  2 scatter
	1, --  3 crown
	1, --  4 cup
	4, --  5 melon
	6, --  6 apple
	6, --  7 orange
	6, --  8 lemon
	6, --  9 plum
	6, -- 10 cherry
}

local function reelgen()
	return makereelhot(symset, 3, {[2]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
