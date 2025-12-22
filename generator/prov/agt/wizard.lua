local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	 6, --  1 wild
	10, --  2 scatter (2, 3, 4 reels only)
	 5, --  3 owl
	 5, --  4 cat
	 5, --  5 cauldron
	 5, --  6 emerald
	 5, --  7 ruby
	 5, --  8 ace
	 5, --  9 king
	 6, -- 10 queen
	 6, -- 11 jack
}

local function reelgen(n)
	local function make()
		return makereelhot(symset, 4, {}, {})
	end
	if n == 1 or n == 5 then
		local n2 = symset[2]
		symset[2] = 0
		local reel, iter = make()
		symset[2] = n2
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
