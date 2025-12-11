-- Novomatic / Jewels RTP calculation

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


-- Get number of total reshuffles and lengths of each reel.
local reshuffles, lens = 1, {}
for i, r in ipairs(REELS) do
	reshuffles = reshuffles * #r
	lens[i] = #r
end

-- Function to count symbol occurrences on each reel and calculate
-- probabilities of landing exactly one symbol in the middle row
local function get_symbol_data(symbol_id)
	local counts, probs = {}, {}
	for i, r in ipairs(REELS) do
		local count = 0
		for _, sym in ipairs(r) do
			if sym == symbol_id then
				count = count + 1
			end
		end
		counts[i], probs[i] = count, count / lens[i]
	end
	return counts, probs
end

-- Function to calculate expected return from line wins for all symbols
local function calculate_line_wins_ev()
	local total_ev_line = 0

	-- Iterate through all symbols that pay on lines
	for symbol_id = 1, #PAYTABLE_LINE do
		local counts, _ = get_symbol_data(symbol_id)
		local pays = PAYTABLE_LINE[symbol_id]

		-- Calculate combinations for 5-of-a-kind (XXXXX)
		local c5 = counts[1] * counts[2] * counts[3] * counts[4] * counts[5]
		total_ev_line = total_ev_line + c5 * pays[5]

		-- 4-of-a-kind (XXXX-) EV on left side
		local cl4 = counts[1] * counts[2] * counts[3] * counts[4] * (lens[5] - counts[5])
		total_ev_line = total_ev_line + cl4 * pays[4]

		-- 4-of-a-kind (-XXXX) EV on right side
		local cr4 = (lens[1] - counts[1]) * counts[2] * counts[3] * counts[4] * counts[5]
		total_ev_line = total_ev_line + cr4 * pays[4]

		-- 3-of-a-kind (XXX--) EV on left side
		local cl3 = counts[1] * counts[2] * counts[3] * (lens[4] - counts[4]) * lens[5]
		total_ev_line = total_ev_line + cl3 * pays[3]

		-- 3-of-a-kind (-XXX-) EV in the middle
		local cm3 = (lens[1] - counts[1]) * counts[2] * counts[3] * counts[4] * (lens[5] - counts[5])
		total_ev_line = total_ev_line + cm3 * pays[3]

		-- 3-of-a-kind (--XXX) EV on right side
		local cr3 = lens[1] * (lens[2] - counts[2]) * counts[3] * counts[4] * counts[5]
		total_ev_line = total_ev_line + cr3 * pays[3]
	end

	return total_ev_line / reshuffles
end

-- Execute calculation
local line_rtp = calculate_line_wins_ev() * 100
local scat_rtp = 0
local total_rtp = line_rtp + scat_rtp
print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", line_rtp, scat_rtp, total_rtp))
