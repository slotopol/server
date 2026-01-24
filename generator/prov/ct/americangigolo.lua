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

local chained = {
	0, --  1 wild
	0, --  2 scatter
	4, --  3 blonde
	4, --  4 brunette
	4, --  5 cat
	4, --  6 dog
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 then
		ss[1] = 0
	end
	return makereelct(ss, 3, {[1]=true, [2]=true}, chained)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4, 5"
printreel(reelgen(2))
