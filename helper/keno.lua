
local function fmt(...)
	print(string.format(...))
end

local function Combin(n, r)
	local mi, mj = 1, 1
	local i, j = n, 1
	for _ = 1, r do
		mi = mi * i
		mj = mj * j
		i = i - 1
		j = j + 1
	end
	return mi / mj
end

local C_80_20 = Combin(80, 20)

local function Prob(n, r)
	return Combin(n, r) * Combin(80-n, 20-r) / C_80_20
end

print "Probability calculation"
do
	local t = {}
	for r = 0, 10 do
		t[#t+1] = string.format("%-10d", r)
	end
	fmt("     %s", table.concat(t, " | "))
end
for n = 1, 10 do
	local t = {}
	for r = 0, n do
		local c = Prob(n, r)
		if c < 1e-8 then
			c = 0
		end
		t[#t+1] = string.format("%.8f", c)
	end
	fmt("[%02d] %s", n, table.concat(t, " | "))
end

-- Keno Luxury
local paytable = {
--       0     1     2     3     4     5     6     7     8     9    10  WIN BALLS
	{    0,    0,    0,    0,    0,    0,    0,    0,    0,    0,    0,},
	{    0,    1,    9,    0,    0,    0,    0,    0,    0,    0,    0,},
	{    0,    0,    2,   47,    0,    0,    0,    0,    0,    0,    0,},
	{    0,    0,    2,    5,   91,    0,    0,    0,    0,    0,    0,},
	{    0,    0,    0,    3,   12,  820,    0,    0,    0,    0,    0,},
	{    0,    0,    0,    3,    4,   70, 1600,    0,    0,    0,    0,},
	{    0,    0,    0,    1,    2,   21,  400, 7000,    0,    0,    0,},
	{    0,    0,    0,    0,    2,   12,  100, 1650,10000,    0,    0,},
	{    0,    0,    0,    0,    1,    6,   44,  335, 4700,10000,    0,},
	{    0,    0,    0,    0,    0,    5,   24,  142, 1000, 4500,10000,},
}

print "RTP calculation"
local grtp = 0
for n = 2, 10 do
	local rtp = 0
	for r = 0, n do
		local pay = paytable[n][r+1]
		rtp = rtp + pay * Prob(n, r)
	end
	fmt("RTP[%2d] = %f", n, rtp*100)
	grtp = grtp + rtp
end
grtp = grtp / 9
fmt("game RTP = %f", grtp*100)
