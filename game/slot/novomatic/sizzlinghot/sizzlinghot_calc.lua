-- Sizzling Hot RTP Calculation

-- 1. REEL STRIPS DATA
local REELS = {
	{1, 4, 4, 4, 3, 1, 3, 3, 8, 6, 6, 6, 7, 7, 7, 6, 6, 2, 2, 5, 2, 5, 5, 5, 4},
	{1, 6, 6, 6, 2, 2, 1, 2, 7, 7, 7, 7, 8, 4, 4, 4, 4, 5, 5, 5, 3, 5, 3, 3, 6},
	{1, 6, 7, 7, 7, 8, 5, 5, 5, 1, 5, 2, 2, 4, 2, 4, 4, 4, 3, 3, 7, 3, 6, 6, 6},
	{1, 5, 5, 5, 5, 1, 5, 4, 4, 4, 8, 3, 3, 6, 6, 6, 7, 6, 7, 7, 7, 4, 4, 2, 2},
	{1, 4, 4, 6, 6, 6, 2, 2, 5, 8, 5, 5, 5, 8, 5, 4, 4, 4, 6, 1, 7, 7, 7, 3, 3}
}
local reshuffles, lens = 1, {}
for i, r in ipairs(REELS) do
	reshuffles = reshuffles * #r
	lens[i] = #r
end

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
local PAYTABLE_SCATTER = {0, 0, 2, 10, 50}

-- 4. CONFIGURATION
-- Screen dimensions
local SCRW = #REELS
local SCRH = 3


-- Function to count symbol occurrences on each reel and calculate probabilities
local function get_symbol_data(symbol_id)
	local counts = {}
	local probabilities = {}
	for i, r in ipairs(REELS) do
		local count = 0
		for _, sym in ipairs(r) do
			if sym == symbol_id then
				count = count + 1
			end
		end
		counts[i] = count
		-- Probability of landing exactly one symbol in the middle row
		probabilities[i] = count / lens[i]
	end
	return counts, probabilities
end

-- Function to calculate expected return from line wins for all symbols
local function calculate_line_wins_ev()
	local total_ev_line = 0

	-- Iterate through all symbols that pay on lines
	for symbol_id = 1, #PAYTABLE_LINE do
		local counts, _ = get_symbol_data(symbol_id)
		local symbol_pays = PAYTABLE_LINE[symbol_id]

		-- Calculate combinations for 5-of-a-kind (XXX..)
		local combinations_5 = counts[1] * counts[2] * counts[3] * counts[4] * counts[5]
		total_ev_line = total_ev_line + combinations_5 * symbol_pays[5]

		-- 4-of-a-kind (XXXX-) EV
		local combinations_4 = counts[1] * counts[2] * counts[3] * counts[4] * (lens[5] - counts[5])
		total_ev_line = total_ev_line + combinations_4 * symbol_pays[4]

		-- 3-of-a-kind (XXX--) EV
		local combinations_3 = counts[1] * counts[2] * counts[3] * (lens[4] - counts[4]) * lens[5]
		total_ev_line = total_ev_line + combinations_3 * symbol_pays[3]

		-- 2-of-a-kind (XX---) EV
		local combinations_2 = counts[1] * counts[2] * (lens[3] - counts[3]) * lens[4] * lens[5]
		total_ev_line = total_ev_line + combinations_2 * symbol_pays[2]
	end

	return total_ev_line / reshuffles
end

-- Function to calculate expected return from scatter wins using probabilities
local function calculate_scatter_ev()
	local counts, _ = get_symbol_data(8)
	local p = {} -- Probability of scatter appearing ANYWHERE on screen for each reel

	for i, len in ipairs(lens) do
		p[i] = (SCRH * counts[i]) / len
	end

	-- We need probability of exactly 3, 4, 5 scatters on screen.
	-- This requires a recursive or iterative function (binomial distribution over 5 reels with different P)

	local ev_scat = 0
	-- This part is complex to write inline in Lua combinatorially.
	-- We assume EV is the sum of (Probability(N_scatters) * Pay(N_scatters))

	-- Using an iterative approach to sum probabilities for exactly N scatters
	local function combinations_prob(reel_index, current_scatters, current_prob)
		if reel_index > SCRW then
			if current_scatters >= 3 then
				local pay_index = math.min(current_scatters, SCRW)
				ev_scat = ev_scat + current_prob * PAYTABLE_SCATTER[pay_index]
			end
			return
		end
		-- Probability of having a scatter on this reel
		combinations_prob(reel_index + 1, current_scatters + 1, current_prob * p[reel_index])
		-- Probability of NOT having a scatter on this reel
		combinations_prob(reel_index + 1, current_scatters, current_prob * (1 - p[reel_index]))
	end

	combinations_prob(1, 0, 1) -- Start recursion

	-- Note: This recursive function correctly calculates the EV per spin (for 1 line bet value)
	return ev_scat
end

-- Execute calculation
local line_rtp = calculate_line_wins_ev() * 100
local scat_rtp = calculate_scatter_ev() * 100
local total_rtp = line_rtp + scat_rtp

print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", line_rtp, scat_rtp, total_rtp))
