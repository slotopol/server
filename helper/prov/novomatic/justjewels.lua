local path = arg[0]:match("(.*[/\\])")
dofile(path.."../../lib/reelgen.lua")

local symset15 = {
	2, -- 1 crown
	3, -- 2 gold
	3, -- 3 money
	3, -- 4 ruby
	3, -- 5 sapphire
	18, -- 6 emerald
	18, -- 7 amethyst
	1, -- 8 euro
}

local symset24 = {
	2, -- 1 crown
	3, -- 2 gold
	3, -- 3 money
	5, -- 4 ruby
	5, -- 5 sapphire
	4, -- 6 emerald
	4, -- 7 amethyst
	1, -- 8 euro
}

local symset3 = {
	5, -- 1 crown
	5, -- 2 gold
	5, -- 3 money
	4, -- 4 ruby
	4, -- 5 sapphire
	3, -- 6 emerald
	3, -- 7 amethyst
	3, -- 8 euro
}

local chunklen = {
	1, -- 1 crown
	8, -- 2 gold
	8, -- 3 money
	8, -- 4 ruby
	8, -- 5 sapphire
	8, -- 6 emerald
	8, -- 7 amethyst
	1, -- 8 euro
}

math.randomseed(os.time())
print "reel 1, 5"
printreel(makereelhot(symset15, 3, {[8]=true}, chunklen, true))
print "reel 2, 4"
printreel(makereelhot(symset24, 3, {[8]=true}, chunklen, true))
print "reel 3"
printreel(makereelhot(symset3, 3, {[8]=true}, chunklen, true))
