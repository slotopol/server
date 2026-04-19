
local function printf(...)
	print(string.format(...))
end

function combin(n, k)
	local mi, mj = 1, 1
	local i, j = n + 0.0, 1
	for _ = 1, k do
		mi = mi * i
		mj = mj * j
		i = i - 1
		j = j + 1
	end
	return mi / mj
end
local c = combin

local c_80_20 = c(80, 20)

function kenoprob(n, k)
	return c(n, k) * c(80-n, 20-k) / c_80_20
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
		for k = 0, n do
			t[#t+1] = string.format(pf, p(n, k))
		end
		printf("[%02d] %s", n, table.concat(t, " | "))
	end
end

function kenortp(paytable, prec)
	print "RTP calculation"
	local rtp, l = 0, 0
	local pf = "RTP[%2d] = %."..prec.."f%%, sigma = %."..prec.."f"
	for n = 1, 10 do
		local ev, e2 = 0, 0
		for k = 0, n do
			local pay = paytable[n][k+1]
			ev = ev + p(n, k) * pay
			e2 = e2 + p(n, k) * pay * pay
		end
		if ev > 0 then
			local D = e2 - ev*ev
			printf(pf, n, ev*100, math.sqrt(D))
			rtp = rtp + ev
			l = l + 1
		end
	end
	rtp = rtp / l
	printf("RTP[game] = %."..prec.."f%%", rtp*100)
end
