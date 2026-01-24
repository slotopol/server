local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	5, --  1 wild (2, 3, 4, 5 reels only)
	2, --  2 scatter
	8, --  3 woman 50/200/1000
	8, --  4 ox    15/50/100
	8, --  5 dog   10/50/100
	8, --  6 fox   10/50/100
	8, --  7 ace   10/20/100
	8, --  8 king  10/20/100
	8, --  9 queen 10/20/100
	9, -- 10 jack  10/20/100
}

local chained = {
	0, --  1 wild
	0, --  2 scatter
	0, --  3 woman
	0, --  4 ox
	0, --  5 dog
	0, --  6 fox
	5, --  7 ace
	5, --  8 king
	5, --  9 queen
	5, -- 10 jack
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
