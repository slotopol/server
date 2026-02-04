-- AGT / Halloween
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{1, 3, 3, 3, 3, 4, 4, 4, 4, 4, 2, 5, 5, 5, 5, 5, 5, 8, 8, 8, 8, 8, 8, 8, 8, 8, 2, 1, 2, 7, 7, 7, 7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 6, 6, 1},
	{5, 5, 5, 5, 5, 2, 1, 2, 6, 6, 6, 6, 6, 6, 6, 6, 1, 8, 8, 8, 8, 8, 8, 8, 8, 8, 3, 3, 3, 3, 7, 7, 7, 7, 7, 7, 7, 7, 7, 4, 4, 4, 4, 4, 1, 2},
	{2, 7, 7, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 5, 8, 8, 8, 8, 8, 8, 8, 8, 8, 2, 3, 3, 3, 3, 2, 1, 4, 4, 4, 4, 4, 1, 6, 6, 6, 6, 6, 6, 6, 1},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = 1000, -- pumpkin
	[2] = 500,  -- witch
	[3] = 200,  -- castle
	[4] = 100,  -- scarecrow
	[5] = 30,   -- ghost
	[6] = 20,   -- spider
	[7] = 10,   -- skeleton
	[8] = 5,    -- candles
}

-- 3. CONFIGURATION
local sx = 3 -- grid width

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
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			local c = counts[sym_id]
			local comb = c[1] * c[2] * c[3]
			ev_sum = ev_sum + comb * pay
		end

		return ev_sum
	end

	-- Execute calculation
	local rtp_total = calculate_line_ev() / reshuffles
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("RTP = %.6f%%", rtp_total*100))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
