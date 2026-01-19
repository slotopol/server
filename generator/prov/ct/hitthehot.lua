local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	4, --  1 wild      (2, 3, 4 reels only)
	2, --  2 scatter
	4, --  3 clover    1500
	5, --  4 horseshoe 1250
	5, --  5 bell      500
	6, --  6 apple     100+
	6, --  7 orange    100
	6, --  8 lemon     100
	6, --  9 plum      100
	6, -- 10 cherry    100
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	5, --  3 clover
	5, --  4 horseshoe
	3, --  5 bell
	3, --  6 apple
	3, --  7 orange
	3, --  8 lemon
	3, --  9 plum
	3, -- 10 cherry
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
