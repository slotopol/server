local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 wild    (2, 3, 4 reels only)
	2, --  2 scatter
	7, --  3 seven   750
	7, --  4 grape   200
	7, --  5 melon   200
	7, --  6 apple   200
	8, --  7 orange  100
	8, --  8 lemon   100
	8, --  9 plum    100
	8, -- 10 cherry  100
}

local bigsym = {
	0, --  1 wild
	0, --  2 scatter
	1, --  3 seven
	1, --  4 grape
	1, --  5 melon
	1, --  6 apple
	1, --  7 orange
	1, --  8 lemon
	1, --  9 plum
	1, -- 10 cherry
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
		ss[2] = ss[2]-1
	end
	return makereelct(ss, 3, {[1]=true, [2]=true}, 4, bigsym)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
