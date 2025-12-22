local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	3, -- 1 wild       1000
	2, -- 2 scatter
	10, -- 3 diamond    400
	12, -- 4 heliodor   200
	12, -- 5 aquamarine 200
	14, -- 6 sapphire   100
	14, -- 7 emerald    100
	16, -- 8 topaz      80
}

local chunklen = {
	3, -- 1 wild
	1, -- 2 scatter
	4, -- 3 diamond
	4, -- 4 heliodor
	4, -- 5 aquamarine
	5, -- 6 sapphire
	5, -- 7 emerald
	6, -- 8 topaz
}

local function reelgen()
	return makereelhot(symset, 3, {[2]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
