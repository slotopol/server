local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	6, --  1 wild (2, 3, 4 reels only)
	2, --  2 scatter
	7, --  3 blue
	7, --  4 red
	7, --  5 swords
	7, --  6 axe
	7, --  7 ace
	7, --  8 king
	8, --  9 queen
	8, -- 10 jack
}

local chained = {
	4, --  1 wild (2, 3, 4 reels only)
	0, --  2 scatter
	3, --  3 blue
	3, --  4 red
	3, --  5 swords
	3, --  6 axe
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local function reelgen(n)
	local ss, cs = tcopy(symset), tcopy(chained)
	if n == 1 or n == 5 then
		ss[1], cs[1] = 0, 0
	end
	return makereelct(ss, 3, {[1]=true, [2]=true}, cs)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
