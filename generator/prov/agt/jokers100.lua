local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	4, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	5, --  3 strawberry
	5, --  4 pear
	5, --  5 greenstar
	5, --  6 redstar
	6, --  7 plum
	6, --  8 peach
	6, --  9 papaya
	6, -- 10 cherry
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	8, --  3 strawberry
	8, --  4 pear
	8, --  5 greenstar
	8, --  6 redstar
	8, --  7 plum
	8, --  8 peach
	8, --  9 papaya
	8, -- 10 cherry
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
