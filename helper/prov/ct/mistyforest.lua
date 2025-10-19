local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset1 = {
	1, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	3, --  3 woman   1000
	3, --  4 man     100
	3, --  5 axe     100
	3, --  6 hummer  100
	4, --  7 bear    75
	4, --  8 wolf    75
	4, --  9 boar    75
	5, -- 10 fox     75
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0,}, --  3 woman
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0,}, --  4 man
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0,}, --  5 axe
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0,}, --  6 hummer
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  7 bear
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  8 wolf
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, --  9 boar
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 10 fox
}

local symset2 = {
	1, --  1 wild
	1, --  2 scatter
	3, --  3 woman
	4, --  4 man
	4, --  5 axe
	4, --  6 hummer
	4, --  7 bear
	4, --  8 wolf
	4, --  9 boar
	4, -- 10 fox
}

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	6, --  3 woman
	6, --  4 man
	6, --  5 axe
	6, --  6 hummer
	6, --  7 bear
	6, --  8 wolf
	6, --  9 boar
	6, -- 10 fox
}

math.randomseed(os.time())

local function batchreel(comment)
	local reel1, iter1 = makereel(symset1, neighbours)
	local reel2, iter2 = makereelhot(symset2, 3, {[2]=true}, chunklen, true)
	print(comment)
	if iter1 >= 1000 then
		print(string.format("iterations: %d, %d", iter1, iter2))
	end
	printreel(tableglue(reel1, reel2))
end

do
	local n11, n21 = symset1[1], symset2[1]
	symset1[1], symset2[1] = 0, 0
	batchreel "reel 1, 5"
	symset1[1], symset2[1] = n11, n21
end

do
	batchreel "reel 2, 3, 4"
end
