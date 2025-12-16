local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset1 = {
	 3, --  1 wild       1000
	 0, --  2 chilli     (3 reel only)
	 2, --  3 scatter
	 5, --  4 seven      500
	 5, --  5 avocado    200
	 5, --  6 peach      200
	 5, --  7 apple      200
	 5, --  8 watermelon 200
	14, --  9 orange     100
	14, -- 10 lemon      100
	14, -- 11 plum       100
	14, -- 12 cherry     100
}

local symset2 = {
	 3, --  1 wild       1000
	 0, --  2 chilli     (3 reel only)
	 2, --  3 scatter
	 7, --  4 seven      500
	14, --  5 avocado    200
	14, --  6 peach      200
	14, --  7 apple      200
	14, --  8 watermelon 200
	 5, --  9 orange     100
	 5, -- 10 lemon      100
	 5, -- 11 plum       100
	 5, -- 12 cherry     100
}

local symset3 = {
	3*3, --  1 wild       1000
	  1, --  2 chilli     (3 reel only)
	7*2, --  3 scatter
	6*4, --  4 seven      500
	6*6, --  5 avocado    200
	6*6, --  6 peach      200
	6*6, --  7 apple      200
	6*6, --  8 watermelon 200
	6*8, --  9 orange     100
	6*8, -- 10 lemon      100
	6*8, -- 11 plum       100
	6*8, -- 12 cherry     100
}

local chunklen = {
	3, --  1 wild
	1, --  2 chilli
	1, --  3 scatter
	3, --  4 seven
	6, --  5 avocado
	6, --  6 peach
	6, --  7 apple
	6, --  8 watermelon
	6, --  9 orange
	6, -- 10 lemon
	6, -- 11 plum
	6, -- 12 cherry
}

local function reelgen(n)
	if n == 1 or n == 5 then
		return makereelhot(symset1, 3, {[3]=true}, chunklen)
	elseif n == 2 or n == 4 then
		return makereelhot(symset2, 3, {[3]=true}, chunklen)
	else -- n == 3
		return makereelhot(symset3, 3, {[3]=true}, chunklen)
	end
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 4"
printreel(reelgen(2))
print "reel 3"
printreel(reelgen(3))
