local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, -- 1 diamond 1000
	4, -- 2 seven   80
	3, -- 3 bar3    40
	5, -- 4 bar2    25
	8, -- 5 bar1    10
	3, -- 6 cherry  10
}

local neighbours = {
	--1, 2, 3, 4, 5, 6,
	{ 1, 0, 0, 0, 0, 1,}, -- 1 diamond
	{ 0, 0, 0, 0, 0, 0,}, -- 2 seven
	{ 0, 0, 0, 0, 0, 0,}, -- 3 bar3
	{ 0, 0, 0, 0, 0, 0,}, -- 4 bar2
	{ 0, 0, 0, 0, 0, 0,}, -- 5 bar1
	{ 1, 0, 0, 0, 0, 1,}, -- 6 cherry
}

local function insertspace(reel)
	local reelsp = {}
	local prev
	for i = 1, rawlen(reel) do
		local sym = rawget(reel, i)
		if prev then
			table.insert(reelsp, 0)
			if prev == 1 or sym == 1 then
				table.insert(reelsp, 0)
			end
		end
		table.insert(reelsp, sym)
		prev = sym
	end
	table.insert(reelsp, 0)
	if prev == 1 or rawget(reel, 1) == 1 then
		table.insert(reelsp, 0)
	end
	return reelsp
end

math.randomseed(os.time())
local reel, iter = makereel(symset, neighbours)
reel = insertspace(reel)
printreel(reel, iter)
