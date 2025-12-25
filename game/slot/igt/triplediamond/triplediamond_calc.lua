-- IGT / Triple Diamond
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{3, 0, 4, 0, 4, 0, 4, 0, 4, 0, 4, 0, 4, 0, 2, 0, 4, 0, 2, 0, 3, 0, 5, 0, 5, 0, 5, 0, 3, 0, 0, 1, 0, 0, 3, 0, 5, 0, 2, 0, 5, 0, 5, 0, 5, 0, 2, 0, 5, 0, 3, 0, 5, 0, 5, 0, 4, 0},
	{4, 0, 4, 0, 4, 0, 5, 0, 3, 0, 5, 0, 2, 0, 3, 0, 5, 0, 3, 0, 2, 0, 4, 0, 5, 0, 5, 0, 0, 1, 0, 0, 4, 0, 4, 0, 2, 0, 3, 0, 5, 0, 5, 0, 5, 0, 4, 0, 3, 0, 5, 0, 2, 0, 5, 0},
	{3, 0, 4, 0, 3, 0, 4, 0, 5, 0, 4, 0, 4, 0, 5, 0, 2, 0, 5, 0, 4, 0, 0, 1, 0, 0, 3, 0, 5, 0, 3, 0, 5, 0, 3, 0, 4, 0, 5, 0, 4, 0, 5, 0, 4, 0, 5, 0, 5, 0, 2, 0, 2, 0, 5, 0, 2, 0},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[0] = 0,    -- space
	[1] = 1199, -- diamond
	[2] = 100,  -- seven
	[3] = 40,   -- bar3
	[4] = 20,   -- bar2
	[5] = 10,   -- bar1
}

-- 3. CONFIGURATION
local sx = 3 -- screen width
-- Symbols names
local space   = 0
local diamond = 1
local seven   = 2
local bar3    = 3
local bar2    = 4
local bar1    = 5

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
	for symbol_id in pairs(PAYTABLE_LINE) do
		counts[symbol_id] = {}
		for i = 1, sx do counts[symbol_id][i] = 0 end
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
		local b1 = counts[bar1]
		local b2 = counts[bar2]
		local b3 = counts[bar3]
		local z = counts[space]
		local comb1w = 0 -- winning combinations with 1 diamond

		-- Iterate through all symbols that pay on lines
		for symbol_id, pay in pairs(PAYTABLE_LINE) do
			if symbol_id ~= space then
				local s = counts[symbol_id]
				local comb = s[1] * s[2] * s[3]
				ev_sum = ev_sum + comb * pay
			end
		end
		-- 1 diamond and any 2 symbols
		for symbol_id, pay in pairs(PAYTABLE_LINE) do
			if symbol_id ~= diamond and symbol_id ~= space then
				local s = counts[symbol_id]
				local comb =
					w[1] * s[2] * s[3] +
					s[1] * w[2] * s[3] +
					s[1] * s[2] * w[3]
				ev_sum = ev_sum + comb * pay * 3
				comb1w = comb1w + comb
			end
		end
		-- 2 diamonds and any 1 symbol
		for symbol_id, pay in pairs(PAYTABLE_LINE) do
			if symbol_id ~= diamond and symbol_id ~= space then
				local s = counts[symbol_id]
				local comb =
					w[1] * w[2] * s[3] +
					s[1] * w[2] * w[3] +
					w[1] * s[2] * w[3]
				ev_sum = ev_sum + comb * pay * 9
			end
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
			ev_sum = ev_sum + comb * 5 * 3
			comb1w = comb1w + comb
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
		-- 1 diamond
		do
			local comb =
				w[1] * (lens[2] - w[2]) * (lens[3] - w[3]) +
				(lens[1] - w[1]) * w[2] * (lens[3] - w[3]) +
				(lens[1] - w[1]) * (lens[2] - w[2]) * w[3] - comb1w
			ev_sum = ev_sum + comb * 2
		end
		-- 2 diamonds
		do
			local comb =
				w[1] * w[2] * z[3] +
				z[1] * w[2] * w[3] +
				w[1] * z[2] * w[3]
			ev_sum = ev_sum + comb * 10
		end

		return ev_sum
	end

	-- Execute calculation
	local rtp_total = calculate_line_ev() / reshuffles * 100
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("RTP = %.6f%%", rtp_total))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
