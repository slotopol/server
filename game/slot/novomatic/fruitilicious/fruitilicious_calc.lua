-- Novomatic / Fruitilicious RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{5, 5, 5, 5, 7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7, 7, 7, 6, 6, 6, 6, 1, 1, 1, 1, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 2, 2, 2, 2, 2, 2, 4, 4, 4, 4},
	{7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 4, 4, 4, 4, 4, 5, 5, 5, 5, 4, 4, 4, 4, 7, 7, 7, 7, 6, 6, 6, 6, 1, 1, 1, 1, 6, 6, 6, 6, 6, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2},
	{6, 6, 6, 6, 6, 7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 4, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 2, 2, 2, 2, 2, 2, 4, 4, 4, 4, 1, 1, 1, 1, 7, 7, 7, 7, 3, 3, 3, 3, 3, 3},
	{7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 7, 7, 7, 7, 5, 5, 5, 5, 5, 4, 4, 4, 4, 4, 2, 2, 2, 2, 2, 2, 4, 4, 4, 4, 1, 1, 1, 1, 6, 6, 6, 6, 5, 5, 5, 5, 3, 3, 3, 3, 3, 3},
	{5, 5, 5, 5, 4, 4, 4, 4, 4, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 5, 5, 5, 5, 5, 4, 4, 4, 4, 6, 6, 6, 6, 6, 7, 7, 7, 7, 6, 6, 6, 6, 7, 7, 7, 7, 7},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 0, 100, 500, 5000}, -- seven
	[2] = {0, 0, 25, 100, 500},   -- melon
	[3] = {0, 0, 25, 100, 500},   -- grapes
	[4] = {0, 0, 10, 30, 125},    -- plum
	[5] = {0, 0, 10, 30, 125},    -- orange
	[6] = {0, 0, 10, 30, 125},    -- lemon
	[7] = {0, 0, 10, 30, 125},    -- cherry
}

-- Performs full RTP calculation for given reels
local function calculate(reels)
	-- Get number of total reshuffles and lengths of each reel.
	local reshuffles, lens = 1, {}
	for i, r in ipairs(reels) do
		reshuffles = reshuffles * #r
		lens[i] = #r
	end

	-- Function to count symbol occurrences on each reel
	local function symbol_counts(symbol_id)
		local counts = {}
		for i, r in ipairs(reels) do
			local count = 0
			for _, sym in ipairs(r) do
				if sym == symbol_id then
					count = count + 1
				end
			end
			counts[i] = count
		end
		return counts
	end

	-- Function to calculate expected return from line wins for all symbols
	local function calculate_line_ev()
		local ev_sum = 0

		-- Iterate through all symbols that pay on lines
		for symbol_id, pays in pairs(PAYTABLE_LINE) do
			local c = symbol_counts(symbol_id)

			-- Calculate combinations for 5-of-a-kind (XXXXX)
			local comb5 = c[1] * c[2] * c[3] * c[4] * c[5]
			ev_sum = ev_sum + comb5 * pays[5]

			-- 4-of-a-kind (XXXX-) EV on left side
			local comb4l = c[1] * c[2] * c[3] * c[4] * (lens[5] - c[5])
			ev_sum = ev_sum + comb4l * pays[4]

			-- 4-of-a-kind (-XXXX) EV on right side
			local comb4r = (lens[1] - c[1]) * c[2] * c[3] * c[4] * c[5]
			ev_sum = ev_sum + comb4r * pays[4]

			-- 3-of-a-kind (XXX--) EV on left side
			local comb3l = c[1] * c[2] * c[3] * (lens[4] - c[4]) * lens[5]
			ev_sum = ev_sum + comb3l * pays[3]

			-- 3-of-a-kind (--XXX) EV on right side
			local comb3r = lens[1] * (lens[2] - c[2]) * c[3] * c[4] * c[5]
			ev_sum = ev_sum + comb3r * pays[3]
		end

		return ev_sum
	end

	-- Execute calculation
	local rtp_line = calculate_line_ev() / reshuffles * 100
	local rtp_scat = 0
	local rtp_total = rtp_line + rtp_scat
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_total))
	return rtp_total
end

calculate(REELS)
