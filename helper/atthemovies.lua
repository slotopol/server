local path = arg[0]:match("(.*[/\\])")
dofile(path.."reelgen.lua")

local symset = {
	1, --  1 oscar
	2, --  2 popcorn
	3, --  3 poster
	4, --  4 a
	4, --  5 dummy
	4, --  6 maw
	5, --  7 starship
	5, --  8 heart
	1, --  9 masks
	1, -- 10 projector
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, }, --  1
	{ 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, }, --  2
	{ 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, }, --  3
	{ 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, }, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, }, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, }, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, }, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, }, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 1, }, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, }, -- 10
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
