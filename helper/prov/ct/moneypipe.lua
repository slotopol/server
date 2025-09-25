local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset24 = {
	0, --  1 wild    (only on reel 1, 3)
	0, --  2 scatter (only on reel 1, 3, 5)
	4, --  3 dollar  750/70/20
	4, --  4 pliers  100/30/20
	4, --  5 hammer  100/30/20
	4, --  6 worker  100/30/15
	6, --  7 ace     100/20/15
	6, --  8 king    100/20/15
	6, --  9 queen   50/20/10
	6, -- 10 jack    50/20/10
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
	1, --  3 dollar
	4, --  4 pliers
	4, --  5 hammer
	4, --  6 worker
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

math.randomseed(os.time())
printreel(makereelhot(symset13, 4, {[2]=true}, chunklen, true))
printreel(makereelhot(symset24, 4, {[2]=true}, chunklen, true))
printreel(makereelhot(symset13, 4, {[2]=true}, chunklen, true))
printreel(makereelhot(symset24, 4, {[2]=true}, chunklen, true))
printreel(makereelhot(symset5, 4, {[2]=true}, chunklen, true))
