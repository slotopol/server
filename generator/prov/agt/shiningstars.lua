local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, --  1 wild (on 2, 3, 4 reels)
	2, --  2 scatter1 (on all reels)
	2, --  3 scatter2 (on 1, 3, 5 reels)
	2, --  4 seven
	4, --  5 grape
	4, --  6 watermelon
	6, --  7 avocado
	7, --  8 pomegranate
	7, --  9 carambola
	7, -- 10 maracuya
	7, -- 11 orange
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	if n == 2 or n == 4 then
		ss[3] = 0
	end
	return makereelhot(ss, 3, {[1]=true, [2]=true, [3]=true}, {})
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
