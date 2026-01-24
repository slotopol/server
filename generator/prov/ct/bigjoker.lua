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

local chained = {
	0, --  1 wild
	0, --  2 scatter
	4, --  3 seven
	4, --  4 grape
	4, --  5 melon
	4, --  6 apple
	4, --  7 orange
	4, --  8 lemon
	4, --  9 plum
	4, -- 10 cherry
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
		ss[2] = ss[2]-1
	end
	return makereelct(ss, 3, {[1]=true, [2]=true}, chained)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
