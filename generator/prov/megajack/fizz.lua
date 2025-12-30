
local t = { -- champagne buttles table
	10, 10, 20, 150, 150,
}

local m = 0
for i = 1,#t do
	m = m + t[i]
end
m = m / #t

local M = 0
local n = 0
for i = 1,#t do
	for j = i+1,#t do
		if t[i] == t[j] then
			M = M + t[i]*4
		else
			M = M + t[i]+t[j]
		end
		n = n + 1
	end
end
M = M / n
print(string.format("len = %d, count = %d, avr bottle gain = %g, M = %g", #t, n, m, M))
