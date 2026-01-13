local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	5, --  1 wild (2, 3, 4, 5 reels only)
	2, --  2 scatter
	8, --  3 blonde   50/150/500
	8, --  4 brunette 15/50/100
	8, --  5 cat      10/50/100
	8, --  6 dog      10/50/100
	8, --  7 ace      10/20/100
	9, --  8 king     10/20/100
	9, --  9 queen    10/20/100
	9, -- 10 jack     10/20/100
}

local bigsym = {
	0, --  1 wild
	0, --  2 scatter
	1, --  3 blonde
	1, --  4 brunette
	1, --  5 cat
	1, --  6 dog
	1, --  7 ace
	1, --  8 king
	1, --  9 queen
	1, -- 10 jack
}

local function reelgen(n)
	local function make()
		return makereelct(symset, 3, {[1]=true, [2]=true}, 4, bigsym)
	end
	if n == 1 then
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
print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4, 5"
printreel(reelgen(2))
