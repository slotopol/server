-- Novomatic / Fruit Sensation
-- RTP calculation

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

-- 3. CONFIGURATION
local sx = 5 -- screen width

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
			local comb4 = c[1] * c[2] * c[3] * c[4] * (lens[5] - c[5])
			ev_sum = ev_sum + comb4 * pays[4]

			-- 3-of-a-kind (XXX--) EV on left side
			local comb3 = c[1] * c[2] * c[3] * (lens[4] - c[4]) * lens[5]
			ev_sum = ev_sum + comb3 * pays[3]
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

if autoscan then
	return calculate
end

calculate(REELS)
