local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	 4, -- 1 wild      1000
	 3, -- 2 scatter
	15, -- 3 seven     400
	16, -- 4 horseshoe 200
	17, -- 5 bell      200
	26, -- 6 grape     100
	26, -- 7 plum      100
	26, -- 8 cherry    100
}

local chunklen = {
	4, -- 1 wild
	1, -- 2 scatter
	6, -- 3 seven
	6, -- 4 horseshoe
	6, -- 5 bell
	6, -- 6 grape
	6, -- 7 plum
	6, -- 8 cherry
}

local function reelgen()
	return makereelhot(symset, 3, {[2]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
