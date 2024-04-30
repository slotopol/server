local path = arg[0]:match("(.*[/\\])")
dofile(path.."reelgen.lua")

local symset = {
	2, --  1 troll1
	2, --  2 troll2
	2, --  3 troll3
	2, --  4 troll4
	3, --  5 troll5
	3, --  6 troll6
	3, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
	4, -- 11 ten
	1, -- 12 wild
	0, -- 13 golden
	1, -- 14 scatter
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,14
	{ 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1
	{ 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2
	{ 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3
	{ 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 11
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 12
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 13
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 14
}

math.randomseed(os.time())
local reel = MakeReel(symset)
print("reel length: " .. #reel)
ShuffleReel(reel)
local iter = CorrectReel(reel, neighbours)
if iter > 1 then
	if iter >= 1000 then
		print"too many neighbours shuffle iterations"
	else
		print(iter.." iterations")
	end
end
RrintReel(reel)