local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 two balls
	1, --  2 white ball
	1, --  3 yellow ball
	4, --  4 electrocar
	6, --  5 golf clubs
	6, --  6 flag
	6, --  7 beer
	6, --  8 slippers
	1, --  9 fitch
	1, -- 10 drake
	1, -- 11 luce
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1
	{ 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2
	{ 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2,}, -- 11
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))

-- bonus game calculation
local golf = {
  4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20, 25, 40, 50, 100,
}
local sum = 0
for _, v in ipairs(golf) do
  sum = sum + v
end
print("golf bonus: "..#golf.." elements, sum = "..sum..", E = "..sum/#golf)
