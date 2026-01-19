local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	2, --  3 seven
	3, --  4 strawberry
	4, --  5 blueberry
	6, --  6 pear
	7, --  7 plum
	7, --  8 peach
	7, --  9 papaya
	7, -- 10 cherry
}

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	1, --  3 seven
	1, --  4 strawberry
	1, --  5 blueberry
	3, --  6 pear
	3, --  7 plum
	3, --  8 peach
	3, --  9 papaya
	3, -- 10 cherry
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
