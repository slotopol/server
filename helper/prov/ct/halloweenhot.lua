local scripts = arg[0]:match("^(.*helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	4, -- 1 wild    (2, 3, 4 reels only)
	1, -- 2 scatter
	4, -- 3 seven   1000
	5, -- 4 dead    150
	5, -- 5 cat     80
	5, -- 6 vampire 80
	6, -- 7 pot     80
	6, -- 8 hat     80
	6, -- 9 scull   80
}

local chunklen = {
	4, -- 1 wild
	1, -- 2 scatter
	4, -- 3 seven
	4, -- 4 dead
	4, -- 5 cat
	4, -- 6 vampire
	4, -- 7 pot
	4, -- 8 hat
	4, -- 9 scull
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
