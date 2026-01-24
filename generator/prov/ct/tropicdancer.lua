local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	6, --  1 wild (2, 3, 4, 5 reels only)
	5, --  2 scatter
	6, --  3 singer
	6, --  4 dancer man
	6, --  5 dancer girl 1
	6, --  6 dancer girl 2
	8, --  7 ace
	8, --  8 king
	8, --  9 queen
	8, -- 10 jack
}

local chained = {
	3, --  1 wild
	3, --  2 scatter
	3, --  3 singer
	3, --  4 dancer man
	3, --  5 dancer girl 1
	3, --  6 dancer girl 2
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local function reelgen(n)
	local ss, cs = tcopy(symset), tcopy(chained)
	if n == 1 then
		ss[1], cs[1] = 0, 0
	end
	return makereelct(ss, 3, {[1]=true, [2]=true}, cs)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4, 5"
printreel(reelgen(2))
