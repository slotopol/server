local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset15 = {
	2, -- 1 wild
	0, -- 2 scatter (always 0 here)
	3, -- 3 strawberry
	4, -- 4 papaya
	4, -- 5 grapes
	7, -- 6 orange
	7, -- 7 plum
	8, -- 8 cherry
	8, -- 9 pear
}

local symset234 = {
	 4, --  1 wild
	 1, --  2 scatter (always 1 here)
	 7, --  3 strawberry
	12, -- 4 papaya
	12, -- 5 grapes
	14, -- 6 orange
	14, -- 7 plum
	15, -- 8 cherry
	15, -- 9 pear
}

local chunklen15 = {
	1, -- 1 wild
	1, -- 2 scatter
	1, -- 3 strawberry
	1, -- 4 papaya
	1, -- 5 grapes
	3, -- 6 orange
	3, -- 7 plum
	3, -- 8 cherry
	3, -- 9 pear
}

local chunklen234 = {
	1, -- 1 wild
	1, -- 2 scatter
	1, -- 3 strawberry
	1, -- 4 papaya
	1, -- 5 grapes
	7, -- 6 orange
	7, -- 7 plum
	5, -- 8 cherry
	5, -- 9 pear
}

local function reelgen(n)
	if n == 1 or n == 5 then
		return makereelhot(symset15, 3, {[1]=true, [2]=true}, chunklen15)
	else
		return makereelhot(symset234, 3, {[1]=true, [2]=true}, chunklen234)
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
