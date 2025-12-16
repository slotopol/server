local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild      9000
	1, --  2 scatter
	4, --  3 dollar    700
	4, --  4 horseshoe 700
	3, --  5 bug       200
	3, --  6 bell      200
	4, --  7 jug       200
	3, --  8 ace       150
	3, --  9 king      150
	3, -- 10 queen     100
	3, -- 11 jack      100
	3, -- 12 ten       100
	3, -- 13 nine      100
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  1 wild
	{ 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  2 scatter
	{ 1, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  3 dollar
	{ 1, 1, 1, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0,}, --  4 horseshoe
	{ 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 0, 0,}, --  5 bug
	{ 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0,}, --  6 bell
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  7 jug
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  8 ace
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  9 king
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, -- 10 queen
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, -- 11 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, -- 12 ten
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 13 nine
}

local function reelgen()
	return makereel(symset, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
