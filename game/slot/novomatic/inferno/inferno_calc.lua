-- Novomatic / Inferno
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{6, 7, 5, 6, 8, 4, 3, 6, 5, 2, 7, 6, 4, 7, 5, 1, 7, 5, 3, 7, 4, 3, 1, 5, 3, 2, 4, 7, 2, 5, 3, 6, 4, 8, 2, 4, 6, 2, 7, 6, 5, 4},
	{7, 6, 4, 3, 2, 8, 4, 5, 6, 1, 5, 3, 7, 8, 2, 7, 3, 5, 4, 6, 5, 2, 4, 6, 7, 3, 5, 4, 7, 6, 1, 2, 4, 3, 6, 7, 4, 5, 7, 6, 5, 2},
	{6, 5, 4, 7, 5, 4, 1, 7, 2, 5, 7, 2, 6, 3, 4, 5, 7, 3, 8, 5, 7, 3, 6, 4, 2, 6, 4, 8, 3, 6, 7, 2, 4, 5, 6, 3, 7, 6, 2, 1, 5, 4},
	{6, 3, 4, 7, 8, 3, 7, 2, 5, 4, 7, 6, 4, 7, 6, 5, 7, 4, 5, 2, 3, 6, 5, 2, 6, 5, 4, 8, 6, 4, 2, 1, 6, 3, 7, 1, 2, 3, 5, 4, 7, 5},
	{2, 7, 5, 6, 1, 3, 6, 7, 5, 6, 4, 7, 5, 6, 2, 5, 1, 3, 4, 8, 2, 4, 6, 5, 7, 3, 2, 6, 4, 7, 8, 5, 4, 7, 3, 6, 4, 2, 5, 4, 7, 3},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 0, 100, 1000, 10000}, -- star
	[2] = {0, 0, 40, 200, 500},     -- bell
	[3] = {0, 0, 40, 200, 500},     -- grapes
	[4] = {0, 0, 20, 50, 200},      -- plum
	[5] = {0, 0, 20, 50, 200},      -- orange
	[6] = {0, 0, 20, 50, 200},      -- lemon
	[7] = {0, 5, 20, 50, 200},      -- cherry
	[8] = {0, 0, 0, 0, 0},          -- crown (scatter)
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 2, 10, 50}
local scat_min = 3 -- minimum scatters to win

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local scat = 8 -- scatter symbol ID

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
	for symbol_id in pairs(PAYTABLE_LINE) do
		counts[symbol_id] = {}
		for i = 1, sx do counts[symbol_id][i] = 0 end
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
		for symbol_id, pays in pairs(PAYTABLE_LINE) do
			local c = counts[symbol_id]

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
		local c = counts[scat]
		local ev_sum = 0

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > sx then
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
