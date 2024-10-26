local path = arg[0]:match("(.*[/\\])")
dofile(path.."reelgen.lua")

local symset = {
	1, --  1 wild (on 2, 3, 4 reels)
	2, --  2 scatter1 (on all reels)
	2, --  3 scatter2 (on 1, 3, 5 reels)
	2, --  4 seven
	4, --  5 grape
	4, --  6 watermelon
	6, --  7 avocado
	7, --  8 pomegranate
	7, --  9 carambola
	7, -- 10 maracuya
	7, -- 11 orange
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1 wild (on 2, 3, 4 reels)
	{ 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 scatter1 (on all reels)
	{ 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3 scatter2 (on 1, 3, 5 reels)
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4 seven
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  5 grape
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  6 watermelon
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  7 avocado
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  8 pomegranate
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  9 carambola
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 10 maracuya
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 11 orange
}

math.randomseed(os.time())
local reel = MakeReel(symset)
print("reel length: " .. #reel)
ShuffleReel(reel)
local iter = CorrectReel(reel, neighbours)
PrintReel(reel, iter)
