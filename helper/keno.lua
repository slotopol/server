
local function fmt(...)
	print(string.format(...))
end

function Combin(n, r)
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
local c = Combin

local C_80_20 = c(80, 20)

function KenoProb(n, r)
	return c(n, r) * c(80-n, 20-r) / C_80_20
end
local p = KenoProb

function KenoProbTable(selmax, prec)
	print "Probability calculation"
	local hf = "%-"..(prec+2).."d"
	do
		local t = {}
		for r = 0, selmax do
			t[#t+1] = string.format(hf, r)
		end
		fmt("     %s", table.concat(t, " | "))
	end
	local pf = "%."..prec.."f"
	for n = 1, selmax do
		local t = {}
		for r = 0, n do
			t[#t+1] = string.format(pf, p(n, r))
		end
		fmt("[%02d] %s", n, table.concat(t, " | "))
	end
end

function KenoRTP(paytable, prec)
	print "RTP calculation"
	local grtp = 0
	local pf = "RTP[%2d] = %."..prec.."f%%"
	for n = 2, 10 do
		local rtp = 0
		for r = 0, n do
			local pay = paytable[n][r+1]
			rtp = rtp + pay * p(n, r)
		end
		fmt(pf, n, rtp*100)
		grtp = grtp + rtp
	end
	grtp = grtp / 9
	fmt("RTP[game] = %."..prec.."f%%", grtp*100)
end
