local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	3, --  1 wild    (2, 3, 4 reels only)
	2, --  2 scatter
	4, --  3 seven   2500
	5, --  4 grape   400
	5, --  5 melon   400
	6, --  6 apple   300
	7, --  7 orange  100
	7, --  8 lemon   100
	7, --  9 plum    100
	7, -- 10 cherry  100
}

local chunklen = {
	3, --  1 wild
	1, --  2 scatter
	4, --  3 seven
	5, --  4 grape
	5, --  5 melon
	3, --  6 apple
	4, --  7 orange
	4, --  8 lemon
	4, --  9 plum
	4, -- 10 cherry
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		return makereelhot(symset, 3, {[2]=true}, chunklen, true)
	end
	if n == 1 or n == 5 then
		local n1 = symset[1]
		symset[1] = 0
		local reel, iter = make()
		symset[1] = n1
		return reel, iter
	else
		return make()
	end
end

if autoscan then
	return reelgen
end

print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
