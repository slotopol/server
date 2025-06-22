local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset1 = {
	3, --  1 wild          (on 2, 3, 4, 5 reels)
	2, --  2 scatter
	3, --  3 singer
	3, --  4 dancer man
	3, --  5 dancer girl 1
	3, --  6 dancer girl 2
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0,}, --  3 singer
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0,}, --  4 dancer man
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0,}, --  5 dancer girl 1
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0,}, --  6 dancer girl 2
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 10 jack
}

local symset2 = {
	3, --  1 wild          (on 2, 3, 4, 5 reels)
	3, --  2 scatter
	3, --  3 singer
	3, --  4 dancer man
	3, --  5 dancer girl 1
	3, --  6 dancer girl 2
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local chunklen = {
	3, --  1 wild
	3, --  2 scatter
	6, --  3 singer
	6, --  4 dancer man
	6, --  5 dancer girl 1
	6, --  6 dancer girl 2
	6, --  7 ace
	6, --  8 king
	6, --  9 queen
	6, -- 10 jack
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
