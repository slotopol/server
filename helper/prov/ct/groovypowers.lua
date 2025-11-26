local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	2, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 glasses 1000
	4, --  4 blonde  200
	4, --  5 curly   200
	4, --  6 bald    200
	5, --  7 ace     100
	5, --  8 king    100
	5, --  9 queen   100
	5, -- 10 jack    100
}

local chunklen = {
	1, --  1 wild    (2, 3, 4 reels only)
	1, --  2 scatter
	6, --  3 glasses
	6, --  4 blonde
	6, --  5 curly
	6, --  6 bald
	6, --  7 ace
	6, --  8 king
	6, --  9 queen
	6, -- 10 jack
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		return makereelhot(symset, 3, {[2]=true}, chunklen, true)
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

print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
