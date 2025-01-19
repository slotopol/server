local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 wild
	1, --  2 scatter
	3, --  3 money bag
	3, --  4 diamonds
	3, --  5 robbery
	3, --  6 picture
	4, --  7 watch
	4, --  8 cop
	4, --  9 jail
	4, -- 10 thief
	4, -- 11 handcuffs
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1
	{ 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  3
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  4
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  5
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  6
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  9
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 10
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 11
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))

local fsn = {10, 15, 15, 20, 25}
local sum = 0
for _, v in ipairs(fsn) do
	sum = sum + v
end
print("free spins Efs = "..sum/#fsn..", len = "..#fsn)
