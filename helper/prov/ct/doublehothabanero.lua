local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	3, --  1 wild     (2, 3, 4 reels only)
	2, --  2 scatter
	6, --  3 woman    1000
	6, --  4 man      500
	8, --  5 crayfish 200
	8, --  6 shrimp   150
	8, --  7 ananas   50
	8, --  8 lime     50
	8, --  9 corn     50
	8, -- 10 banana   50
}

local chunklen = {
	1, --  1 wild
	1, --  2 scatter
	1, --  3 woman
	1, --  4 man
	3, --  5 crayfish
	3, --  6 shrimp
	4, --  7 ananas
	4, --  8 lime
	4, --  9 corn
	4, -- 10 banana
}

math.randomseed(os.time())

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

print "reel 1, 5"
printreel(reelgen(1))
print "reel 2, 3, 4"
printreel(reelgen(2))
