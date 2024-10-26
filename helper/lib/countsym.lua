-- This script helps to count symbols on the reels by their kind.

local reel = {} -- place here some reel

local sn = 0
for _, sym in pairs(reel) do
	if sym > sn then
		sn = sym
	end
end

local symset = {}
for i = 1, sn do
	symset[i] = 0
end

for _, sym in pairs(reel) do
	symset[sym] = symset[sym]+1
end

for sym, n in pairs(symset) do
	print(string.format("\t%d, -- %2d", n, sym))
end
print("reel length: "..#reel)
