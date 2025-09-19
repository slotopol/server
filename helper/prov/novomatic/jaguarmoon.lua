local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symsetbon = {
	0, --  1 wild      (2, 3, 4 reel)
	0, --  2 scatter   (not used)
	3, --  3 wooman    800
	3, --  4 panther   200
	3, --  5 footprint 100
	3, --  6 rings     100
	3, --  7 ace       50
	3, --  8 king      50
	3, --  9 queen     50
	3, -- 10 jack      40
	3, -- 11 ten       40
	3, -- 12 nine      40
}

local symset = {
	0, --  1 wild      (2, 3, 4 reel)
	0, --  2 scatter   (insert directly)
	3, --  3 wooman    800
	3, --  4 panther   200
	3, --  5 footprint 100
	3, --  6 rings     100
	3, --  7 ace       50
	3, --  8 king      50
	3, --  9 queen     50
	3, -- 10 jack      40
	3, -- 11 ten       40
	3, -- 12 nine      40
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, }, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, }, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, }, --  3 wooman
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, }, --  4 panther
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, }, --  5 footprint
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, }, --  6 rings
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, }, -- 11 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, }, -- 12 nine
}

local scat = 2
local function ins1(reel)
	table.insert(reel, 1, scat)
	table.insert(reel, 1, scat)
end
local function ins2(reel)
	for i = 5, #reel do
		if neighbours[scat][reel[i - 2]] < 2 and
			neighbours[scat][reel[i - 1]] < 1 and
			neighbours[scat][reel[i + 1]] < 1 and
			neighbours[scat][reel[i - 2]] < 2 then
			table.insert(reel, i, scat)
			table.insert(reel, i, scat)
			break
		end
	end
end

math.randomseed(os.time())
printreel(makereel(symsetbon, neighbours))
local reel, iter = makereel(symset, neighbours)
printreel(reel, iter)
ins1(reel)
ins2(reel)
printreel(reel, iter)
