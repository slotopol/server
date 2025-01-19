local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	3, --  1 booth
	3, --  2 vip
	3, --  3 food
	3, --  4 bell
	4, --  5 ace
	4, --  6 king
	4, --  7 queen
	4, --  8 jack
	1, --  9 bonus
	1, -- 10 wild
	1, -- 11 scatter
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3
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

local dd = {10,20,20,40,40,80,120,120,450}
local sum = 0
for _, v in ipairs(dd) do
	sum = sum + v
end
print("bonus M = "..sum/#dd..", E = 3*M = "..3*sum/#dd)
