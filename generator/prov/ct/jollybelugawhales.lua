local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild    (2, 3, 4 reel)
	2, --  2 bear    40/500/10000
	4, --  3 wolf    20/80/200
	4, --  4 owl     20/80/200
	4, --  5 walrus  15/40/100
	4, --  6 puffin  15/40/100
	4, --  7 ace     10/20/100
	5, --  8 king    10/20/100
	5, --  9 queen   10/20/100
	5, -- 10 jack    10/20/100
}

local chunklen = {
	1, --  1 wild
	1, --  2 bear
	3, --  3 wolf
	3, --  4 owl
	3, --  5 walrus
	3, --  6 puffin
	3, --  7 ace
	3, --  8 king
	3, --  9 queen
	3, -- 10 jack
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	return makereelhot(ss, 3, {[1]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
