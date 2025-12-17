-- CT Interactive / Groovy Automat
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{7, 7, 7, 7, 3, 3, 3, 8, 8, 8, 8, 5, 5, 5, 5, 5, 5, 2, 5, 1, 6, 2, 4, 4, 4, 6, 6, 6, 6, 6, 6},
	{4, 8, 8, 8, 8, 8, 8, 2, 1, 3, 2, 8, 7, 4, 4, 4, 6, 6, 6, 5, 5, 5, 3, 3, 3, 7, 7, 7, 7, 7, 7},
	{4, 4, 4, 8, 8, 8, 8, 8, 6, 6, 6, 6, 6, 6, 3, 3, 3, 2, 1, 2, 5, 5, 5, 5, 5, 5, 7, 7, 7, 7, 7},
	{3, 3, 3, 2, 5, 5, 5, 4, 7, 7, 7, 7, 7, 7, 6, 6, 6, 1, 7, 8, 8, 8, 8, 8, 8, 4, 4, 4, 8, 3, 2},
	{4, 4, 4, 8, 8, 8, 8, 2, 5, 5, 5, 5, 5, 5, 3, 3, 3, 6, 2, 5, 1, 7, 7, 7, 7, 6, 6, 6, 6, 6, 6},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 0, 0, 0, 0},         -- scatter
	[2] = {0, 0, 100, 1000, 5000}, -- 2 dollar
	[3] = {0, 0, 50, 200, 1000},   -- 3 melon
	[4] = {0, 0, 50, 200, 1000},   -- 4 apple
	[5] = {0, 0, 20, 50, 200},     -- 5 orange
	[6] = {0, 0, 20, 50, 200},     -- 6 lemon
	[7] = {0, 0, 20, 40, 200},     -- 7 plum
	[8] = {0, 5, 20, 40, 200},     -- 8 cherry
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 10, 50, 250}
local scat_min = 3 -- minimum scatters to win

-- 4. CONFIGURATION
local sy = 3 -- screen height
local scat = 1 -- scatter symbol ID

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

		-- Iterate through all symbols that pay on lines
		for symbol_id, pays in pairs(PAYTABLE_LINE) do
			local c = symbol_counts(symbol_id)

			-- 5-of-a-kind (XXXXX) EV
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

calculate(REELS)
