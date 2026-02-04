-- Novomatic / Jewels
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 0, 20, 200, 2000}, -- crown
	[2] = {0, 0, 15, 100, 500},  -- gold
	[3] = {0, 0, 15, 100, 500},  -- money
	[4] = {0, 0, 10, 50, 200},   -- ruby
	[5] = {0, 0, 10, 50, 200},   -- sapphire
	[6] = {0, 0, 5, 25, 100},    -- emerald
	[7] = {0, 0, 5, 25, 100},    -- amethyst
}

-- 3. CONFIGURATION
local sx = 5 -- grid width

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

		-- Iterate through all symbols that pay on lines
		for sym_id, pays in pairs(PAYTABLE_LINE) do
			local c = counts[sym_id]

			-- 5-of-a-kind (XXXXX) EV
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

			-- 3-of-a-kind (-XXX-) EV in the middle
			local comb3m = (lens[1] - c[1]) * c[2] * c[3] * c[4] * (lens[5] - c[5])
			ev_sum = ev_sum + comb3m * pays[3]

			-- 3-of-a-kind (--XXX) EV on right side
			local comb3r = lens[1] * (lens[2] - c[2]) * c[3] * c[4] * c[5]
			ev_sum = ev_sum + comb3r * pays[3]
		end

		return ev_sum
	end

	-- Execute calculation
	local rtp_line = calculate_line_ev() / reshuffles
	local rtp_scat = 0
	local rtp_total = rtp_line + rtp_scat
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line*100, rtp_scat*100, rtp_total*100))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
