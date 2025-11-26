local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	4, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 seven     25/50/1000
	4, --  4 coin      5/15/100
	4, --  5 bell      5/15/100
	4, --  6 horseshoe 5/15/100
	4, --  7 apple     5/10/50
	4, --  8 lemon     5/10/50
	5, --  9 plum      5/10/50
	5, -- 10 cherry    5/10/50
}

local chunklen = {
	4, --  1 wild (2, 3, 4 reels only)
	1, --  2 scatter
	4, --  3 seven
	4, --  4 coin
	4, --  5 bell
	4, --  6 horseshoe
	6, --  7 apple
	6, --  8 lemon
	6, --  9 plum
	6, -- 10 cherry
}

math.randomseed(os.time())

local function reelgen(n)
	local function make()
		return makereelhot(symset, 3, {[1]=true, [2]=true}, chunklen, true)
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

print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
