local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	9, --  1 wild    (2, 3, 4, 5 reels only)
	5, --  2 scatter (2, 3, 4 reels only)
	8, --  3 witch   300
	8, --  4 cat     100
	6, --  5 banana  100
	6, --  6 grape   100
	6, --  7 apple   50
	6, --  8 melon   50
	6, --  9 orange  30
	6, -- 10 lemon   30
	6, -- 11 plum    30
	6, -- 12 cherry  30
}

local chained = {
	4, --  1 wild
	0, --  2 scatter
	6, --  3 witch
	6, --  4 cat
	4, --  5 banana
	4, --  6 grape
	4, --  7 apple
	4, --  8 melon
	4, --  9 orange
	4, -- 10 lemon
	4, -- 11 plum
	4, -- 12 cherry
}

local function reelgen(n)
	local ss, cs = tcopy(symset), tcopy(chained)
	if n == 1 then
		ss[1], cs[1] = 0, 0
	end
	if n == 1 or n == 5 then
		ss[2] = 0
	end
	return makereelct(ss, 3, {[2]=true}, cs)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
print "reel 5"
printreel(reelgen(5))
