-- CT Interactive / Purple Hot 2 RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{5, 6, 6, 6, 6, 8, 8, 8, 8, 6, 5, 5, 5, 5, 1, 4, 4, 4, 2, 8, 2, 7, 3, 3, 3, 7, 7, 7, 7},
	{8, 8, 8, 8, 2, 1, 3, 3, 3, 2, 4, 4, 4, 5, 5, 5, 5, 7, 7, 7, 7, 6, 6, 6, 6, 8, 6, 5, 7, 4},
	{8, 8, 8, 8, 1, 3, 3, 3, 2, 1, 2, 4, 3, 5, 5, 5, 5, 6, 4, 4, 4, 6, 6, 6, 6, 5, 7, 7, 7, 7},
	{8, 2, 8, 8, 8, 8, 2, 3, 3, 3, 5, 6, 1, 4, 4, 4, 6, 6, 6, 6, 7, 4, 5, 5, 5, 5, 7, 7, 7, 7},
	{4, 4, 4, 6, 2, 8, 6, 6, 6, 6, 7, 7, 7, 7, 5, 2, 3, 3, 3, 7, 8, 8, 8, 8, 1, 5, 5, 5, 5},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 0, 0, 0, 0},         -- scatter
	[2] = {0, 0, 100, 1000, 5000}, -- seven
	[3] = {0, 0, 50, 200, 500},    -- apple
	[4] = {0, 0, 50, 200, 500},    -- orange
	[5] = {0, 0, 20, 50, 200},     -- plum
	[6] = {0, 0, 20, 50, 200},     -- lemon
	[7] = {0, 0, 20, 50, 200},     -- melon
	[8] = {0, 5, 20, 50, 200},     -- cherry
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 2, 10, 50}

-- 4. CONFIGURATION
local SCRH = 3 -- screen height
local scat = 1 -- scatter symbol ID


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

		-- 4-of-a-kind (XXXX-) EV
		local c4 = counts[1] * counts[2] * counts[3] * counts[4] * (lens[5] - counts[5])
		total_ev_line = total_ev_line + c4 * pays[4]

		-- 3-of-a-kind (XXX--) EV
		local c3 = counts[1] * counts[2] * counts[3] * (lens[4] - counts[4]) * lens[5]
		total_ev_line = total_ev_line + c3 * pays[3]

		-- 2-of-a-kind (XX---) EV
		local c2 = counts[1] * counts[2] * (lens[3] - counts[3]) * lens[4] * lens[5]
		total_ev_line = total_ev_line + c2 * pays[2]
	end

	return total_ev_line / reshuffles
end

-- Function to calculate expected return from scatter wins using probabilities
local function calculate_scatter_ev()
	local counts, _ = get_symbol_data(scat)
	local p = {} -- Probability of scatter appearing ANYWHERE on screen for each reel
	for i, len in ipairs(lens) do
		p[i] = (SCRH * counts[i]) / len
	end

	-- We assume EV is the sum of (Probability(N_scatters) * Pay(N_scatters))
	local ev_scat = 0

	-- Using an recursive approach to sum probabilities for exactly N scatters
	local function combinations_prob(reel_index, current_scatters, current_prob)
		if reel_index > #REELS then
			if current_scatters >= 3 then
				local pay_index = math.min(current_scatters, #REELS)
				ev_scat = ev_scat + current_prob * PAYTABLE_SCAT[pay_index]
			end
			return
		end
		-- Probability of having a scatter on this reel
		combinations_prob(reel_index + 1, current_scatters + 1, current_prob * p[reel_index])
		-- Probability of NOT having a scatter on this reel
		combinations_prob(reel_index + 1, current_scatters, current_prob * (1 - p[reel_index]))
	end
	combinations_prob(1, 0, 1) -- Start recursion
	return ev_scat
end

-- Execute calculation
local line_rtp = calculate_line_wins_ev() * 100
local scat_rtp = calculate_scatter_ev() * 100
local total_rtp = line_rtp + scat_rtp
print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", line_rtp, scat_rtp, total_rtp))
