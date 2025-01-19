local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

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

math.randomseed(os.time())
printreel(makereelhot(symset, 3, {[1]=true, [2]=true, [3]=true}, {}))
