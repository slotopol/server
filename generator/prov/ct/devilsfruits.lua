local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 wild   1500
	4, --  2 seven  100
	4, --  3 pike   35
	4, --  4 bell   25
	4, --  5 orange 25
	4, --  6 plum   25
	5, --  7 bar3   25
	5, --  8 bar2   20
	5, --  9 bar1   15
	7, -- 10 cherry 10
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,}, --  1 wild
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0,}, --  2 seven
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,}, --  3 pike
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  4 bell
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  5 orange
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  6 plum
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  7 bar3
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  8 bar2
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,}, --  9 bar1
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, -- 10 cherry
}

local function insertspace(reel)
	local reelsp = {}
	local n = 0
	for i = 1, rawlen(reel) do
		table.insert(reelsp, rawget(reel, i))
		n = n + 1
		if n >= 3 then
			table.insert(reelsp, 0)
			n = 0
		end
	end
	if reelsp[#reelsp] ~= 0 then
		table.insert(reelsp, 0)
	end
	return reelsp
end

local function reelgen()
	local reel, iter = makereel(symset, neighbours)
	reel = insertspace(reel)
	return reel, iter
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen())
