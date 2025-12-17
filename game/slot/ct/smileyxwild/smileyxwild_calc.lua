-- CT Interactive / Smiley x Wild
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{8, 8, 3, 6, 6, 6, 5, 4, 4, 4, 3, 7, 7, 7, 2, 3, 6, 6, 8, 8, 8, 7, 7, 4, 5, 5, 5},
	{4, 4, 4, 7, 7, 7, 8, 8, 8, 7, 7, 7, 6, 6, 6, 8, 8, 8, 2, 5, 5, 5, 3, 6, 6, 6, 3, 1, 4, 4, 4, 7, 7, 7, 5, 5, 5, 6, 6, 6, 3, 8, 8, 8, 3},
	{5, 4, 4, 4, 2, 6, 6, 6, 8, 8, 8, 4, 5, 5, 5, 3, 6, 6, 8, 8, 3, 7, 7, 7, 3, 7, 7},
	{7, 7, 7, 3, 5, 5, 5, 2, 8, 8, 8, 1, 6, 6, 6, 8, 8, 8, 6, 6, 6, 3, 7, 7, 7, 4, 4, 4, 8, 8, 8, 5, 5, 5, 3, 7, 7, 7, 6, 6, 6, 3, 4, 4, 4},
	{8, 8, 8, 2, 7, 7, 7, 6, 6, 4, 3, 6, 6, 6, 3, 5, 5, 5, 3, 5, 8, 8, 4, 4, 4, 7, 7},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 0, 0, 0, 0},       -- wild (2, 4 reels only)
	[2] = {0, 0, 0, 0, 0},       -- scatter
	[3] = {0, 0, 35, 100, 1000}, -- heart
	[4] = {0, 0, 15, 50, 300},   -- sun
	[5] = {0, 0, 15, 50, 300},   -- beer
	[6] = {0, 0, 10, 30, 100},   -- pizza
	[7] = {0, 0, 10, 30, 100},   -- bomb
	[8] = {0, 0, 10, 30, 100},   -- flower
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 5, 20, 100}
local scat_min = 3 -- minimum scatters to win

-- 4. CONFIGURATION
local sy = 3 -- screen height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local M = 3 -- average multiplier on wilds (1+2+3+4+5)/5

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
		-- count wilds on each reel
		local cw = symbol_counts(wild)
		assert(cw[1] == 0 and cw[3] == 0 and cw[5] == 0,
			"wilds should not appear on reels 1, 3, 5")

		-- Iterate through all symbols that pay on lines
		for symbol_id, pays in pairs(PAYTABLE_LINE) do
			-- count symbol occurrences without wilds
			local c = symbol_counts(symbol_id)

			-- 5-of-a-kind (XXXXX) EV: W on R2 and W on R4
			local comb5_ww = c[1] * cw[2] * c[3] * cw[4] * c[5]
			ev_sum = ev_sum + comb5_ww * pays[5] * M * M

			-- 5-of-a-kind (XXXXX) EV: W on R2
			local comb5_w2 = c[1] * cw[2] * c[3] * c[4] * c[5]
			ev_sum = ev_sum + comb5_w2 * pays[5] * M

			-- 5-of-a-kind (XXXXX) EV: W on R4
			local comb5_w4 = c[1] * c[2] * c[3] * cw[4] * c[5]
			ev_sum = ev_sum + comb5_w4 * pays[5] * M

			-- 5-of-a-kind (XXXXX) EV: no W
			local comb5_x1 = c[1] * c[2] * c[3] * c[4] * c[5]
			ev_sum = ev_sum + comb5_x1 * pays[5] * 1 -- no multiplier

			-- 4-of-a-kind (XXXX-) EV: W on R2 and W on R4
			local comb4_ww = c[1] * cw[2] * c[3] * cw[4] * (lens[5] - c[5])
			ev_sum = ev_sum + comb4_ww * pays[4] * M * M

			-- 4-of-a-kind (XXXX-) EV: W on R2
			local comb4_w2 = c[1] * cw[2] * c[3] * c[4] * (lens[5] - c[5])
			ev_sum = ev_sum + comb4_w2 * pays[4] * M

			-- 4-of-a-kind (XXXX-) EV: W on R4
			local comb4_w4 = c[1] * c[2] * c[3] * cw[4] * (lens[5] - c[5])
			ev_sum = ev_sum + comb4_w4 * pays[4] * M

			-- 4-of-a-kind (XXXX-) EV: no W
			local comb4_x1 = c[1] * c[2] * c[3] * c[4] * (lens[5] - c[5])
			ev_sum = ev_sum + comb4_x1 * pays[4] * 1 -- no multiplier

			-- 3-of-a-kind (XXX--) EV: W on R2
			local comb3_w2 = c[1] * cw[2] * c[3] * (lens[4] - c[4] - cw[4]) * lens[5]
			ev_sum = ev_sum + comb3_w2 * pays[3] * M

			-- 3-of-a-kind (XXX--) EV: no W
			local comb3_x1 = c[1] * c[2] * c[3] * (lens[4] - c[4] - cw[4]) * lens[5]
			ev_sum = ev_sum + comb3_x1 * pays[3] * 1 -- no multiplier
		end

		return ev_sum
	end

	-- Function to calculate expected return from scatter wins
	local function calculate_scat_ev()
		local c = symbol_counts(scat)
		local ev_sum = 0

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > #reels then
				if scat_sum >= scat_min then
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum]
				end
				return
			end
			-- Step 1: having a scatter on this reel
			find_scatter_combs(reel_index + 1, scat_sum + 1,
				current_comb * c[reel_index] * sy)
			-- Step 2: NOT having a scatter on this reel
			find_scatter_combs(reel_index + 1, scat_sum,
				current_comb * (lens[reel_index] - c[reel_index] * sy))
		end
		find_scatter_combs(1, 0, 1) -- Start recursion

		return ev_sum
	end

	-- Execute calculation
	local rtp_line = calculate_line_ev() / reshuffles * 100
	local rtp_scat = calculate_scat_ev() / reshuffles * 100
	local rtp_total = rtp_line + rtp_scat
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_total))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
