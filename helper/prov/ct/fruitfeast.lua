local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	4, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 banana  1000
	4, --  4 grape   300
	6, --  5 melon   200
	6, --  6 apple   200
	8, --  7 orange  100
	8, --  8 lemon   100
	8, --  9 plum    100
	8, -- 10 cherry  100
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	4, --  3 banana
	4, --  4 grape
	4, --  5 melon
	4, --  6 apple
	4, --  7 orange
	4, --  8 lemon
	4, --  9 plum
	4, -- 10 cherry
}

math.randomseed(os.time())

do
	print "reel 1, 5"
	local n1 = symset[1]
	symset[1] = 0
	symset[2] = symset[2] + 1
	printreel(makereelhot(symset, 4, {[2]=true}, chunklen, true))
	symset[1] = n1
	symset[2] = symset[2] - 1
end

do
	print "reel 2, 4"
	printreel(makereelhot(symset, 4, {[2]=true}, chunklen, true))
end

do
	print "reel 3"
	symset[2] = symset[2] + 1
	printreel(makereelhot(symset, 4, {[2]=true}, chunklen, true))
	symset[2] = symset[2] - 1
end
