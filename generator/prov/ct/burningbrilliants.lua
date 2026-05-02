local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	4, --  1 wild      (2, 3, 4 reel)
	1, --  2 scatter
	4, --  3 star      35/100/2000
	4, --  4 ruby      10/25/100
	4, --  5 emerald   10/25/100
	4, --  6 topaz     10/25/100
	5, --  7 spader    7/10/100
	5, --  8 heart     7/10/100
	5, --  9 diamond   5/10/100
	5, -- 10 club      5/10/100
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	6, --  3 star
	6, --  4 ruby
	6, --  5 emerald
	6, --  6 topaz
	6, --  7 spader
	6, --  8 heart
	6, --  9 diamond
	6, -- 10 club
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	return makereelhot(ss, 4, {[2]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
