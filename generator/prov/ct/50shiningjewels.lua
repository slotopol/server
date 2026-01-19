local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	4, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 seven   1000
	5, --  4 heart   300
	6, --  5 diamond 200
	6, --  6 emerald 200
	7, --  7 gold    100
	7, --  8 blue    100
	8, --  9 yellow  100
	8, -- 10 lilac   100
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	6, --  3 seven
	6, --  4 heart
	6, --  5 diamond
	6, --  6 emerald
	8, --  7 gold
	8, --  8 blue
	8, --  9 yellow
	8, -- 10 lilac
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	if n == 1 or n == 3 or n == 5 then
		ss[2] = ss[2] + 1
	end
	return makereelhot(ss, 3, {[2]=true}, chunklen)
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
