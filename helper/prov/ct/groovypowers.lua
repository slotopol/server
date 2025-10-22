local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	2, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 glasses 1000
	4, --  4 blonde  200
	4, --  5 curly   200
	4, --  6 bald    200
	5, --  7 ace     100
	5, --  8 king    100
	5, --  9 queen   100
	5, -- 10 jack    100
}

local chunklen = {
	1, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	6, --  3 glasses
	6, --  4 blonde
	6, --  5 curly
	6, --  6 bald
	6, --  7 ace
	6, --  8 king
	6, --  9 queen
	6, -- 10 jack
}

math.randomseed(os.time())

do
	print "reel 1, 5"
	local n1 = symset[1]
	symset[1] = 0
	printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
	symset[1] = n1
end

do
	print "reel 2, 3, 4"
	printreel(makereelhot(symset, 3, {[2]=true}, chunklen, true))
end
