local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset1 = {
	4, -- 1 bonus
	2, -- 2 leprechaun 5000
	2, -- 3 clover     1000
	2, -- 4 pot        500
	2, -- 5 horseshoe  250
	2, -- 6 bell       150
}

local neighbours = {
	--1, 2, 3, 4, 5, 6,
	{ 1, 0, 0, 0, 0, 0,}, -- 1 bonus
	{ 0, 1, 1, 0, 0, 0,}, -- 2 leprechaun
	{ 0, 1, 1, 0, 0, 0,}, -- 3 clover
	{ 0, 0, 0, 1, 0, 0,}, -- 4 pot
	{ 0, 0, 0, 0, 1, 0,}, -- 5 horseshoe
	{ 0, 0, 0, 0, 0, 1,}, -- 6 bell
}

local symset2 = {
	3, -- 1 bonus
	6, -- 2 leprechaun
	6, -- 3 clover
	6, -- 4 pot
	6, -- 5 horseshoe
	6, -- 6 bell
}

local chunklen = {
	3, -- 1 bonus
	6, -- 2 leprechaun
	6, -- 3 clover
	6, -- 4 pot
	6, -- 5 horseshoe
	6, -- 6 bell
}

math.randomseed(os.time())

local function reelgen()
	local reel1, iter1 = makereel(symset1, neighbours)
	local reel2, iter2 = makereelhot(symset2, 3, {}, chunklen, true)
	return reelglue(reel1, reel2), iter1, iter2
end

if autoscan then
	return reelgen
end

printreel(reelgen())
