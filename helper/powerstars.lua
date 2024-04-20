local path = arg[0]:match("(.*[/\\])")
dofile(path.."reelgen.lua")

local symset = {
	2, --  1 seven
	3, --  2 bell
	4, --  3 melon
	4, --  4 grapes
	5, --  5 plum
	5, --  6 orange
	5, --  7 lemon
	5, --  8 cherry
	0, --  9 star
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0 }, --  1
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0 }, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0 }, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0 }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0 }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0 }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0 }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0 }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2 }, --  9
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
