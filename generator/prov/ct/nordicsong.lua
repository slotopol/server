local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild    (2, 3, 4, 5 reels only)
	2, --  2 scatter (1, 3, 5 reels only)
	5, --  3 man     1000
	4, --  4 woman   500
	3, --  5 owl     400
	3, --  6 dog     400
	4, --  7 ace     200
	4, --  8 king    200
	4, --  9 queen   100
	4, -- 10 jack    100
	4, -- 11 ten     100
}

local chained = {
	0, --  1 wild    (2, 3, 4, 5 reels only)
	0, --  2 scatter (1, 3, 5 reels only)
	4, --  3 man
	0, --  4 woman
	0, --  5 owl
	0, --  6 dog
	0, --  7 ace
	0, --  8 king
	0, --  9 queen
	0, -- 10 jack
	0, -- 11 ten
}

local function reelgen(n)
	local ss, cs = tcopy(symset), tcopy(chained)
	if n == 1 then
		ss[1], cs[1] = 0, 0
	end
	if n == 2 or n == 4 then
		ss[2] = 0
	end
	return makereelct(ss, 3, {[2]=true}, cs)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1"
printreel(reelgen(1))
print "reel 2, 4"
printreel(reelgen(2))
print "reel 3, 5"
printreel(reelgen(3))
