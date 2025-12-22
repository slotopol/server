local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, -- 1 seven
	3, -- 2 blueberry
	3, -- 3 strawberry
	3, -- 4 plum
	4, -- 5 pear
	5, -- 6 peach
	5, -- 7 cherry
	1, -- 8 bell
}

local function reelgen()
	return makereelhot(symset, 3, {[1]=true, [8]=true}, {})
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
