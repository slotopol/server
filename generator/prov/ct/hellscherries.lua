local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 seven  500
	2, --  2 bar    250
	2, --  3 melon  200
	3, --  4 bell   100
	4, --  5 apple  50
	4, --  6 pear   50
	4, --  7 plum   50
	5, --  8 lemon  10
	5, --  9 orange 10
	6, -- 10 cherry 10
}

local chunklen = {
	1, --  1 seven
	1, --  2 bar
	1, --  3 melon
	2, --  4 bell
	6, --  5 apple
	6, --  6 pear
	6, --  7 plum
	6, --  8 lemon
	6, --  9 orange
	6, -- 10 cherry
}

local function reelgen()
	return makereelhot(symset, 3, {[1]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
