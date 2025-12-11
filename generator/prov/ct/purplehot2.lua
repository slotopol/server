local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, -- 1 scatter
	2, -- 2 seven   5000
	3, -- 3 apple   500
	3, -- 4 orange  500
	5, -- 5 plum    200
	5, -- 6 lemon   200
	5, -- 7 melon   200
	5, -- 8 cherry  200
}

local chunklen = {
	1, -- 1 scatter
	1, -- 2 seven
	3, -- 3 apple
	3, -- 4 orange
	4, -- 5 plum
	4, -- 6 lemon
	4, -- 7 melon
	4, -- 8 cherry
}

math.randomseed(os.time())

local function reelgen()
	return makereelhot(symset, 3, {[1]=true}, chunklen)
end

if autoscan then
	return reelgen
end

printreel(reelgen())
