local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."lib/reelgen.lua")

local symset = {
	1, --  1 wild    (2, 3, 4 reel)
	2, --  2 bear    40/500/10000
	4, --  3 wolf    20/80/200
	4, --  4 owl     20/80/200
	4, --  5 walrus  15/40/100
	4, --  6 puffin  15/40/100
	4, --  7 ace     10/20/100
	5, --  8 king    10/20/100
	5, --  9 queen   10/20/100
	5, -- 10 jack    10/20/100
}

local chunklen = {
	1, --  1 wild
	1, --  2 bear
	3, --  3 wolf
	3, --  4 owl
	3, --  5 walrus
	3, --  6 puffin
	3, --  7 ace
	3, --  8 king
	3, --  9 queen
	3, -- 10 jack
}

math.randomseed(os.time())

function reelgen(n)
	if n == 1 or n == 5 then
		local n1 = symset[1]
		symset[1] = 0
		local reel, iter = makereelhot(symset, 3, {[1]=true}, chunklen, true)
		symset[1] = n1
		return reel, iter
	else
		return makereelhot(symset, 3, {[1]=true}, chunklen, true)
	end
end

if not autoscan then
	print "reel 1, 5"
	printreel(reelgen(1))
	print "reel 2, 3, 4"
	printreel(reelgen(2))
end
