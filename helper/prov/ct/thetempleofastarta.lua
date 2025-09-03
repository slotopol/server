local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset1 = {
	1, --  1 wild    (on 2, 3, 4, 5 reels)
	1, --  2 scatter
	3, --  3 woman   1000
	3, --  4 kentaur 200
	3, --  5 cup     150
	3, --  6 horn    150
	4, --  7 ace     100
	4, --  8 king    100
	4, --  9 queen   100
	4, -- 10 jack    100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0,}, --  3 woman
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0,}, --  4 kentaur
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0,}, --  5 cup
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0,}, --  6 horn
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 10 jack
}

local symset2 = {
	4, --  1 wild    (on 2, 3, 4, 5 reels)
	1, --  2 scatter
	4, --  3 woman
	4, --  4 kentaur
	4, --  5 cup
	4, --  6 horn
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	4, --  3 woman
	4, --  4 kentaur
	4, --  5 cup
	4, --  6 horn
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local function tableglue(t1, t2)
	local t = {}
	for _, v in ipairs(t1) do
		table.insert(t, v)
	end
	for _, v in ipairs(t2) do
		table.insert(t, v)
	end
	return t
end

math.randomseed(os.time())
local reel1, iter1 = makereel(symset1, neighbours)
local reel2, iter2 = makereelhot(symset2, 3, {[2]=true}, chunklen, true)
print(string.format("iterations: %d, %d", iter1, iter2))
printreel(tableglue(reel1, reel2))
