local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	3, -- 1 strawberry 1000
	4, -- 2 blueberry  200
	5, -- 3 plum       40
	5, -- 4 pear       40
	5, -- 5 peach      40
	5, -- 6 cherry     32
	6, -- 7 apple      20
	2, -- 8 scatter
}

local chunklen = {
	1, -- 1 strawberry
	1, -- 2 blueberry
	6, -- 3 plum
	6, -- 4 pear
	6, -- 5 peach
	6, -- 6 cherry
	6, -- 7 apple
	1, -- 8 scatter
}

local function reelgen()
	return makereelhot(symset, 3, {[8]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
