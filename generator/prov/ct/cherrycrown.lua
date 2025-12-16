local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset24 = {
	0, --  1 wild    (1, 3 reels only)
	0, --  2 scatter (1, 3, 5 reels only)
	3, --  3 seven   500/60/20
	4, --  4 melon   100/30/20
	4, --  5 apple   100/30/20
	5, --  6 pear    100/30/15
	6, --  7 orange  100/20/15
	6, --  8 lemon   100/20/15
	6, --  9 plum    50/20/10
	6, -- 10 cherry  50/20/10
}

local symset13 = tcopy(symset24)
symset13[2] = 2 -- scatter
for i, v in ipairs(symset13) do
	symset13[i] = v * 3
end
symset13[1] = 1 -- only 1 wild on the reel

local symset5 = tcopy(symset24)
symset5[2] = 2 -- scatter

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	1, --  3 seven
	4, --  4 melon
	4, --  5 apple
	4, --  6 pear
	4, --  7 orange
	4, --  8 lemon
	4, --  9 plum
	4, -- 10 cherry
}

local function reelgen(n)
	if n == 1 or n == 3 then
		return makereelhot(symset13, 3, {[2]=true}, chunklen)
	elseif n == 2 or n == 4 then
		return makereelhot(symset24, 3, {[2]=true}, chunklen)
	else -- n == 5
		return makereelhot(symset5, 3, {[2]=true}, chunklen)
	end
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
