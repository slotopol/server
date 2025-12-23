local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset135 = {
	1, -- 1 diamond   5000
	1, -- 2 star
	2, -- 3 topaz     500
	9, -- 4 sapphire  200
	2, -- 5 heliodor  200
	9, -- 6 ruby      200
	2, -- 7 tanzanite 200
	9, -- 8 emerald   200
}

local symset24 = {
	1, -- 1 diamond   5000
	1, -- 2 star
	9, -- 3 topaz     500
	2, -- 4 sapphire  200
	9, -- 5 heliodor  200
	2, -- 6 ruby      200
	9, -- 7 tanzanite 200
	2, -- 8 emerald   200
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8
	{ 2, 2, 0, 0, 0, 0, 0, 0 }, -- 1 diamond
	{ 2, 2, 0, 0, 0, 0, 0, 0 }, -- 2 star
	{ 0, 0, 2, 0, 0, 0, 0, 0 }, -- 3 topaz
	{ 0, 0, 0, 2, 0, 0, 0, 0 }, -- 4 sapphire
	{ 0, 0, 0, 0, 2, 0, 0, 0 }, -- 5 heliodor
	{ 0, 0, 0, 0, 0, 2, 0, 0 }, -- 6 ruby
	{ 0, 0, 0, 0, 0, 0, 2, 0 }, -- 7 tanzanite
	{ 0, 0, 0, 0, 0, 0, 0, 2 }, -- 8 emerald
}

local function reelgen(n)
	if n == 2 or n == 4 then
		return makereel(symset24, neighbours)
	else
		return makereel(symset135, neighbours)
	end
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 3, 5"
printreel(reelgen(1))
print "reel 2, 4"
printreel(reelgen(2))
