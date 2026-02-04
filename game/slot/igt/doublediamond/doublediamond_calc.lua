-- IGT / Double Diamond
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{2, 0, 6, 0, 4, 0, 5, 0, 4, 0, 6, 0, 2, 0, 5, 0, 3, 0, 3, 0, 2, 0, 5, 0, 4, 0, 5, 0, 5, 0, 3, 0, 0, 1, 0, 0, 4, 0, 4, 0, 6, 0, 5, 0, 3, 0},
	{2, 0, 2, 0, 4, 0, 5, 0, 5, 0, 5, 0, 4, 0, 4, 0, 6, 0, 3, 0, 4, 0, 6, 0, 2, 0, 5, 0, 5, 0, 0, 1, 0, 0, 5, 0, 3, 0, 2, 0, 3, 0, 3, 0, 6, 0},
	{6, 0, 4, 0, 4, 0, 4, 0, 2, 0, 6, 0, 3, 0, 2, 0, 6, 0, 4, 0, 5, 0, 2, 0, 3, 0, 5, 0, 5, 0, 3, 0, 0, 1, 0, 0, 3, 0, 5, 0, 5, 0, 5, 0, 4, 0},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[0] = 0,    -- space
	[1] = 1000, -- diamond
	[2] = 80,   -- seven
	[3] = 40,   -- bar3
	[4] = 25,   -- bar2
	[5] = 10,   -- bar1
	[6] = 10,   -- cherry
}

-- 3. CONFIGURATION
local sx = 3 -- grid width
-- Symbols names
local space   = 0
local diamond = 1
local seven   = 2
local bar3    = 3
local bar2    = 4
local bar1    = 5
local cherry  = 6

local _ = seven -- not used

-- Performs full RTP calculation for given reels
local function calculate(reels)
	assert(#reels == sx, "unexpected number of reels")

	-- Get number of total reshuffles and lengths of each reel.
	local reshuffles, lens = 1, {}
	for i, r in ipairs(reels) do
		reshuffles = reshuffles * #r
		lens[i] = #r
	end

	-- Count symbols occurrences on each reel
	local counts = {}
	for sym_id in pairs(PAYTABLE_LINE) do
		counts[sym_id] = {}
		for i = 1, sx do counts[sym_id][i] = 0 end
	end
	for i, r in ipairs(reels) do
		for _, sym in ipairs(r) do
			counts[sym][i] = counts[sym][i] + 1
		end
	end

	-- Function to calculate expected return from line wins for all symbols
	local function calculate_line_ev()
		local ev_sum = 0
		local w = counts[diamond]
		local c = counts[cherry]
		local a = {}
		for i = 1, sx do a[i] = lens[i] - c[i] - w[i] end
		local b1 = counts[bar1]
		local b2 = counts[bar2]
		local b3 = counts[bar3]

		-- Iterate through all symbols that pay on lines
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			if sym_id ~= space then
				local s = counts[sym_id]
				local comb = s[1] * s[2] * s[3]
				ev_sum = ev_sum + comb * pay
			end
		end
		-- 1 diamond and any 2 symbols
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			if sym_id ~= diamond and sym_id ~= space then
				local s = counts[sym_id]
				local comb =
					w[1] * s[2] * s[3] +
					s[1] * w[2] * s[3] +
					s[1] * s[2] * w[3]
				ev_sum = ev_sum + comb * pay * 2
			end
		end
		-- 2 diamonds and any 1 symbol
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			if sym_id ~= diamond and sym_id ~= space then
				local s = counts[sym_id]
				local comb =
					w[1] * w[2] * s[3] +
					s[1] * w[2] * w[3] +
					w[1] * s[2] * w[3]
				ev_sum = ev_sum + comb * pay * 4
			end
		end
		-- 1 diamond, 1 cherry, and any other symbol
		do
			local comb =
				w[1] * c[2] * a[3] +
				w[1] * a[2] * c[3] +
				c[1] * w[2] * a[3] +
				a[1] * w[2] * c[3] +
				c[1] * a[2] * w[3] +
				a[1] * c[2] * w[3]
			ev_sum = ev_sum + comb * 2 * 2
		end
		-- 1 cherry, and 2 any other symbols
		do
			local comb =
				c[1] * a[2] * a[3] +
				a[1] * c[2] * a[3] +
				a[1] * a[2] * c[3]
			ev_sum = ev_sum + comb * 2
		end
		-- 2 cherry, and 1 any other symbol
		do
			local comb =
				c[1] * c[2] * a[3] +
				a[1] * c[2] * c[3] +
				c[1] * a[2] * c[3]
			ev_sum = ev_sum + comb * 5
		end
		-- any bars with diamond
		do
			local comb =
				w[1] * b1[2] * b2[3] +
				w[1] * b1[2] * b3[3] +
				w[1] * b2[2] * b1[3] +
				w[1] * b2[2] * b3[3] +
				w[1] * b3[2] * b1[3] +
				w[1] * b3[2] * b2[3] +
				b1[1] * w[2] * b2[3] +
				b1[1] * w[2] * b3[3] +
				b2[1] * w[2] * b1[3] +
				b2[1] * w[2] * b3[3] +
				b3[1] * w[2] * b1[3] +
				b3[1] * w[2] * b2[3] +
				b1[1] * b2[2] * w[3] +
				b1[1] * b3[2] * w[3] +
				b2[1] * b1[2] * w[3] +
				b2[1] * b3[2] * w[3] +
				b3[1] * b1[2] * w[3] +
				b3[1] * b2[2] * w[3]
			ev_sum = ev_sum + comb * 5 * 2
		end
		-- any bar without diamond
		do
			local b = {}
			for i = 1, sx do b[i] = b1[i] + b2[i] + b3[i] end
			local comb = b[1] * b[2] * b[3] -
				b1[1] * b1[2] * b1[3] -
				b2[1] * b2[2] * b2[3] -
				b3[1] * b3[2] * b3[3]
			ev_sum = ev_sum + comb * 5
		end

		return ev_sum
	end

	-- Execute calculation
	local rtp_total = calculate_line_ev() / reshuffles
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("RTP = %.6f%%", rtp_total*100))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
