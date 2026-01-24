local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	3, --  1 wild (2, 3, 4, 5 reels only)
	4, --  2 scatter (1, 3, 5 reels only)
	4, --  3 cleopatra 1000
	5, --  4 cat       500
	6, --  5 ankh      400
	6, --  6 eye       400
	7, --  7 ace       300
	7, --  8 king      300
	8, --  9 queen     200
	8, -- 10 jack      200
	8, -- 11 ten       200
}

local chained = {
	0, --  1 wild (2, 3, 4, 5 reels only)
	0, --  2 scatter (1, 3, 5 reels only)
	0, --  3 cleopatra
	0, --  4 cat
	0, --  5 ankh
	0, --  6 eye
	3, --  7 ace
	3, --  8 king
	3, --  9 queen
	3, -- 10 jack
	3, -- 11 ten
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 then
		ss[1] = 0
	end
	if n == 2 or n == 4 then
		ss[2] = 0
	end
	return makereelct(ss, 3, {[1]=true, [2]=true}, chained)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen(1))
printreel(reelgen(2))
printreel(reelgen(3))
printreel(reelgen(4))
printreel(reelgen(5))
