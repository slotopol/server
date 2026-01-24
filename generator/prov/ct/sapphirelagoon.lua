local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	6, --  1 wild (2, 3, 4, 5 reels only)
	4, --  2 scatter
	6, --  3 man
	7, --  4 woman
	7, --  5 flask
	7, --  6 hook
	7, --  7 ace
	7, --  8 king
	7, --  9 queen
	7, -- 10 jack
}

local chained = {
	4, --  1 wild
	3, --  2 scatter
	4, --  3 man
	4, --  4 woman
	4, --  5 flask
	4, --  6 hook
	4, --  7 ace
	4, --  8 king
	4, --  9 queen
	4, -- 10 jack
}

local function reelgen(n)
	local ss, cs = tcopy(symset), tcopy(chained)
	if n == 1 then
		ss[1], cs[1] = 0, 0
	end
	return makereelct(ss, 3, {[1]=true, [2]=true}, cs)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1"
printreel(reelgen(1))
print "reel 2, 3, 4, 5"
printreel(reelgen(2))
