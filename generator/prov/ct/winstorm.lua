local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	4, --  1 wild      (2, 3, 4 reel)
	1, --  2 scatter
	4, --  3 seven     1000
	4, --  4 coin      100+
	4, --  5 horseshoe 100+
	4, --  6 bell      100+
	4, --  7 ace       100
	4, --  8 king      100
	4, --  9 queen     100
	4, -- 10 jack      100
}

local chunklen = {
	4, --  1 wild
	1, --  2 scatter
	3, --  3 seven
	6, --  4 coin
	6, --  5 horseshoe
	6, --  6 bell
	6, --  7 ace
	6, --  8 king
	6, --  9 queen
	6, -- 10 jack
}

local function reelgen(n)
	local function make()
		return makereelhot(symset, 3, {[2]=true}, chunklen)
	end
	if n == 1 or n == 5 then
		local n1 = symset[1]
		symset[1] = 0
		local reel, iter = make()
		symset[1] = n1
		return reel, iter
	else
		return make()
	end
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
