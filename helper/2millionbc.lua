local path = arg[0]:match("(.*[/\\])")
dofile(path.."reelgen.lua")

local symset = {
	2, --  1 girl
	2, --  2 lion
	3, --  3 bee
	3, --  4 stone
	3, --  5 wheel
	3, --  6 club
	3, --  7 chaplet
	3, --  8 gold
	3, --  9 vase
	3, -- 10 ruby
	1, -- 11 fire
	0, -- 12 acorn
	3, -- 13 diamond
}

local symset5 = {
	8, --  1 girl
	8, --  2 lion
	9, --  3 bee
	10, --  4 stone
	10, --  5 wheel
	10, --  6 club
	10, --  7 chaplet
	10, --  8 gold
	10, --  9 vase
	10, -- 10 ruby
	4, -- 11 fire
	1, -- 12 acorn
	10, -- 13 diamond
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, --  1
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  5
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  6
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 11
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 12
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 13
}

math.randomseed(os.time())

local reel = MakeReel(symset)
print("reel length: " .. #reel)
ShuffleReel(reel)
local iter = CorrectReel(reel, neighbours)
RrintReel(reel, iter)

local reel5 = MakeReel(symset5)
print("5-reel length: " .. #reel5)
ShuffleReel(reel5)
local iter5 = CorrectReel(reel5, neighbours)
RrintReel(reel5, iter5)
