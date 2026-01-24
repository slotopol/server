local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	5, --  1 wild    (2, 3, 4, 5 reels only)
	2, --  2 scatter
	7, --  3 aryan   1000
	7, --  4 nymph   200
	7, --  5 pot     150
	7, --  6 vase    150
	8, --  7 ace     100
	8, --  8 king    100
	8, --  9 queen   100
	8, -- 10 jack    100
}

local chained = {
	0, --  1 wild
	0, --  2 scatter
	4, --  3 aryan
	4, --  4 nymph
	4, --  5 pot
	4, --  6 vase
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
