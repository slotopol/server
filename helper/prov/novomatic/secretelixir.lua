local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 lover     1000
	2, --  2 wife      200
	2, --  3 husband   200
	2, --  4 owl       80
	2, --  5 gargoyle1 60
	2, --  6 gargoyle2 60
	3, --  7 ace       40
	3, --  8 king      40
	4, --  9 queen     20
	4, -- 10 jack      20
	4, -- 11 ten       20
	2, -- 12 scatter
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,
	{ 2, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  1 lover
	{ 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  2 wife
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0,}, --  3 husband
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0,}, --  4 owl
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0,}, --  5 gargoyle1
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0,}, --  6 gargoyle2
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  7 ace
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  8 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  9 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 10 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 11 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 12 scatter
}

math.randomseed(os.time())
printreel(makereel(symset, neighbours))

local s
--local r5reg = {1, 2, 1, 1, 3, 1, 1, 2, 1, 1, 5, 1, 1, 5, 1} -- Er5reg = 1.8
local r5reg = {1, 5, 1, 1, 3, 1, 1, 5, 1, 1, 2, 1, 1, 5, 1} -- Er5reg = 2
s = 0
for _, m in pairs(r5reg) do
	s = s + m
end
print("Er5reg = "..s/#r5reg)

local r5bon = {2, 5, 3, 5, 3, 3, 2, 5, 20, 3, 2, 5, 10, 5, 2} -- Er5bon = 5
s = 0
for _, m in pairs(r5bon) do
	s = s + m
end
print("Er5bon = "..s/#r5bon)
