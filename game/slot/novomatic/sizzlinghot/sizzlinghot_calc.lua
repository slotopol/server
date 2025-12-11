-- Novomatic / Sizzling Hot RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{1, 4, 4, 4, 3, 1, 3, 3, 8, 6, 6, 6, 7, 7, 7, 6, 6, 2, 2, 5, 2, 5, 5, 5, 4},
	{1, 6, 6, 6, 2, 2, 1, 2, 7, 7, 7, 7, 8, 4, 4, 4, 4, 5, 5, 5, 3, 5, 3, 3, 6},
	{1, 6, 7, 7, 7, 8, 5, 5, 5, 1, 5, 2, 2, 4, 2, 4, 4, 4, 3, 3, 7, 3, 6, 6, 6},
	{1, 5, 5, 5, 5, 1, 5, 4, 4, 4, 8, 3, 3, 6, 6, 6, 7, 6, 7, 7, 7, 4, 4, 2, 2},
	{1, 4, 4, 6, 6, 6, 2, 2, 5, 8, 5, 5, 5, 8, 5, 4, 4, 4, 6, 1, 7, 7, 7, 3, 3},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 0, 100, 1000, 5000}, -- seven
	[2] = {0, 0, 50, 200, 500},    -- melon
	[3] = {0, 0, 50, 200, 500},    -- grapes
	[4] = {0, 0, 20, 50, 200},     -- plum
	[5] = {0, 0, 20, 50, 200},     -- orange
	[6] = {0, 0, 20, 50, 200},     -- lemon
	[7] = {0, 5, 20, 50, 200},     -- cherry
	[8] = {0, 0, 0, 0, 0},         -- star (scatter)
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 2, 10, 50}

-- 4. CONFIGURATION
local SCRH = 3 -- screen height
local scat = 8 -- scatter symbol ID


-- Get number of total reshuffles and lengths of each reel.
local reshuffles, lens = 1, {}
for i, r in ipairs(REELS) do
	reshuffles = reshuffles * #r
	lens[i] = #r
end

-- Function to count symbol occurrences on each reel
local function get_symbol_data(symbol_id)
	local counts = {}
	for i, r in ipairs(REELS) do
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
local function calculate_line_wins_ev()
	local ev_sum = 0

	-- Iterate through all symbols that pay on lines
	for symbol_id, pays in pairs(PAYTABLE_LINE) do
		local c = get_symbol_data(symbol_id)

		-- Calculate combinations for 5-of-a-kind (XXXXX)
		local comb5 = c[1] * c[2] * c[3] * c[4] * c[5]
		ev_sum = ev_sum + comb5 * pays[5]

		-- 4-of-a-kind (XXXX-) EV
		local comb4 = c[1] * c[2] * c[3] * c[4] * (lens[5] - c[5])
		ev_sum = ev_sum + comb4 * pays[4]

		-- 3-of-a-kind (XXX--) EV
		local comb3 = c[1] * c[2] * c[3] * (lens[4] - c[4]) * lens[5]
		ev_sum = ev_sum + comb3 * pays[3]

		-- 2-of-a-kind (XX---) EV
		local comb2 = c[1] * c[2] * (lens[3] - c[3]) * lens[4] * lens[5]
		ev_sum = ev_sum + comb2 * pays[2]
	end

	return ev_sum / reshuffles
end

-- Function to calculate expected return from scatter wins using probabilities
local function calculate_scatter_ev()
	local c = get_symbol_data(scat)
	local p = {} -- Probability of scatter appearing ANYWHERE on screen for each reel
	for i, len in ipairs(lens) do
		p[i] = (SCRH * c[i]) / len
	end

	-- We assume EV is the sum of (Probability(N_scatters) * Pay(N_scatters))
	local ev_scat = 0

	-- Using an recursive approach to sum probabilities for exactly N scatters
	local function combinations_prob(reel_index, current_scatters, current_prob)
		if reel_index > #REELS then
			if current_scatters >= 3 then
				ev_scat = ev_scat + current_prob * PAYTABLE_SCAT[current_scatters]
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
