local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	2, --  1 wild    (2, 3, 4 reels only)
	2, --  2 scatter
	6, --  3 woman   1000
	7, --  4 man     100
	7, --  5 axe     100
	7, --  6 hummer  100
	8, --  7 bear    75
	8, --  8 wolf    75
	8, --  9 boar    75
	9, -- 10 fox     75
}

local chained = {
	0, --  1 wild
	0, --  2 scatter
	3, --  3 woman
	4, --  4 man
	4, --  5 axe
	4, --  6 hummer
	4, --  7 bear
	4, --  8 wolf
	4, --  9 boar
	4, -- 10 fox
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	return makereelct(ss, 3, {[1]=true, [2]=true}, chained)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
