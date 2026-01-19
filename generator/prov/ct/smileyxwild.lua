local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, -- 1 wild    (2, 4 reel)
	1, -- 2 scatter
	4, -- 3 heart   1000
	7, -- 4 sun     300
	7, -- 5 beer    300
	11, -- 6 pizza   100
	11, -- 7 bomb    100
	11, -- 8 flower  100
}

local chunklen = {
	1, -- 1 wild
	1, -- 2 scatter
	1, -- 3 heart
	3, -- 4 sun
	3, -- 5 beer
	3, -- 6 pizza
	3, -- 7 bomb
	3, -- 8 flower
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 3 or n == 5 then
		ss[1] = 0
	end
	return makereelhot(ss, 3, {[1]=true, [2]=true}, chunklen)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
print "reel 1, 3, 5"
printreel(reelgen(1))
print "reel 2, 4"
printreel(reelgen(2))
