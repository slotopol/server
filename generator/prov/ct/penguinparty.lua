local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild    (2, 3, 4 reel)
	3, --  2 seal    40/300/1000
	4, --  3 shark   20/80/200
	4, --  4 dolphin 20/80/200
	4, --  5 medusa  15/40/200
	4, --  6 tuna    15/40/200
	4, --  7 ace     10/20/100
	4, --  8 king    10/20/100
	4, --  9 queen   10/20/100
	5, -- 10 jack    10/20/100
}

local chunklen = {
	1, --  1 wild
	1, --  2 seal
	3, --  3 shark
	3, --  4 dolphin
	3, --  5 medusa
	3, --  6 tuna
	3, --  7 ace
	3, --  8 king
	3, --  9 queen
	3, -- 10 jack
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	return makereelhot(ss, 3, {[1]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
