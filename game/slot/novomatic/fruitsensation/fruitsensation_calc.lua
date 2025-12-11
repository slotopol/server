-- Novomatic / Fruit Sensation RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{4, 4, 4, 1, 6, 6, 6, 3, 3, 3, 6, 6, 6, 5, 5, 5, 6, 1, 2, 2, 2, 1, 5, 5, 5, 7, 4, 4, 4, 3, 3, 7, 7, 7, 2, 2, 7, 7, 7},
	{1, 4, 4, 4, 5, 5, 5, 3, 3, 6, 6, 6, 2, 2, 2, 5, 5, 5, 1, 7, 7, 7, 1, 4, 4, 4, 3, 3, 3, 6, 6, 6, 7, 7, 7, 6, 2, 2, 7},
	{5, 5, 5, 6, 6, 6, 7, 7, 7, 1, 6, 6, 6, 3, 3, 1, 6, 1, 7, 3, 3, 3, 5, 5, 5, 7, 7, 7, 4, 4, 4, 2, 2, 4, 4, 4, 2, 2, 2},
	{6, 6, 6, 4, 4, 4, 7, 7, 7, 1, 2, 2, 2, 7, 7, 7, 2, 2, 3, 3, 3, 5, 5, 5, 1, 6, 6, 6, 7, 5, 5, 5, 1, 6, 4, 4, 4, 3, 3},
	{2, 2, 1, 5, 5, 5, 1, 6, 6, 6, 7, 7, 7, 6, 6, 6, 4, 4, 4, 6, 7, 7, 7, 3, 3, 7, 4, 4, 4, 1, 5, 5, 5, 3, 3, 3, 2, 2, 2},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 0, 100, 1000, 5000}, -- seven
	[2] = {0, 0, 50, 200, 500},    -- bells
	[3] = {0, 0, 50, 200, 500},    -- melon
	[4] = {0, 0, 20, 50, 200},     -- plum
	[5] = {0, 0, 20, 50, 200},     -- orange
	[6] = {0, 0, 20, 50, 200},     -- lemon
	[7] = {0, 0, 20, 50, 200},     -- cherry
}


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

		-- 4-of-a-kind (XXXX-) EV on left side
		local comb4 = c[1] * c[2] * c[3] * c[4] * (lens[5] - c[5])
		ev_sum = ev_sum + comb4 * pays[4]

		-- 3-of-a-kind (XXX--) EV on left side
		local comb3 = c[1] * c[2] * c[3] * (lens[4] - c[4]) * lens[5]
		ev_sum = ev_sum + comb3 * pays[3]
	end

	return ev_sum / reshuffles
end

-- Execute calculation
local line_rtp = calculate_line_wins_ev() * 100
local scat_rtp = 0
local total_rtp = line_rtp + scat_rtp
print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", line_rtp, scat_rtp, total_rtp))
