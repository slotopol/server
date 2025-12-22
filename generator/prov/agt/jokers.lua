local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	4, --  1 wild (2, 3, 4 reels only)
	2, --  2 scatter
	5, --  3 strawberry
	5, --  4 pear
	5, --  5 greenstar
	5, --  6 redstar
	6, --  7 plum
	6, --  8 peach
	6, --  9 papaya
	6, -- 10 cherry
}

local chunklen = {
	3, --  1 wild
	1, --  2 scatter
	3, --  3 strawberry
	3, --  4 pear
	3, --  5 greenstar
	3, --  6 redstar
	3, --  7 plum
	3, --  8 peach
	3, --  9 papaya
	3, -- 10 cherry
}

local function reelgen(n)
	local function make()
		return makereelhot(symset, 3, {[2]=true}, chunklen)
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

math.randomseed(os.time())
print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
