local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 bee   5000
	3, --  2 snail 1000
	4, --  3 fly   500
	4, --  4 worm  250
	4, --  5 ace   200
	4, --  6 king  200
	6, --  7 queen 100
	6, --  8 jack  100
	6, --  9 ten   100
	2, -- 10 note
	3, -- 11 jazzbee (3 reel only)
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, --  1 bee
	{ 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2,}, --  2 snail
	{ 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 2,}, --  3 fly
	{ 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 2,}, --  4 worm
	{ 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,}, --  5 ace
	{ 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0,}, --  6 king
	{ 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0,}, --  7 queen
	{ 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,}, --  8 jack
	{ 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0,}, --  9 ten
	{ 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2,}, -- 10 note
	{ 2, 2, 2, 2, 0, 0, 0, 0, 0, 2, 2,}, -- 11 jazzbee
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n ~= 3 then
		ss[11] = 0
	else
		ss[10] = ss[10]+1
	end
	return makereel(ss, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 2, 4, 5"
printreel(reelgen(1))
print "reel 3"
printreel(reelgen(3))
