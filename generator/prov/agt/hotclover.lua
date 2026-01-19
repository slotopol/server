local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	 1, --  1 wild (2, 3, 4 reels only)
	 2, --  2 scatter
	 4, --  3 seven
	 8, --  4 strawberr
	 9, --  5 grapes
	 8, --  6 bar
	10, --  7 plum
	10, --  8 orange
	10, --  9 lemon
	10, -- 10 cherry
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	return makereelhot(ss, 4, {[1]=true, [2]=true}, {})
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
