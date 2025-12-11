
local function printf(...)
	print(string.format(...))
end

function combin(n, r)
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
local c = combin

local c_80_20 = c(80, 20)

function kenoprob(n, r)
	return c(n, r) * c(80-n, 20-r) / c_80_20
end
local p = kenoprob

function kenoprobtable(selmax, prec)
	print "Probability calculation"
	local hf = "%-"..(prec+2).."d"
	do
		local t = {}
		for r = 0, selmax do
			t[#t+1] = string.format(hf, r)
		end
		printf("     %s", table.concat(t, " | "))
	end
	local pf = "%."..prec.."f"
	for n = 1, selmax do
		local t = {}
		for r = 0, n do
			t[#t+1] = string.format(pf, p(n, r))
		end
		printf("[%02d] %s", n, table.concat(t, " | "))
	end
end

function kenortp(paytable, prec)
	print "RTP calculation"
	local grtp = 0
	local pf = "RTP[%2d] = %."..prec.."f%%"
	for n = 2, 10 do
		local rtp = 0
		for r = 0, n do
			local pay = paytable[n][r+1]
			rtp = rtp + pay * p(n, r)
		end
		printf(pf, n, rtp*100)
		grtp = grtp + rtp
	end
	grtp = grtp / 9
	printf("RTP[game] = %."..prec.."f%%", grtp*100)
end
