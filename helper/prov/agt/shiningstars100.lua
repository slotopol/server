local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

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

math.randomseed(os.time())
printreel(makereelhot(symset, 4, {[1]=true, [2]=true, [3]=true}, {}))
