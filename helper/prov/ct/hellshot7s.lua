local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset = {
	1, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 bar    1000
	4, --  4 grapes 500
	4, --  5 apple  200
	4, --  6 pear   150
	4, --  7 plum   50
	4, --  8 orange 50
	4, --  9 lemon  50
	4, -- 10 cherry 50
}

local chunklen = {
	1, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	1, --  3 bar
	1, --  4 grapes
	1, --  5 apple
	1, --  6 pear
	4, --  7 plum
	4, --  8 orange
	4, --  9 lemon
	4, -- 10 cherry
}

math.randomseed(os.time())

local function batchreel()
	printreel(makereelhot(symset, 3, {[1]=true, [2]=true}, chunklen, true))
end

do
	print "reel 1, 5"
	local n1 = symset[1]
	symset[1] = 0
	batchreel()
	symset[1] = n1
end

do
	print "reel 2, 3, 4"
	batchreel()
end
