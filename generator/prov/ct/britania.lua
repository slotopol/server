local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	6, --  1 wild (2, 3, 4 reels only)
	2, --  2 scatter
	7, --  3 blue
	7, --  4 red
	7, --  5 swords
	7, --  6 axe
	7, --  7 ace
	7, --  8 king
	8, --  9 queen
	8, -- 10 jack
}

local bigsym = {
	1, --  1 wild (2, 3, 4 reels only)
	0, --  2 scatter
	1, --  3 blue
	1, --  4 red
	1, --  5 swords
	1, --  6 axe
	1, --  7 ace
	1, --  8 king
	1, --  9 queen
	1, -- 10 jack
}

local function reelgen(n)
	local function make()
		return makereelct(symset, 3, {[1]=true, [2]=true}, 4, bigsym)
	end
	if n == 1 or n == 5 then
		local n1 = symset[1]
		symset[1] = 0
		bigsym[1] = 0
		local reel, iter = make()
		bigsym[1] = 1
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
