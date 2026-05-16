local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	5, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 man    30/100/1000
	4, --  4 banana 5/20/400
	5, --  5 apple  5/20/100
	5, --  6 melon  5/20/100
	5, --  7 orange 5/10/100
	5, --  8 lemon  5/10/100
	5, --  9 plum   5/10/100
	5, -- 10 cherry 5/10/100
}

local chunklen = {
	4, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	3, --  3 man
	3, --  4 banana
	6, --  5 apple
	6, --  6 melon
	6, --  7 orange
	6, --  8 lemon
	6, --  9 plum
	6, -- 10 cherry
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	return makereelhot(ss, 3, {[2]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
