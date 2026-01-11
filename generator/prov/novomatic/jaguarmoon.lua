local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symsetbon = {
	1, --  1 wild      (2, 3, 4 reels only)
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

local symsetreg = {
	1, --  1 wild      (2, 3, 4 reels only)
	0, --  2 scatter   (insert directly)
	2, --  3 wooman    800
	2, --  4 panther   200
	3, --  5 footprint 100
	3, --  6 rings     100
	3, --  7 ace       50
	3, --  8 king      50
	3, --  9 queen     50
	3, -- 10 jack      40
	4, -- 11 ten       40
	4, -- 12 nine      40
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
	local mt = getmetatable(reel)
	setmetatable(reel, nil)
	table.insert(reel, 1, scat)
	table.insert(reel, 1, scat)
	setmetatable(reel, mt)
end
local function ins2(reel)
	local mt = getmetatable(reel)
	setmetatable(reel, nil)
	for i = 5, rawlen(reel) do
		if neighbours[scat][reel[i - 2]] < 2 and
			neighbours[scat][reel[i - 1]] < 1 and
			neighbours[scat][reel[i + 1]] < 1 and
			neighbours[scat][reel[i - 2]] < 2 then
			table.insert(reel, i, scat)
			table.insert(reel, i, scat)
			break
		end
	end
	setmetatable(reel, mt)
end

local function reelgen(n, isbon)
	if isbon then
		local symset = symsetbon
		local function make()
			return makereel(symset, neighbours)
		end
		if n == 1 or n == 5 then
			local n1 = symset[1]
			symset[1] = 0
			local reel, iter = make()
			symset[1] = n1
			return reel, iter
		else
			return make()
		end
	else
		local symset = symsetreg
		local n1, n2 = symset[1], symset[2]
		if n == 1 or n == 5 then
			symset[1] = 0
		end
		local reel, iter = makereel(symset, neighbours)
		if n == 1 or n == 2 or n == 3 then
			ins1(reel)
			if n == 2 then
				ins2(reel)
			end
		end
		symset[1], symset[2] = n1, n2
		return reel, iter
	end
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
local isbon = false
printreel(reelgen(1, isbon))
printreel(reelgen(2, isbon))
printreel(reelgen(3, isbon))
printreel(reelgen(4, isbon))
printreel(reelgen(5, isbon))
