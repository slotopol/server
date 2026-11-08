local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	4, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 seven     25/50/1000
	4, --  4 coin      5/15/100
	4, --  5 bell      5/15/100
	4, --  6 horseshoe 5/15/100
	5, --  7 apple     5/10/50
	5, --  8 lemon     5/10/50
	6, --  9 plum      5/10/50
	6, -- 10 cherry    5/10/50
}

local chunklen = {
	4, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 seven
	4, --  4 coin
	4, --  5 bell
	4, --  6 horseshoe
	6, --  7 apple
	6, --  8 lemon
	6, --  9 plum
	6, -- 10 cherry
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
printreel(reelgen(1))
printreel(reelgen(2))
printreel(reelgen(3))
printreel(reelgen(4))
printreel(reelgen(5))
