local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 wild    (2, 3, 4, 5 reels only)
	2, --  2 scatter
	6, --  3 samurai 2000
	6, --  4 geisha  400
	6, --  5 bowl    300
	6, --  6 coins   300
	6, --  7 ace     200
	6, --  8 king    200
	6, --  9 queen   200
	6, -- 10 jack    100
}

local chained = {
	0, --  1 wild
	0, --  2 scatter
	4, --  3 samurai
	4, --  4 geisha
	4, --  5 bowl
	4, --  6 coins
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 then
		ss[1] = 0
	end
	return makereelct(ss, 3, {[1]=true, [2]=true}, chained)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4, 5"
printreel(reelgen(2))
