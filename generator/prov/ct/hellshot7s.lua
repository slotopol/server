local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 bar    1000
	4, --  4 grapes 500
	4, --  5 apple  200
	4, --  6 pear   150
	4, --  7 plum   50
	4, --  8 orange 50
	5, --  9 lemon  50
	5, -- 10 cherry 50
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

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	return makereelhot(ss, 3, {[1]=true, [2]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
