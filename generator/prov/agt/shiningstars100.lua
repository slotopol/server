local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild (on 2, 3, 4 reels)
	2, --  2 scatter1 (on all reels)
	3, --  3 scatter2 (on 1, 3, 5 reels)
	3, --  4 seven
	6, --  5 grape
	6, --  6 watermelon
	8, --  7 avocado
	11, --  8 pomegranate
	11, --  9 carambola
	11, -- 10 maracuya
	11, -- 11 orange
}

local function reelgen(n)
	local n1, n3 = symset[1], symset[3]
	if n == 1 or n == 5 then
		symset[1] = 0
	end
	if n == 2 or n == 4 then
		symset[3] = 0
	end
	local reel, iter = makereelhot(symset, 4, {[1]=true, [2]=true, [3]=true}, {})
	symset[1], symset[3] = n1, n3
	return reel, iter
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
