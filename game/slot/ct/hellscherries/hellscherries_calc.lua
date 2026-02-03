-- CT Interactive / Hell's Cherries
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{3, 8, 8, 8, 8, 8, 2, 3, 4, 10, 10, 10, 10, 10, 10, 4, 4, 9, 9, 9, 9, 9, 7, 7, 7, 7, 1, 6, 6, 6, 6, 2, 5, 5, 5, 5},
	{3, 5, 5, 5, 5, 8, 8, 8, 8, 8, 6, 6, 6, 6, 2, 4, 4, 10, 10, 10, 10, 10, 10, 9, 9, 9, 9, 9, 2, 7, 7, 7, 7, 3, 1, 4},
	{5, 5, 5, 5, 8, 8, 8, 8, 8, 4, 4, 10, 10, 10, 10, 10, 10, 6, 6, 6, 6, 2, 3, 9, 9, 9, 9, 9, 4, 3, 1, 7, 7, 7, 7, 2},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = 500, --  1 seven
	[ 2] = 250, --  2 bar
	[ 3] = 200, --  3 melon
	[ 4] = 100, --  4 bell
	[ 5] = 50,  --  5 apple
	[ 6] = 50,  --  6 pear
	[ 7] = 50,  --  7 plum
	[ 8] = 10,  --  8 lemon
	[ 9] = 10,  --  9 orange
	[10] = 10,  -- 10 cherry
}

-- 3. CONFIGURATION
local sx = 3 -- grid width
local wild = 1 -- wild symbol ID

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
		local w = counts[wild]

		local comb_w3 = w[1] * w[2] * w[3]
		ev_sum = ev_sum + comb_w3 * PAYTABLE_LINE[wild]

		-- Iterate through all symbols that pay on lines
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			if sym_id ~= wild then
				local s = counts[sym_id]
				local comb = (s[1] + w[1]) * (s[2] + w[2]) * (s[3] + w[3]) - comb_w3
				ev_sum = ev_sum + comb * pay
			end
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
